package wise

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
)

// scaClient extends Client with Strong Customer Authentication support.
// Note: SCA is currently not working due to signature rejection by Wise.
// This code is preserved for when SCA is enabled on the account.
type scaClient struct {
	*Client
	privateKey *rsa.PrivateKey
}

func newSCAClient(apiKey, privateKeyPath string) (*scaClient, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &scaClient{
		Client:     NewClient(apiKey),
		privateKey: privateKey,
	}, nil
}

func newSandboxSCAClient(apiKey, privateKeyPath string) (*scaClient, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &scaClient{
		Client:     NewSandboxClient(apiKey),
		privateKey: privateKey,
	}, nil
}

func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("reading private key: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("decoding PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	return privateKey, nil
}

func signData(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("signing: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (c *scaClient) fundTransfer(profileID, transferID int64) error {
	return c.fundTransferVerbose(profileID, transferID, false)
}

func (c *scaClient) fundTransferVerbose(profileID, transferID int64, verbose bool) error {
	url := fmt.Sprintf("%s/v3/profiles/%d/transfers/%d/payments", c.baseURL, profileID, transferID)

	payload := map[string]string{"type": "BALANCE"}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Initial request
	if verbose {
		fmt.Printf("[SCA] Initial request to %s\n", url)
	}
	resp, ott, err := c.makeFundRequest(url, body, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Printf("[SCA] Response: %d\n", resp.StatusCode)
		fmt.Printf("[SCA] Headers: %v\n", resp.Header)
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("[SCA] Body: %s\n", string(respBody))
		// Reset body for further processing
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	// Handle SCA challenge
	if resp.StatusCode == http.StatusForbidden && ott != "" {
		if verbose {
			fmt.Printf("[SCA] Got OTT: %s\n", ott)
			// Get OTT status to see available challenges
			status, err := c.getOTTStatus(ott)
			if err == nil {
				fmt.Printf("[SCA] OTT Status: %s\n", string(status))
			}
		}

		ottSignature, err := signData(c.privateKey, []byte(ott))
		if err != nil {
			return fmt.Errorf("signing OTT: %w", err)
		}

		if verbose {
			fmt.Printf("[SCA] Signature: %s\n", ottSignature)
			fmt.Printf("[SCA] Retrying with x-2fa-approval=%s, X-Signature=%s...\n", ott, ottSignature[:40]+"...")
		}

		resp, _, err = c.makeFundRequestWithSCA(url, body, ott, ottSignature)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if verbose {
			fmt.Printf("[SCA] Retry response: %d\n", resp.StatusCode)
			fmt.Printf("[SCA] Retry headers: %v\n", resp.Header)
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *scaClient) makeFundRequest(url string, body []byte, ottSignature string) (*http.Response, string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if ottSignature != "" {
		req.Header.Set("x-2fa-Approval", ottSignature)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	ott := resp.Header.Get("x-2fa-Approval")
	if ott == "" {
		ott = resp.Header.Get("X-2fa-Approval")
	}

	return resp, ott, nil
}

func (c *scaClient) makeFundRequestWithSCA(url string, body []byte, ott, ottSignature string) (*http.Response, string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-2FA-Approval", ott)
	req.Header.Set("X-Signature", ottSignature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	newOTT := resp.Header.Get("x-2fa-Approval")
	if newOTT == "" {
		newOTT = resp.Header.Get("X-2fa-Approval")
	}

	return resp, newOTT, nil
}

// getOTTStatus retrieves the status of a one-time token (for debugging).
func (c *scaClient) getOTTStatus(ott string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/v1/one-time-token/status", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("One-Time-Token", ott)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// triggerSMS triggers an SMS challenge for the given OTT.
func (c *scaClient) triggerSMS(ott string) error {
	req, err := http.NewRequest("POST", c.baseURL+"/v1/one-time-token/sms/trigger", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("One-Time-Token", ott)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("SMS trigger failed %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// verifySMS verifies an SMS OTP code. In sandbox, the code is always "111111".
func (c *scaClient) verifySMS(ott, otpCode string) error {
	payload := map[string]string{"otpCode": otpCode}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", c.baseURL+"/v1/one-time-token/sms/verify", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("One-Time-Token", ott)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("SMS verify failed %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// fundTransferWithSMS funds a transfer using SMS verification (sandbox: OTP is "111111").
func (c *scaClient) fundTransferWithSMS(profileID, transferID int64, verbose bool) error {
	url := fmt.Sprintf("%s/v3/profiles/%d/transfers/%d/payments", c.baseURL, profileID, transferID)

	payload := map[string]string{"type": "BALANCE"}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Initial request to get OTT
	if verbose {
		fmt.Printf("[SMS] Initial request to %s\n", url)
	}
	resp, ott, err := c.makeFundRequest(url, body, "")
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden || ott == "" {
		return fmt.Errorf("expected 403 with OTT, got %d", resp.StatusCode)
	}

	if verbose {
		fmt.Printf("[SMS] Got OTT: %s\n", ott)
		status, _ := c.getOTTStatus(ott)
		fmt.Printf("[SMS] OTT Status: %s\n", string(status))
	}

	// Trigger SMS
	if verbose {
		fmt.Println("[SMS] Triggering SMS challenge...")
	}
	err = c.triggerSMS(ott)
	if err != nil {
		return fmt.Errorf("triggering SMS: %w", err)
	}

	// Verify with sandbox OTP
	if verbose {
		fmt.Println("[SMS] Verifying with OTP 111111...")
	}
	err = c.verifySMS(ott, "111111")
	if err != nil {
		return fmt.Errorf("verifying SMS: %w", err)
	}

	// Retry the funding request with cleared OTT
	if verbose {
		fmt.Println("[SMS] Retrying fund request with cleared OTT...")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-2FA-Approval", ott)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Printf("[SMS] Final response: %d\n", resp.StatusCode)
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("[SMS] Body: %s\n", string(respBody))
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
