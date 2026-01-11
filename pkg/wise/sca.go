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

// SCAClient extends Client with Strong Customer Authentication support.
// Note: SCA is currently not working due to signature rejection by Wise.
// This code is preserved for when SCA is enabled on the account.
type SCAClient struct {
	*Client
	privateKey *rsa.PrivateKey
}

// NewSCAClient creates a client with SCA support using a private key file.
func NewSCAClient(apiKey, privateKeyPath string) (*SCAClient, error) {
	privateKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &SCAClient{
		Client:     NewClient(apiKey),
		privateKey: privateKey,
	}, nil
}

// NewSandboxSCAClient creates a sandbox client with SCA support.
func NewSandboxSCAClient(apiKey, privateKeyPath string) (*SCAClient, error) {
	privateKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &SCAClient{
		Client:     NewSandboxClient(apiKey),
		privateKey: privateKey,
	}, nil
}

// LoadPrivateKey loads an RSA private key from a PEM file.
func LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
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

// SignData signs data with SHA256 and RSA PKCS1v15, returning base64.
func SignData(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("signing: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// FundTransfer funds a transfer from the account balance.
func (c *SCAClient) FundTransfer(profileID, transferID int64) error {
	return c.FundTransferVerbose(profileID, transferID, false)
}

// FundTransferVerbose funds a transfer with optional debug output.
func (c *SCAClient) FundTransferVerbose(profileID, transferID int64, verbose bool) error {
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
			status, err := c.GetOTTStatus(ott)
			if err == nil {
				fmt.Printf("[SCA] OTT Status: %s\n", string(status))
			}
		}

		ottSignature, err := SignData(c.privateKey, []byte(ott))
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

func (c *SCAClient) makeFundRequest(url string, body []byte, ottSignature string) (*http.Response, string, error) {
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

func (c *SCAClient) makeFundRequestWithSCA(url string, body []byte, ott, ottSignature string) (*http.Response, string, error) {
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

// GetOTTStatus retrieves the status of a one-time token (for debugging).
func (c *SCAClient) GetOTTStatus(ott string) ([]byte, error) {
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
