<!-- Source: https://docs.wise.com/api-reference/payin-deposit-detail -->

# Payin Deposit Detail
Operations 

[GET/v1/profiles/{{profileId}}/transfers/{{transferId}}/deposit-details/bank-transfer](/api-reference/payin-deposit-detail#get)

Get the bank transfer deposit details for a transfer

## The Payin Deposit Detail resource 

**payinBank** object

Information about the receiving bank.

**payinBank.bankName** text

Bank name.

**payinBank.bankAddress** object (can be null)

Bank address.

**payinBank.bankAddress.country** text

Country ISO 2 code.

**payinBank.bankAddress.firstLine** text

Street address.

**payinBank.bankAddress.postCode** text

Post code / zip code.

**payinBank.bankAddress.city** text

City.

**payinBank.bankAddress.state** text (can be null)

State. Can be null.

**payinBank.bankAddress.phone** text (can be null)

Phone Number. Can be null.

**payinBankAccount** object

Bank account details to use to send the payment to.

**payinBankAccount.currency** text

ISO 4217 source currency code.

**payinBankAccount.bankAccountType** text

The type of bank account to use for the payin, e.g., RECIPIENT, ESCROW, BALANCE, etc.

**payinBankAccount.details** array of objects

Account details.

**payinBankAccount.details

[n].type** text

Account details type e.g. accountNumber, iban etc.

**payinBankAccount.details[n].label** text

Account details label that should be displayed in your user interface.

**payinBankAccount.details[n].value** text

Account details value - the value of account details (like iban, account number etc). This field is deprecated. Use `rawValue` and `formattedValue` instead

**payinBankAccount.details[n].rawValue** text

Account details value - the value of account details (like iban, account number etc). This field will always be unformatted and can be used for processing.

**payinBankAccount.details[n].formattedValue** text

Account details value - the value of account details (like iban, account number etc). This field is subject to change and intended to be human-readable, and is formatted according to the type of the field. This can be null.

**payinBankAccount.details[n].description** text

Description of the account detail type.

**wiseInformation** object

Information about the receiving Wise entity, the owner of the bank account.

**wiseInformation.localCompanyName** text

Wise local company name.

**wiseInformation.localAddress** object (can be null)

Wise local address.

**wiseInformation.localAddress.country** text

Country ISO 2 code.

**wiseInformation.localAddress.firstLine** text

Street address.

**wiseInformation.localAddress.postCode** text

Post code / zip code

**wiseInformation.localAddress.city** text

City.

**wiseInformation.localAddress.state** text (can be null)

State. Can be null.

**wiseInformation.localAddress.phone** text (can be null)

Phone Number. Can be null.

```json
{
  "payinBank": {
    "bankName": "Wise Europe SA/NV",
    "bankAddress": {
      "country": "BE",
      "countryIso2Code": "BE",
      "countryName": "Belgium",
      "firstLine": "Rue Du Trône 100, box 3",
      "postCode": "1050",
      "city": "Brussels",
      "state": null,
      "phone": "+32 472 123 456"
    }
  },
  "payinBankAccount": {
    "currency": "EUR",
    "bankAccountType": "RECIPIENT",
    "details": [
      {
        "type": "iban",
        "label": "IBAN",
        "value": "BE11111111111111",
        "rawValue": "BE11111111111111",
        "formattedValue": "BE11 1111 1111 1111",
        "description": "International Bank Account Number"
      },
      {
        "type": "bic",
        "label": "Bank code (BIC/Swift)",
        "value": "TRWIBEB1XXX",
        "rawValue": "TRWIBEB1XXX",
        "formattedValue": "TRWIBEB1XXX",
        "description": "BIC/SWIFT code of the receiving bank"
      },
      {
        "type": "recipientName",
        "label": "Recipient name",
        "value": "TransferWise Europe SA",
        "rawValue": "TransferWise Europe SA",
        "formattedValue": "TransferWise Europe SA",
        "description": "Name of the account holder"
      }
    ]
  },
  "wiseInformation": {
    "localCompanyName": "Wise Europe SA",
    "localAddress": {
      "country": "BE",
      "countryIso2Code": "BE",
      "countryName": "Belgium",
      "firstLine": "Rue du Trône 100, 3rd floor",
      "postCode": "1050",
      "city": "Brussels",
      "state": null,
      "phone": null
    }
  }
}
```

## Retrieve the bank transfer deposit details for a transfer 

**`GET /v1/profiles/{{profileId}}/transfers/{{transferId}}/deposit-details/bank-transfer`**

The payin deposit details API allows you to get the bank details for the account that the customer should send funds to when paying for a Wise transfer via a bank transfer. These details will be provided in the local format for that currency and usually contain bank account information - like iban, swift code etc. It also includes the name and address of the receiving bank (`payinBank`) and the name and address of the Wise entity that owns the bank account (`wiseInformation`) as sometimes these are required to make a payment.

{{profileId}} in the request url refers to the profile that created the transfer. It can be either the personal profile ID, or the business profile ID.

The `payinBankAccount` field allows the bank details to be displayed dynamically in a user interface, by displaying the label and value fields.

Currently, this API supports the following currencies:

- AUD
- BGN
- BRL
- CAD
- CHF
- CZK
- DKK
- EUR
- GBP
- HKD
- HRK
- HUF
- IDR
- INR
- JPY
- MYR
- NOK
- NZD
- PLN
- RON
- SEK
- SGD
- TRY
- USD
### Request

**profileId** text

Profile ID.

**transferId** text

Transfer ID.

#### Response

Returns a [payin deposit detail object](/api-reference/payin-deposit-detail#object).

```bash
curl -X GET \
  https://api.transferwise.com/v1/profiles/{{profileId}}/transfers/{{transferId}}/deposit-details/bank-transfer \
  -H 'Authorization: Bearer <your api token>' \
```