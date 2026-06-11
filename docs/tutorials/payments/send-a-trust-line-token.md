---
seo:
  description: Send a trust line token from an issuer account to a holder account on the XRP Ledger.
metadata:
  indexPage: true
labels:
  - Tokens
  - Payments
---

# Send a Trust Line Token

This tutorial shows you how to send a fungible [trust line token](../../concepts/tokens/fungible-tokens/trust-line-tokens.md) from an issuer account to a holder account on the XRP Ledger. The example creates a trust line from the holder to the issuer, sends the token in a [Payment transaction][], and verifies that the holder received the balance.

## Goals

By the end of this tutorial, you will be able to:

- Create a trust line from a holder account to an issuer using a [TrustSet transaction][].
- Send a trust line token from the issuer to the holder using a [Payment transaction][].
- Verify the holder's balance and the issuer's outstanding obligations on the ledger.

## Prerequisites

To complete this tutorial, you should:

- Have a basic understanding of the XRP Ledger.
- Have an XRP Ledger client library set up in your development environment. This page provides examples for the following:
  - **JavaScript** with the [xrpl.js library][]. See [Get Started Using JavaScript][] for setup steps.
  - **Python** with the [xrpl-py library][]. See [Get Started Using Python][] for setup steps.
  - **Go** with the [xrpl-go library][]. See [Get Started Using Go][] for setup steps.

## Source Code

You can find the complete source code for this tutorial's examples in the {% repo-link path="_code-samples/send-a-trust-line-token/" %}code samples section of this website's repository{% /repo-link %}.

## Steps

### 1. Install dependencies

{% tabs %}
{% tab label="JavaScript" %}
From the code sample folder, use `npm` to install dependencies.

```bash
npm install
```
{% /tab %}
{% tab label="Python" %}
From the code sample folder, set up a virtual environment and use `pip` to install dependencies.

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```
{% /tab %}
{% tab label="Go" %}
From the code sample folder, use `go` to install dependencies.

```bash
go mod tidy
```
{% /tab %}
{% /tabs %}

### 2. Set up client and accounts

To get started, import the necessary libraries and instantiate a client to connect to the XRPL. This example imports:

{% tabs %}
{% tab label="JavaScript" %}
- `xrpl`: Used for XRPL client connection, transaction submission, and wallet handling.
- `fs` and `child_process`: Used to run the tutorial setup script.

{% code-snippet file="/_code-samples/send-a-trust-line-token/js/sendTrustLineToken.js" language="js" before="// This step checks" /%}
{% /tab %}
{% tab label="Python" %}
- `xrpl`: Used for XRPL client connection, transaction submission, and wallet handling.
- `json`: Used for loading and formatting JSON data.
- `os`, `subprocess`, and `sys`: Used to run the tutorial setup script.

{% code-snippet file="/_code-samples/send-a-trust-line-token/py/send_trust_line_token.py" language="py" before="# This step checks" /%}
{% /tab %}
{% tab label="Go" %}
- `xrpl-go`: Used for XRPL client connection, transaction submission, and wallet handling.
- `encoding/json` and `fmt`: Used for formatting and printing results to the console.
- `os` and `os/exec`: Used to run the tutorial setup script.

{% code-snippet file="/_code-samples/send-a-trust-line-token/go/send-trust-line-token/main.go" language="go" before="// This step checks" /%}
{% /tab %}
{% /tabs %}

Next, load the preconfigured issuer and holder accounts. The setup script funds two Testnet wallets and enables the `Default Ripple` flag on the issuer with an [AccountSet transaction][], so holders of the issued token can transfer it between each other through the issuer's trust lines. The main script auto-invokes the setup script if `setup.json` is not present.

{% tabs %}
{% tab label="JavaScript" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/js/sendTrustLineToken.js" language="js" from="// This step checks" before="// Create trust line" /%}

This example uses preconfigured accounts from the `setup.js` script, but you can replace `issuer` and `holder` with your own values.
{% /tab %}
{% tab label="Python" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/py/send_trust_line_token.py" language="py" from="# This step checks" before="# Create trust line" /%}

This example uses preconfigured accounts from the `setup.py` script, but you can replace `issuer` and `holder` with your own values.
{% /tab %}
{% tab label="Go" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/go/send-trust-line-token/main.go" language="go" from="// This step checks" before="// Create trust line" /%}

This example uses preconfigured accounts from the `setup` script, but you can replace `issuer` and `holder` with your own values.
{% /tab %}
{% /tabs %}

### 3. Create a trust line from the holder to the issuer

The holder must opt in to receive the token by creating a trust line to the issuer with a [TrustSet transaction][]. The `LimitAmount` field specifies the maximum quantity of the token the holder is willing to hold, the currency code, and the issuer's address.

{% tabs %}
{% tab label="JavaScript" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/js/sendTrustLineToken.js" language="js" from="// Create trust line" before="// Send issued token" /%}
{% /tab %}
{% tab label="Python" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/py/send_trust_line_token.py" language="py" from="# Create trust line" before="# Send issued token" /%}
{% /tab %}
{% tab label="Go" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/go/send-trust-line-token/main.go" language="go" from="// Create trust line" before="// Send issued token" /%}
{% /tab %}
{% /tabs %}

### 4. Send the token to the holder

Once the trust line exists, the issuer can send the token to the holder with a [Payment transaction][]. The `Amount` field is an [issued currency amount](../../references/protocol/data-types/currency-formats.md#token-amounts) object that specifies the currency code, the issuer's address, and the value to send.

{% tabs %}
{% tab label="JavaScript" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/js/sendTrustLineToken.js" language="js" from="// Send issued token" before="// Verify balances" /%}
{% /tab %}
{% tab label="Python" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/py/send_trust_line_token.py" language="py" from="# Send issued token" before="# Verify balances" /%}
{% /tab %}
{% tab label="Go" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/go/send-trust-line-token/main.go" language="go" from="// Send issued token" before="// Verify balances" /%}
{% /tab %}
{% /tabs %}

### 5. Verify the balances

Confirm the token transfer by reading the holder's trust lines with the [account_lines method][] and the issuer's outstanding obligations with the [gateway_balances method][]. The holder's balance should match the amount sent, and the issuer's obligations should show the same value as a negative balance on the issuer's side of the trust line.

{% tabs %}
{% tab label="JavaScript" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/js/sendTrustLineToken.js" language="js" from="// Verify balances" /%}
{% /tab %}
{% tab label="Python" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/py/send_trust_line_token.py" language="py" from="# Verify balances" /%}
{% /tab %}
{% tab label="Go" %}
{% code-snippet file="/_code-samples/send-a-trust-line-token/go/send-trust-line-token/main.go" language="go" from="// Verify balances" before="func printJSON" /%}
{% /tab %}
{% /tabs %}

## See Also

- **Concepts**:
  - [Trust Line Tokens](../../concepts/tokens/fungible-tokens/trust-line-tokens.md)
  - [Rippling](../../concepts/tokens/fungible-tokens/rippling.md)
  - [Issuing and Operational Addresses](../../concepts/accounts/account-types.md)
- **Tutorials**:
  - [Issue a Fungible Token](../tokens/fungible-tokens/issue-a-fungible-token.md)
  - [Send XRP](./send-xrp.md)
- **References**:
  - [AccountSet transaction][]
  - [TrustSet transaction][]
  - [Payment transaction][]
  - [account_lines method][]
  - [gateway_balances method][]

{% raw-partial file="/docs/_snippets/common-links.md" /%}
