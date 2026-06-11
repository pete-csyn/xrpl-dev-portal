// IMPORTANT: This example sends a trust line token on the XRP Ledger Testnet
// from a preconfigured issuer to a holder. It creates a trust line from the
// holder to the issuer, sends a Payment of the issued currency, and verifies
// the resulting balances. The token sent here is "FOO".

import fs from 'fs'
import { execSync } from 'child_process'
import xrpl from 'xrpl'

// Set up client ----------------------
const client = new xrpl.Client('wss://s.altnet.rippletest.net:51233')
await client.connect()

// This step checks for the necessary setup data to run the tutorial.
// If missing, setup.js will fund the issuer and holder wallets and
// enable Default Ripple on the issuer.
if (!fs.existsSync('setup.json')) {
  console.log(`\n=== Tutorial setup data doesn't exist. Running setup script... ===\n`)
  execSync('node setup.js', { stdio: 'inherit' })
}

// Load preconfigured issuer and holder accounts.
const setupData = JSON.parse(fs.readFileSync('setup.json', 'utf8'))
const issuer = xrpl.Wallet.fromSeed(setupData.issuer.seed)
const holder = xrpl.Wallet.fromSeed(setupData.holder.seed)

console.log(`Issuer address: ${issuer.address}`)
console.log(`Holder address: ${holder.address}`)

const currencyCode = 'FOO'

// Create trust line ----------------------
// The holder authorizes the issuer to send up to the limit_amount of the
// issued currency to the holder's account.
console.log('\n=== Preparing TrustSet transaction ===\n')
const trustSetTx = {
  TransactionType: 'TrustSet',
  Account: holder.address,
  LimitAmount: {
    currency: currencyCode,
    issuer: issuer.address,
    value: '1000000000'
  }
}

xrpl.validate(trustSetTx)
console.log(JSON.stringify(trustSetTx, null, 2))

console.log('\n=== Submitting TrustSet transaction ===\n')
const trustSetResponse = await client.submitAndWait(trustSetTx, {
  wallet: holder,
  autofill: true
})
if (trustSetResponse.result.meta.TransactionResult !== 'tesSUCCESS') {
  const resultCode = trustSetResponse.result.meta.TransactionResult
  console.error('Error: Unable to create trust line:', resultCode)
  await client.disconnect()
  process.exit(1)
}
console.log('Trust line created from holder to issuer.')

// Send issued token ----------------------
// The issuer sends a Payment of the issued currency to the holder. The
// Amount object specifies the currency, issuer, and value.
const issueQuantity = '500'
console.log('\n=== Preparing Payment transaction ===\n')
const paymentTx = {
  TransactionType: 'Payment',
  Account: issuer.address,
  Destination: holder.address,
  Amount: {
    currency: currencyCode,
    issuer: issuer.address,
    value: issueQuantity
  }
}

xrpl.validate(paymentTx)
console.log(JSON.stringify(paymentTx, null, 2))

console.log('\n=== Submitting Payment transaction ===\n')
const paymentResponse = await client.submitAndWait(paymentTx, {
  wallet: issuer,
  autofill: true
})
if (paymentResponse.result.meta.TransactionResult !== 'tesSUCCESS') {
  const resultCode = paymentResponse.result.meta.TransactionResult
  console.error('Error: Unable to send issued token:', resultCode)
  await client.disconnect()
  process.exit(1)
}
console.log(`Sent ${issueQuantity} ${currencyCode} from issuer to holder.`)

// Verify balances ----------------------
// Confirm the holder received the token by reading its trust line balances,
// and confirm the issuer's outstanding obligation through gateway_balances.
console.log('\n=== Verifying balances ===\n')
const holderLines = await client.request({
  command: 'account_lines',
  account: holder.address,
  ledger_index: 'validated'
})
console.log('Holder trust lines:')
console.log(JSON.stringify(holderLines.result.lines, null, 2))

const issuerObligations = await client.request({
  command: 'gateway_balances',
  account: issuer.address,
  ledger_index: 'validated'
})
console.log('\nIssuer obligations:')
console.log(JSON.stringify(issuerObligations.result.obligations, null, 2))

await client.disconnect()
