# IMPORTANT: This example sends a trust line token on the XRP Ledger Testnet
# from a preconfigured issuer to a holder. It creates a trust line from the
# holder to the issuer, sends a Payment of the issued currency, and verifies
# the resulting balances. The token sent here is "FOO".

import json
import os
import subprocess
import sys

from xrpl.clients import JsonRpcClient
from xrpl.models import AccountLines, GatewayBalances
from xrpl.models.amounts import IssuedCurrencyAmount
from xrpl.models.transactions import Payment, TrustSet
from xrpl.transaction import submit_and_wait
from xrpl.wallet import Wallet

# Set up client ----------------------
client = JsonRpcClient("https://s.altnet.rippletest.net:51234")

# This step checks for the necessary setup data to run the tutorial.
# If missing, setup.py will fund the issuer and holder wallets and
# enable Default Ripple on the issuer.
if not os.path.exists("setup.json"):
    print("\n=== Tutorial setup data doesn't exist. Running setup script... ===\n")
    subprocess.run([sys.executable, "setup.py"], check=True)

# Load preconfigured issuer and holder accounts.
with open("setup.json") as f:
    setup_data = json.load(f)

issuer = Wallet.from_seed(setup_data["issuer"]["seed"])
holder = Wallet.from_seed(setup_data["holder"]["seed"])

print(f"Issuer address: {issuer.address}")
print(f"Holder address: {holder.address}")

currency_code = "FOO"

# Create trust line ----------------------
# The holder authorizes the issuer to send up to the limit_amount of the
# issued currency to the holder's account.
print("\n=== Preparing TrustSet transaction ===\n")
trust_set_tx = TrustSet(
    account=holder.address,
    limit_amount=IssuedCurrencyAmount(
        currency=currency_code,
        issuer=issuer.address,
        value="1000000000",
    ),
)

print(json.dumps(trust_set_tx.to_xrpl(), indent=2))

print("\n=== Submitting TrustSet transaction ===\n")
trust_set_response = submit_and_wait(trust_set_tx, client, holder)
result_code = trust_set_response.result["meta"]["TransactionResult"]
if result_code != "tesSUCCESS":
    print(f"Error: Unable to create trust line: {result_code}")
    sys.exit(1)
print("Trust line created from holder to issuer.")

# Send issued token ----------------------
# The issuer sends a Payment of the issued currency to the holder. The
# amount object specifies the currency, issuer, and value.
issue_quantity = "500"
print("\n=== Preparing Payment transaction ===\n")
payment_tx = Payment(
    account=issuer.address,
    destination=holder.address,
    amount=IssuedCurrencyAmount(
        currency=currency_code,
        issuer=issuer.address,
        value=issue_quantity,
    ),
)

print(json.dumps(payment_tx.to_xrpl(), indent=2))

print("\n=== Submitting Payment transaction ===\n")
payment_response = submit_and_wait(payment_tx, client, issuer)
result_code = payment_response.result["meta"]["TransactionResult"]
if result_code != "tesSUCCESS":
    print(f"Error: Unable to send issued token: {result_code}")
    sys.exit(1)
print(f"Sent {issue_quantity} {currency_code} from issuer to holder.")

# Verify balances ----------------------
# Confirm the holder received the token by reading its trust line balances,
# and confirm the issuer's outstanding obligation through gateway_balances.
print("\n=== Verifying balances ===\n")
holder_lines = client.request(AccountLines(
    account=holder.address,
    ledger_index="validated",
))
print("Holder trust lines:")
print(json.dumps(holder_lines.result["lines"], indent=2))

issuer_obligations = client.request(GatewayBalances(
    account=issuer.address,
    ledger_index="validated",
))
print("\nIssuer obligations:")
print(json.dumps(issuer_obligations.result.get("obligations", {}), indent=2))
