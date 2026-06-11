# Send a Trust Line Token Examples (Go)

Send a trust line token on the XRP Ledger Testnet from an issuer to a holder. Creates a trust line, sends a Payment of the issued currency, and verifies balances.

## Setup

```sh
go mod tidy
```

## Run the example

First, run the setup script to fund an issuer and holder wallet and enable Default Ripple on the issuer:

```sh
go run ./setup
```

Then run the main tutorial script:

```sh
go run ./send-trust-line-token
```

The main script auto-invokes the setup script if `setup.json` is not present. It creates a trust line from the holder, sends 500 FOO from issuer to holder, and prints the holder's trust lines and the issuer's outstanding obligations.

```sh
Issuer address: rEXAMPLEISSUERaddressFromTheTestnetFaucet
Holder address: rEXAMPLEHOLDERaddressFromTheTestnetFaucet

=== Preparing TrustSet transaction ===

{
  "Account": "rEXAMPLEHOLDERaddressFromTheTestnetFaucet",
  "Fee": "12",
  "LastLedgerSequence": 18147680,
  "LimitAmount": {
    "currency": "FOO",
    "issuer": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
    "value": "1000000000"
  },
  "Sequence": 18147660,
  "TransactionType": "TrustSet"
}

=== Submitting TrustSet transaction ===

Trust line created from holder to issuer.

=== Preparing Payment transaction ===

{
  "Account": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
  "Amount": {
    "currency": "FOO",
    "issuer": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
    "value": "500"
  },
  "Destination": "rEXAMPLEHOLDERaddressFromTheTestnetFaucet",
  "Fee": "12",
  "LastLedgerSequence": 18147682,
  "Sequence": 18147660,
  "TransactionType": "Payment"
}

=== Submitting Payment transaction ===

Sent 500 FOO from issuer to holder.

=== Verifying balances ===

Holder trust lines:
[
  {
    "account": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
    "balance": "500",
    "currency": "FOO",
    "limit": "1000000000",
    "limit_peer": "0",
    "quality_in": 0,
    "quality_out": 0
  }
]

Issuer obligations:
{
  "FOO": "500"
}
```
