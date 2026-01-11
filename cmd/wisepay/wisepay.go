// wisepay is a CLI tool for creating Wise transfers.
package main

import (
	"fmt"
	"os"

	wise "github.com/gwillem/wise-go"
	flags "github.com/jessevdk/go-flags"
)

type opts struct {
	Sandbox        bool    `short:"s" long:"sandbox" description:"Use sandbox API"`
	ProfileID      int64   `short:"p" long:"profile" description:"Profile ID (required for transfers)"`
	Amount         float64 `short:"a" long:"amount" description:"Amount to send"`
	SourceCurrency string  `long:"source" default:"EUR" description:"Source currency"`
	TargetCurrency string  `long:"target" default:"EUR" description:"Target currency"`
	IBAN           string  `long:"iban" description:"Recipient IBAN"`
	Name           string  `short:"n" long:"name" description:"Recipient name"`
	Reference      string  `long:"ref" default:"Payment" description:"Transfer reference"`
	List           bool    `short:"l" long:"list" description:"List profiles and balances"`
	Verbose        bool    `short:"v" long:"verbose" description:"Show verbose output"`
}

func main() {
	var o opts
	_, err := flags.Parse(&o)
	if err != nil {
		os.Exit(1)
	}

	apiKey := os.Getenv("WISE_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: WISE_KEY environment variable is required")
		os.Exit(1)
	}

	var client *wise.Client
	if o.Sandbox {
		client = wise.NewSandboxClient(apiKey)
	} else {
		client = wise.NewClient(apiKey)
	}

	if o.List {
		if err := listProfilesAndBalances(client); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if o.ProfileID == 0 {
		fmt.Fprintln(os.Stderr, "Error: --profile is required")
		os.Exit(1)
	}
	if o.Amount <= 0 {
		fmt.Fprintln(os.Stderr, "Error: --amount must be positive")
		os.Exit(1)
	}
	if o.IBAN == "" {
		fmt.Fprintln(os.Stderr, "Error: --iban is required")
		os.Exit(1)
	}
	if o.Name == "" {
		fmt.Fprintln(os.Stderr, "Error: --name is required")
		os.Exit(1)
	}

	result, err := client.CreateTransfer(wise.TransferRequest{
		ProfileID:      o.ProfileID,
		SourceCurrency: o.SourceCurrency,
		TargetCurrency: o.TargetCurrency,
		Amount:         o.Amount,
		RecipientName:  o.Name,
		IBAN:           o.IBAN,
		Reference:      o.Reference,
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
