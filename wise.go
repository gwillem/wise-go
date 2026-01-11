// Package wise provides a client for the Wise (TransferWise) API.
package wise

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	productionBaseURL = "https://api.transferwise.com"
	sandboxBaseURL    = "https://api.sandbox.transferwise.tech"
)

// Client is a Wise API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Wise API client for production.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL:    productionBaseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// NewSandboxClient creates a new Wise API client for sandbox.
func NewSandboxClient(apiKey string) *Client {
	return &Client{
		baseURL:    sandboxBaseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Profile represents a Wise profile.
type Profile struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"fullName"`
}

// Balance represents an account balance.
type Balance struct {
	ID       int64  `json:"id"`
	Currency string `json:"currency"`
	Amount   struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	} `json:"amount"`
}

// Quote represents a transfer quote.
type Quote struct {
	ID             string  `json:"id"`
	SourceAmount   float64 `json:"sourceAmount"`
	TargetAmount   float64 `json:"targetAmount"`
	Rate           float64 `json:"rate"`
	Fee            float64 `json:"fee"`
	SourceCurrency string  `json:"sourceCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
}

// Recipient represents a transfer recipient.
type Recipient struct {
	ID                int64  `json:"id"`
	AccountHolderName string `json:"accountHolderName"`
	Currency          string `json:"currency"`
}

// Transfer represents a money transfer.
type Transfer struct {
	ID           int64   `json:"id"`
	Status       string  `json:"status"`
	Reference    string  `json:"reference"`
	SourceAmount float64 `json:"sourceAmount"`
	TargetAmount float64 `json:"targetAmount"`
}

// TransferRequest contains the parameters for creating a transfer.
//
// For IBAN-based transfers (European), set IBAN.
// For non-IBAN transfers, set RecipientType and the relevant fields:
//   - "indian": IFSCCode + AccountNumber
//   - "swift_code": SwiftCode + AccountNumber
//   - other local types: AccountNumber (+ SwiftCode if needed)
type TransferRequest struct {
	ProfileID      int64
	SourceCurrency string
	TargetCurrency string
	Amount         float64
	RecipientName  string
	Reference      string

	// IBAN-based (European)
	IBAN string

	// Non-IBAN
	RecipientType string // e.g. "indian", "aba", "swift_code"
	AccountNumber string
	SwiftCode     string
	IFSCCode      string
	LegalType     string // "PRIVATE" or "BUSINESS", defaults to "PRIVATE"
}

// TransferResult contains the result of creating a transfer.
type TransferResult struct {
	Quote     Quote
	Recipient Recipient
	Transfer  Transfer
}

// GetProfiles returns all profiles for the authenticated user.
func (c *Client) GetProfiles() ([]Profile, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/v2/profiles", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var profiles []Profile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

// GetBalances returns all balances for a profile.
func (c *Client) GetBalances(profileID int64) ([]Balance, error) {
	url := fmt.Sprintf("%s/v4/profiles/%d/balances?types=STANDARD", c.baseURL, profileID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var balances []Balance
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, err
	}
	return balances, nil
}

// CreateTransfer creates a new transfer that can be approved in the Wise app.
// Returns the quote, recipient, and transfer details.
func (c *Client) CreateTransfer(req TransferRequest) (*TransferResult, error) {
	// Step 1: Create quote
	quote, err := c.createQuote(req.ProfileID, req.Amount, req.SourceCurrency, req.TargetCurrency)
	if err != nil {
		return nil, fmt.Errorf("creating quote: %w", err)
	}

	// Step 2: Create recipient
	recipient, err := c.createRecipient(req)
	if err != nil {
		return nil, fmt.Errorf("creating recipient: %w", err)
	}

	// Step 3: Create transfer
	transfer, err := c.createTransfer(quote.ID, recipient.ID, req.Reference)
	if err != nil {
		return nil, fmt.Errorf("creating transfer: %w", err)
	}

	return &TransferResult{
		Quote:     *quote,
		Recipient: *recipient,
		Transfer:  *transfer,
	}, nil
}

func (c *Client) createQuote(profileID int64, amount float64, sourceCurrency, targetCurrency string) (*Quote, error) {
	payload := map[string]any{
		"sourceCurrency": sourceCurrency,
		"targetCurrency": targetCurrency,
		"sourceAmount":   amount,
		"profile":        profileID,
		"preferredPayIn": "BALANCE",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v3/profiles/%d/quotes", c.baseURL, profileID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var quote Quote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, err
	}
	return &quote, nil
}

func (c *Client) createRecipient(tr TransferRequest) (*Recipient, error) {
	legalType := tr.LegalType
	if legalType == "" {
		legalType = "PRIVATE"
	}

	recipientType := tr.RecipientType
	details := map[string]any{"legalType": legalType}

	switch {
	case tr.IBAN != "":
		recipientType = "iban"
		details["iban"] = tr.IBAN
	case recipientType == "indian":
		details["ifscCode"] = tr.IFSCCode
		details["accountNumber"] = tr.AccountNumber
	case recipientType == "swift_code":
		details["swiftCode"] = tr.SwiftCode
		details["accountNumber"] = tr.AccountNumber
	case recipientType != "":
		details["accountNumber"] = tr.AccountNumber
		if tr.SwiftCode != "" {
			details["swiftCode"] = tr.SwiftCode
		}
	default:
		return nil, fmt.Errorf("either IBAN or RecipientType must be set")
	}

	payload := map[string]any{
		"currency":          tr.TargetCurrency,
		"type":              recipientType,
		"profile":           tr.ProfileID,
		"ownedByCustomer":   false,
		"accountHolderName": tr.RecipientName,
		"details":           details,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v1/accounts", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var recipient Recipient
	if err := json.NewDecoder(resp.Body).Decode(&recipient); err != nil {
		return nil, err
	}
	return &recipient, nil
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type simulationResult struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func (c *Client) simulateTransferProcessing(transferID int64) (*simulationResult, error) {
	return c.simulateTransfer(transferID, "processing")
}

func (c *Client) simulateTransferFundsConverted(transferID int64) (*simulationResult, error) {
	return c.simulateTransfer(transferID, "funds_converted")
}

func (c *Client) simulateTransferOutgoingPayment(transferID int64) (*simulationResult, error) {
	return c.simulateTransfer(transferID, "outgoing_payment_sent")
}

func (c *Client) simulateTransfer(transferID int64, status string) (*simulationResult, error) {
	url := fmt.Sprintf("%s/v1/simulation/transfers/%d/%s", c.baseURL, transferID, status)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("simulation error %d: %s", resp.StatusCode, string(body))
	}

	var result simulationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) createTransfer(quoteID string, recipientID int64, reference string) (*Transfer, error) {
	payload := map[string]any{
		"targetAccount":         recipientID,
		"quoteUuid":             quoteID,
		"customerTransactionId": generateUUID(),
		"details": map[string]string{
			"reference": reference,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v1/transfers", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var transfer Transfer
	if err := json.NewDecoder(resp.Body).Decode(&transfer); err != nil {
		return nil, err
	}
	return &transfer, nil
}
