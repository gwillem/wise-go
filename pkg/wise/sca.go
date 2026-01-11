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
// This requires SCA which is currently not working.
func (c *SCAClient) FundTransfer(profileID, transferID int64) error {
	url := fmt.Sprintf("%s/v3/profiles/%d/transfers/%d/payments", c.baseURL, profileID, transferID)

	payload := map[string]string{"type": "BALANCE"}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Initial request
	resp, ott, err := c.makeFundRequest(url, body, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle SCA challenge
	if resp.StatusCode == http.StatusForbidden && ott != "" {
		ottSignature, err := SignData(c.privateKey, []byte(ott))
		if err != nil {
			return fmt.Errorf("signing OTT: %w", err)
		}

		resp, _, err = c.makeFundRequestWithSCA(url, body, ott, ottSignature)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
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
	req.Header.Set("x-2fa-approval", ott)
	req.Header.Set("x-signature", ottSignature)

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
