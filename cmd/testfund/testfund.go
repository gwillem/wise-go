// testfund tests SCA funding in the Wise sandbox.
package main

import (
	"fmt"
	"os"

	"wise/pkg/wise"
)

const (
	sandboxKey     = "cbbc7006-576c-48b3-92db-2bedc9c67fda"
	privateKeyPath = "private.pem"
	profileID      = int64(29137484) // sandbox business profile
)

func main() {
	fmt.Println("=== Wise SCA Funding Test ===")
	fmt.Println()

	// Create SCA client
	client, err := wise.NewSandboxSCAClient(sandboxKey, privateKeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	// Create a transfer
	fmt.Println("Creating transfer...")
	result, err := client.CreateTransfer(wise.TransferRequest{
		ProfileID:      profileID,
		SourceCurrency: "EUR",
		TargetCurrency: "EUR",
		Amount:         1.00,
		RecipientName:  "Test Recipient",
		IBAN:           "NL44BUNQ2070060462",
		Reference:      "SCA Test",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating transfer: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Transfer created: ID=%d, Status=%s\n", result.Transfer.ID, result.Transfer.Status)
	fmt.Printf("Quote: %.2f %s -> %.2f %s (fee: %.2f)\n",
		result.Quote.SourceAmount, result.Quote.SourceCurrency,
		result.Quote.TargetAmount, result.Quote.TargetCurrency,
		result.Quote.Fee)
	fmt.Println()

	// Fund the transfer - first get the OTT to see available challenges
	fmt.Println("Attempting to fund transfer...")
	err = client.FundTransferVerbose(profileID, result.Transfer.ID, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error funding transfer: %v\n", err)
		fmt.Println("\nNote: The SIGNATURE challenge type may not be enabled for this account.")
		fmt.Println("Contact Wise to enable SCA signature-based authentication.")
		os.Exit(1)
	}

	fmt.Println("Transfer funded successfully!")
}
