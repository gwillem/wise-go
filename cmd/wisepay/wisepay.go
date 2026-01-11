// wisepay is a CLI tool for creating Wise transfers.
package main

import (
	"flag"
	"fmt"
	"os"

	wise "github.com/gwillem/wise-go"
)

func main() {
	var (
		sandbox       = flag.Bool("sandbox", false, "Use sandbox API")
		profileID     = flag.Int64("profile", 0, "Profile ID (required)")
		amount        = flag.Float64("amount", 0, "Amount to send (required)")
		sourceCurrency = flag.String("source", "EUR", "Source currency")
		targetCurrency = flag.String("target", "EUR", "Target currency")
		iban          = flag.String("iban", "", "Recipient IBAN (required)")
		name          = flag.String("name", "", "Recipient name (required)")
		reference     = flag.String("ref", "Payment", "Transfer reference")
		listProfiles  = flag.Bool("list", false, "List profiles and balances")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: wisepay [options]\n\n")
		fmt.Fprintf(os.Stderr, "Environment:\n")
		fmt.Fprintf(os.Stderr, "  WISE_KEY    Wise API token (required)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  wisepay -list\n")
		fmt.Fprintf(os.Stderr, "  wisepay -profile 123 -amount 100 -iban DE89370400440532013000 -name \"John Doe\"\n")
	}

	flag.Parse()

	apiKey := os.Getenv("WISE_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: WISE_KEY environment variable is required")
		os.Exit(1)
	}

	var client *wise.Client
	if *sandbox {
		client = wise.NewSandboxClient(apiKey)
	} else {
		client = wise.NewClient(apiKey)
	}

	if *listProfiles {
		if err := listProfilesAndBalances(client); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Validate required fields for transfer
	if *profileID == 0 {
		fmt.Fprintln(os.Stderr, "Error: -profile is required")
		flag.Usage()
		os.Exit(1)
	}
	if *amount <= 0 {
		fmt.Fprintln(os.Stderr, "Error: -amount must be positive")
		flag.Usage()
		os.Exit(1)
	}
	if *iban == "" {
		fmt.Fprintln(os.Stderr, "Error: -iban is required")
		flag.Usage()
		os.Exit(1)
	}
	if *name == "" {
		fmt.Fprintln(os.Stderr, "Error: -name is required")
		flag.Usage()
		os.Exit(1)
	}

	result, err := client.CreateTransfer(wise.TransferRequest{
		ProfileID:      *profileID,
		SourceCurrency: *sourceCurrency,
		TargetCurrency: *targetCurrency,
		Amount:         *amount,
		RecipientName:  *name,
		IBAN:           *iban,
		Reference:      *reference,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Transfer created successfully!\n\n")
	fmt.Printf("Quote:     %s (rate: %.4f, fee: %.2f %s)\n",
		result.Quote.ID, result.Quote.Rate, result.Quote.Fee, result.Quote.SourceCurrency)
	fmt.Printf("Recipient: %d (%s)\n", result.Recipient.ID, result.Recipient.AccountHolderName)
	fmt.Printf("Transfer:  %d (status: %s)\n", result.Transfer.ID, result.Transfer.Status)
	fmt.Printf("Amount:    %.2f %s -> %.2f %s\n",
		result.Quote.SourceAmount, result.Quote.SourceCurrency,
		result.Quote.TargetAmount, result.Quote.TargetCurrency)
	fmt.Printf("\nApprove the transfer in the Wise app to complete.\n")
}

func listProfilesAndBalances(client *wise.Client) error {
	profiles, err := client.GetProfiles()
	if err != nil {
		return fmt.Errorf("getting profiles: %w", err)
	}

	fmt.Println("Profiles and Balances")
	fmt.Println("=====================")

	for _, profile := range profiles {
		fmt.Printf("\n%s (ID: %d, Type: %s)\n", profile.FullName, profile.ID, profile.Type)

		balances, err := client.GetBalances(profile.ID)
		if err != nil {
			fmt.Printf("  Error getting balances: %v\n", err)
			continue
		}

		if len(balances) == 0 {
			fmt.Println("  No balances")
		}
		for _, b := range balances {
			fmt.Printf("  %s: %.2f\n", b.Amount.Currency, b.Amount.Value)
		}
	}

	return nil
}
