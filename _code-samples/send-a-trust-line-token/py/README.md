# Send a Trust Line Token Examples (Python)

Send a trust line token on the XRP Ledger Testnet from an issuer to a holder. Creates a trust line, sends a Payment of the issued currency, and verifies balances.

## Setup

```sh
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## Run the example

First, run the setup script to fund an issuer and holder wallet and enable Default Ripple on the issuer:

```sh
python setup.py
```

Then run the main tutorial script:

```sh
python send_trust_line_token.py
```

The main script auto-invokes `setup.py` if `setup.json` is not present. It creates a trust line from the holder, sends 500 FOO from issuer to holder, and prints the holder's trust lines and the issuer's outstanding obligations.

```sh
Issuer address: rEXAMPLEISSUERaddressFromTheTestnetFaucet
Holder address: rEXAMPLEHOLDERaddressFromTheTestnetFaucet

=== Preparing TrustSet transaction ===

{
  "TransactionType": "TrustSet",
  "Account": "rEXAMPLEHOLDERaddressFromTheTestnetFaucet",
  "LimitAmount": {
    "currency": "FOO",
    "issuer": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
    "value": "1000000000"
  }
}

=== Submitting TrustSet transaction ===

Trust line created from holder to issuer.

=== Preparing Payment transaction ===

{
  "TransactionType": "Payment",
  "Account": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
  "Destination": "rEXAMPLEHOLDERaddressFromTheTestnetFaucet",
  "Amount": {
    "currency": "FOO",
    "issuer": "rEXAMPLEISSUERaddressFromTheTestnetFaucet",
    "value": "500"
  }
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
