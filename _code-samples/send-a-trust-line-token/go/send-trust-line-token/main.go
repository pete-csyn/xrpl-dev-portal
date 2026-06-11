// IMPORTANT: This example sends a trust line token on the XRP Ledger Testnet
// from a preconfigured issuer to a holder. It creates a trust line from the
// holder to the issuer, sends a Payment of the issued currency, and verifies
// the resulting balances. The token sent here is "FOO".

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/Peersyst/xrpl-go/xrpl/queries/account"
	accountv1 "github.com/Peersyst/xrpl-go/xrpl/queries/account/v1"
	"github.com/Peersyst/xrpl-go/xrpl/queries/common"
	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	transactions "github.com/Peersyst/xrpl-go/xrpl/transaction"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
	"github.com/Peersyst/xrpl-go/xrpl/wallet"
)

const currencyCode = "FOO"

func main() {
	// Set up client ----------------------
	cfg, err := rpc.NewClientConfig("https://s.altnet.rippletest.net:51234/")
	if err != nil {
		panic(err)
	}
	client := rpc.NewClient(cfg)

	// This step checks for the necessary setup data to run the tutorial.
	// If missing, the setup script funds the issuer and holder wallets
	// and enables Default Ripple on the issuer.
	if _, err := os.Stat("setup.json"); os.IsNotExist(err) {
		fmt.Printf("\n=== Tutorial setup data doesn't exist. Running setup script... ===\n\n")
		cmd := exec.Command("go", "run", "./setup")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	// Load preconfigured issuer and holder accounts.
	data, err := os.ReadFile("setup.json")
	if err != nil {
		panic(err)
	}
	var setup map[string]any
	if err := json.Unmarshal(data, &setup); err != nil {
		panic(err)
	}
	issuer, err := wallet.FromSecret(setup["issuer"].(map[string]any)["seed"].(string))
	if err != nil {
		panic(err)
	}
	holder, err := wallet.FromSecret(setup["holder"].(map[string]any)["seed"].(string))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Issuer address: %s\n", issuer.ClassicAddress)
	fmt.Printf("Holder address: %s\n", holder.ClassicAddress)

	// Create trust line ----------------------
	// The holder authorizes the issuer to send up to the limit_amount of the
	// issued currency to the holder's account.
	fmt.Println("\n=== Preparing TrustSet transaction ===\n")
	trustSetTx := &transactions.TrustSet{
		BaseTx: transactions.BaseTx{
			Account: types.Address(holder.ClassicAddress),
		},
		LimitAmount: types.IssuedCurrencyAmount{
			Currency: currencyCode,
			Issuer:   types.Address(issuer.ClassicAddress),
			Value:    "1000000000",
		},
	}

	flatTx := trustSetTx.Flatten()
	if err := client.Autofill(&flatTx); err != nil {
		panic(err)
	}
	printJSON(flatTx)

	fmt.Println("\n=== Submitting TrustSet transaction ===\n")
	txBlob, _, err := holder.Sign(flatTx)
	if err != nil {
		panic(err)
	}
	trustSetResp, err := client.SubmitTxBlobAndWait(txBlob, false)
	if err != nil {
		panic(err)
	}
	if !trustSetResp.Validated {
		fmt.Println("Error: Unable to create trust line.")
		return
	}
	fmt.Println("Trust line created from holder to issuer.")

	// Send issued token ----------------------
	// The issuer sends a Payment of the issued currency to the holder. The
	// Amount object specifies the currency, issuer, and value.
	const issueQuantity = "500"
	fmt.Println("\n=== Preparing Payment transaction ===\n")
	paymentTx := &transactions.Payment{
		BaseTx: transactions.BaseTx{
			Account: types.Address(issuer.ClassicAddress),
		},
		Destination: types.Address(holder.ClassicAddress),
		Amount: types.IssuedCurrencyAmount{
			Currency: currencyCode,
			Issuer:   types.Address(issuer.ClassicAddress),
			Value:    issueQuantity,
		},
	}

	flatTx = paymentTx.Flatten()
	if err := client.Autofill(&flatTx); err != nil {
		panic(err)
	}
	printJSON(flatTx)

	fmt.Println("\n=== Submitting Payment transaction ===\n")
	txBlob, _, err = issuer.Sign(flatTx)
	if err != nil {
		panic(err)
	}
	paymentResp, err := client.SubmitTxBlobAndWait(txBlob, false)
	if err != nil {
		panic(err)
	}
	if !paymentResp.Validated {
		fmt.Println("Error: Unable to send issued token.")
		return
	}
	fmt.Printf("Sent %s %s from issuer to holder.\n", issueQuantity, currencyCode)

	// Verify balances ----------------------
	// Confirm the holder received the token by reading its trust line balances,
	// and confirm the issuer's outstanding obligation through gateway_balances.
	fmt.Println("\n=== Verifying balances ===\n")
	linesResp, err := client.Request(&accountv1.LinesRequest{
		Account:     types.Address(holder.ClassicAddress),
		LedgerIndex: common.Validated,
	})
	if err != nil {
		panic(err)
	}
	var linesResult accountv1.LinesResponse
	if err := linesResp.GetResult(&linesResult); err != nil {
		panic(err)
	}
	fmt.Println("Holder trust lines:")
	printJSON(linesResult.Lines)

	balancesResp, err := client.Request(&account.GatewayBalancesRequest{
		Account:     types.Address(issuer.ClassicAddress),
		LedgerIndex: common.Validated,
	})
	if err != nil {
		panic(err)
	}
	var balancesResult account.GatewayBalancesResponse
	if err := balancesResp.GetResult(&balancesResult); err != nil {
		panic(err)
	}
	fmt.Println("\nIssuer obligations:")
	printJSON(balancesResult.Obligations)
}

func printJSON(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
