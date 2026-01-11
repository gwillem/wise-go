<!-- Source: https://docs.wise.com/api-reference/balance-statement -->

# Balance Statement
Balance statements contains transactional activities on a Wise Multi-Currency Account.

Operations 

[GET/v1/profiles/{{profileId}}/balance‑statements/{{balanceId}}/statement.json](/api-reference/balance-statement#get)

Retrieve a balance statement in JSON

## The Balance Statement resource 

**accountHolder.type** text

Account holder type (`PERSONAL`, `BUSINESS`)

**accountHolder.address.addressFirstLine** text

Account holder address street

**accountHolder.address.city** text

Account holder address city

**accountHolder.address.postCode** text

Account holder address zipcode

**accountHolder.address.stateCode** text

Account holder address state

**accountHolder.address.countryName** text

Account holder address country

**accountHolder.firstName** text

Account holder first name

**accountHolder.lastName** text

Account holder last name

**issuer.name** text

Account issuer name

**issuer.firstLine** text

Account issuer address street

**issuer.city** text

Account issuer address city

**issuer.postCode** text

Account issuer address zipcode

**issuer.stateCode** text

Account issuer address state

**issuer.country** text

Account issuer address country

**bankDetails** group

Your local bank details

**transactions

[n].type** text

Type (`DEBIT`, `CREDIT`)

**transactions[n].date** timestamp

Date of when transaction was created

**transactions[n].amount.value** decimal

Transaction amount

**transactions[n].amount.currency** text

Transaction currency code

**transactions[n].totalFees.value** decimal

Transaction fee amount

**transactions[n].totalFees.currency** text

Transaction fee currency code

**transactions[n].details.type** text

Type (`CARD`, `CONVERSION`, `DEPOSIT`, `TRANSFER`, `MONEY_ADDED`, `INCOMING_CROSS_BALANCE`, `OUTGOING_CROSS_BALANCE`, `DIRECT_DEBIT`, `BALANCE_INTEREST`, `BALANCE_ADJUSTMENT`, `UNKNOWN`, `ACCRUAL_CHARGE`, `INVESTMENT_TRADE_ORDER`, `ACQUIRING_PAYMENT`, `CARD_CASHBACK`, `CARD_ORDER_CHECKOUT`)

**transactions[n].details.description** text

Human readable explanation about the transaction

**transactions[n].details.amount.value** decimal

Amount in original currency (card transactions abroad)

**transactions[n].details.amount.currency** text

Original currency code (ISO 4217 Alphabetic Code)

**transactions[n].details.senderName** text

Deposit sender name

**transactions[n].details.senderAccount** text

Deposit sender bank account details

**transactions[n].details.paymentReference** text

Deposit payment reference text

**transactions[n].details.category** text

Card transaction category

**transactions[n].details.merchant.name** text

Card transaction merchant name

**transactions[n].details.merchant.firstLine** text

Merchant address street

**transactions[n].details.merchant.postCode** text

Merchant address zipcode

**transactions[n].details.merchant.city** text

Merchant address city

**transactions[n].details.merchant.state** text

Merchant address state

**transactions[n].details.merchant.country** text

Merchant address country

**transactions[n].details.merchant.category** text

Merchant category

**transactions[n].exchangeDetails.toAmount.value** decimal

Exchange target amount

**transactions[n].exchangeDetails.toAmount.currency** text

Exchange currency code (ISO 4217 Alphabetic Code)

**transactions[n].exchangeDetails.fromAmount.value** decimal

Exchange source amount

**transactions[n].exchangeDetails.fromAmount.currency** text

Exchange currency code (ISO 4217 Alphabetic Code)

**transactions[n].exchangeDetails.rate** decimal

Exchange rate

**transactions[n].runningBalance.value** decimal

Running balance after the transaction

**transactions[n].runningBalance.currency** text

Running balance currency code (ISO 4217 Alphabetic Code)

**transactions[n].referenceNumber** text

Wise assigned unique transaction reference number, this number can be used to map the refunds to the transfer that was refunded.

**endOfStatementBalance.value** decimal

Closing balance for specified time period

**endOfStatementBalance.currency** text

Closing balance currency code (ISO 4217 Alphabetic Code)

**query.intervalStart** timestamp

Query parameter repeated

**query.intervalEnd** timestamp

Query parameter repeated

**query.currency** text

Query parameter repeated

**query.accountId** integer

Query parameter repeated

```json
{
  "accountHolder": {
    "type": "PERSONAL",
    "address": {
      "addressFirstLine": "Veerenni 24",
      "city": "Tallinn",
      "postCode": "12112",
      "stateCode": "",
      "countryName": "Estonia"
    },
    "firstName": "Oliver",
    "lastName": "Wilson"
  },
  "issuer": {
    "name": "TransferWise Ltd.",
    "firstLine": "56 Shoreditch High Street",
    "city": "London",
    "postCode": "E1 6JJ",
    "stateCode": "",
    "country": "United Kingdom"
  },
  "bankDetails": null,
  "transactions": [
    {
      "type": "DEBIT",
      "date": "2018-04-30T08:47:05.832Z",
      "amount": {
        "value": -7.76,
        "currency": "EUR"
      },
      "totalFees": {
        "value": 0.04,
        "currency": "EUR"
      },
      "details": {
        "type": "CARD",
        "description": "Card transaction of 6.80 GBP issued by Tfl.gov.uk/cp TFL TRAVEL CH",
        "amount": {
          "value": 6.8,
          "currency": "GBP"
        },
        "category": "Transportation Suburban and Loca",
        "merchant": {
          "name": "Tfl.gov.uk/cp",
          "firstLine": null,
          "postCode": "SW1H 0TL  ",
          "city": "TFL TRAVEL CH",
          "state": "   ",
          "country": "GB",
          "category": "Transportation Suburban and Loca"
        }
      },
      "exchangeDetails": {
        "forAmount": {
          "value": 6.8,
          "currency": "GBP"
        },
        "rate": null
      },
      "runningBalance": {
        "value": 16.01,
        "currency": "EUR"
      },
      "referenceNumber": "CARD-249281"
    },
    {
      "type": "CREDIT",
      "date": "2018-04-17T07:47:00.227Z",
      "amount": {
        "value": 200,
        "currency": "EUR"
      },
      "totalFees": {
        "value": 0,
        "currency": "EUR"
      },
      "details": {
        "type": "DEPOSIT",
        "description": "Received money from HEIN LAURI with reference SVWZ+topup card",
        "senderName": "HEIN LAURI",
        "senderAccount": "EE76 1700 0170 0049 6704 ",
        "paymentReference": "SVWZ+topup card"
      },
      "exchangeDetails": null,
      "runningBalance": {
        "value": 207.69,
        "currency": "EUR"
      },
      "referenceNumber": "TRANSFER-34188888"
    },
    {
      "type": "CREDIT",
      "date": "2018-04-10T05:58:34.681Z",
      "amount": {
        "value": 9.94,
        "currency": "EUR"
      },
      "totalFees": {
        "value": 0,
        "currency": "EUR"
      },
      "details": {
        "type": "CONVERSION",
        "description": "Converted 8.69 GBP to 9.94 EUR",
        "sourceAmount": {
          "value": 8.69,
          "currency": "GBP"
        },
        "targetAmount": {
          "value": 9.94,
          "currency": "EUR"
        },
        "fee": {
          "value": 0.03,
          "currency": "GBP"
        },
        "rate": 1.147806
      },
      "exchangeDetails": null,
      "runningBalance": {
        "value": 9.94,
        "currency": "EUR"
      },
      "referenceNumber": "CONVERSION-1511237"
    }
  ],
  "endOfStatementBalance": {
    "value": 9.94,
    "currency": "EUR"
  },
  "query": {
    "intervalStart": "2018-03-01T00:00:00Z",
    "intervalEnd": "2018-04-30T23:59:59.999Z",
    "currency": "EUR",
    "accountId": 64
  }
}
```

## Retrieving a balance account statement 

**`GET /v1/profiles/{{profileId}}/balance-statements/{{balanceId}}/statement.json?currency=EUR&intervalStart=2018-03-01T00:00:00.000Z&intervalEnd=2018-03-15T23:59:59.999Z&type=COMPACT`**

This endpoint allows for statements to be generated for the provided balanceId, with the response in JSON. To generate in CSV, PDF, XLSX, CAMT.053, MT940 or QIF, replace statement.json with statement.csv, statement.pdf, statement.xlsx, statement.xml, statement.mt940 or statement.qif respectively in the above URL. Note that the PDF includes Wise branding.

The period between intervalStart and intervalEnd cannot exceed 469 days (around 1 year 3 months).

This endpoint is SCA protected when it applies. If your profile is registered within the UK and/or EEA, SCA most likely applies to you. The additional authentication is only required once every 90 days, viewing the statement on the website or in the mobile app counts towards that as well.
 [Learn more](/guides/developer/auth-and-security/sca-and-2fa)

### Request

**currency** text

Currency of the balance statement requested

**intervalStart** timestamp

Statement start time in UTC time

**intervalEnd** timestamp

Statement end time in UTC time

**type** text

- `COMPACT` for a single statement line per transaction
- `FLAT` for accounting statements where transaction fees are on a separate line
**statementLocale** text

Language that you wish the statement to be in. Supports 2 character language codes

#### Response

Returns a [balance statement object](#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/balance-statements/{{balanceId}}/statement.json
    ?currency=EUR
    &intervalStart=2018-03-01T00:00:00.000Z
    &intervalEnd=2018-03-15T23:59:59.999Z
    &type=COMPACT \
  -H 'Authorization: Bearer <your api token>'
```