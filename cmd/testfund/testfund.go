// testfund tests transfer funding in the Wise sandbox using simulation.
package main

import (
	"fmt"
	"os"
	"time"

	wise "github.com/gwillem/wise-go"
)

const (
	sandboxKey = "cbbc7006-576c-48b3-92db-2bedc9c67fda"
	profileID  = int64(29137484) // sandbox business profile
)

func main() {
	fmt.Println("=== Wise Sandbox Transfer Test ===")
	fmt.Println()

	// Create client (no SCA needed - using simulation)
	client := wise.NewSandboxClient(sandboxKey)

	// Create a transfer
	fmt.Println("Creating transfer...")
	result, err := client.CreateTransfer(wise.TransferRequest{
		ProfileID:      profileID,
		SourceCurrency: "EUR",
		TargetCurrency: "EUR",
		Amount:         1.00,
		RecipientName:  "Test Recipient",
		IBAN:           "NL44BUNQ2070060462",
		Reference:      "Test Transfer",
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

	// Use simulation to process transfer (bypasses SCA in sandbox)
	fmt.Println("Using simulation to process transfer...")
	simResult, err := client.SimulateTransferProcessing(result.Transfer.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Simulation failed: %v\n", err)
		fmt.Println("\nNote: In production, use SCA to fund the transfer.")
		fmt.Println("The transfer was created and can be approved in the Wise app.")
		os.Exit(0)
	}

	fmt.Printf("Transfer status: %s\n", simResult.Status)

	// Wait and try next steps
	time.Sleep(2 * time.Second)

	simResult, err = client.SimulateTransferFundsConverted(result.Transfer.ID)
	if err != nil {
		fmt.Printf("Funds conversion simulation: %v\n", err)
		fmt.Println("Transfer may need manual verification in sandbox.")
	} else {
		fmt.Printf("Transfer status: %s\n", simResult.Status)

		time.Sleep(2 * time.Second)

		simResult, err = client.SimulateTransferOutgoingPayment(result.Transfer.ID)
		if err != nil {
			fmt.Printf("Outgoing payment simulation: %v\n", err)
		} else {
			fmt.Printf("Transfer status: %s\n", simResult.Status)
		}
	}

	fmt.Println("\nDone!")
}
