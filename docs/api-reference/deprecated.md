<!-- Source: https://docs.wise.com/api-reference/deprecated -->

# Deprecated APIs
Deprecated APIs are those that are no longer the standard version and have been marked as such. They will continue to work for the foreseeable future, but should not be used for new integrations.

All endpoints below are deprecated. Please see the recommended alternative and use that instead of the endpoint included here.

Operations 

[GET/v1/borderless-accounts?profileId={{profileId}}](/api-reference/deprecated#get-account-balance-v1)

Get Account Balance - v1

[POST/v1/borderless-accounts/{{borderlessAccountId}}/conversions](/api-reference/deprecated#convert-currencies-v1)

Convert Currencies - v1

[GET/v1/identity/one-time-token/status](/api-reference/deprecated#get-ott-status-v1)

Get OTT status

[PUT/v3/spend/profiles/{{profileId}}/cards/{{cardToken}}/phone-number](/api-reference/deprecated#update-phone-number)

Update phone number

[POST/v3/spend/profiles/{{profileId}}/dispute-form/file](/api-reference/deprecated#dispute-file-v3)

Upload file

## Get account balance 

**`GET /v1/borderless-accounts?profileId={{profileId}}`**

This endpoint is deprecated. Please see the 

[v4 Balances endpoints](/api-reference/balance#get).

Get available balances for all activated currencies in your multi-currency account.

Use profile ID obtained earlier to make this call.

### Response

**id** integer

Multi-currency account ID

**profileId** integer

Personal or business profile ID

**recipientId** integer

Recipient ID you can use for a multi-currency account deposit

**creationTime** timestamp

Date when multi-currency account was opened

**modificationTime** timestamp

Date when multi-currency account setup was modified.

**active** boolean

Is multi-currency account active or inactive

**eligible** boolean

Ignore

**balances[n].balanceType** text

AVAILABLE

**balances[n].currency** text

Currency code

**balances[n].amount.value** decimal

Available balance in specified currency

**balances[n].amount.currency** text

Currency code

**balances[n].reservedAmount.value** decimal

Reserved amount from your balance

**balances[n].reservedAmount.currency** text

Reserved amount currency code

**balances[n].bankDetails** group

Bank account details assigned to your multi-currency account. Available for EUR, GBP, USD, AUD, NZD

**balances[n].bankDetails.id** integer

Bank account details ID

**balances[n].bankDetails.currency** text

Bank account currency

**balances[n].bankDetails.bankCode** text

Bank account code

**balances[n].bankDetails.accountNumber** text

Bank account number

**balances[n].bankDetails.swift** text

Bank account swift code

**balances[n].bankDetails.iban** text

Bank account iban

**balances[n].bankDetails.bankName**

Bank name

**balances[n].bankDetails.accountHolderName** text

Bank account holder name

**balances[n].bankDetails.bankAddress.addressFirstLine** text

Bank address street and house

**balances[n].bankDetails.bankAddress.postCode** text

Bank address zip code

**balances[n].bankDetails.bankAddress.city** text

Bank address city

**balances[n].bankDetails.bankAddress.country** text

Bank address country

**balances[n].bankDetails.bankAddress.stateCode** text

Bank address state code

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/borderless-accounts?profileId={{profileId}} \
  -H 'Authorization: Bearer <your api token>'
```

```json
[
  {
    "id": 64,
    "profileId": 33333333,
    "recipientId": 13828530,
    "creationTime": "2018-03-14T12:31:15.678Z",
    "modificationTime": "2018-03-19T15:19:42.111Z",
    "active": true,
    "eligible": true,
    "balances": [
      {
        "balanceType": "AVAILABLE",
        "currency": "GBP",
        "amount": {
          "value": 10999859,
          "currency": "GBP"
        },
        "reservedAmount": {
          "value": 0,
          "currency": "GBP"
        },
        "bankDetails": null
      },
      {
        "balanceType": "AVAILABLE",
        "currency": "EUR",
        "amount": {
          "value": 9945236.2,
          "currency": "EUR"
        },
        "reservedAmount": {
          "value": 0,
          "currency": "EUR"
        },
        "bankDetails": {
          "id": 90,
          "currency": "EUR",
          "bankCode": "DEKTDE7GXXX",
          "accountNumber": "DE51 7001 1110 6050 1008 91",
          "swift": "DEKTDE7GXXX",
          "iban": "DE51 7001 1110 6050 0008 91",
          "bankName": "Handelsbank",
          "accountHolderName": "Oliver Wilson",
          "bankAddress": {
            "addressFirstLine": "Elsenheimer Str. 41",
            "postCode": "80687",
            "city": "München",
            "country": "Germany",
            "stateCode": null
          }
        }
      }
    ]
  }
]
```

## Convert currencies between balances 

**`POST /v1/borderless-accounts/{{borderlessAccountId}}/conversions`**

This endpoint is deprecated. Please check new [balance endpoints](/api-reference/balance#convert) to move and convert between balances.

Convert funds between your multi-currency account currencies. Quote which is used in this call must be created with `"payOut": "BALANCE"`.

Note that this call needs an extra field in header called "X-idempotence-uuid".

### Request

**quoteId** integer

Conversion quote ID

**X-idempotence-uuid** uuid

Unique identifier assigned by you. Used for idempotency check purposes. Should your call fail for technical reasons then you can use the same value again for making retry call.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/borderless-accounts/{{borderlessAccountId}}/conversions \
  -H 'Authorization: Bearer <your api token>'  \
  -H 'Content-Type: application/json' \
  -H 'X-idempotence-uuid: <generated uuid>' \
  -d '{
    "quoteId": <conversion quote ID>
  }'
```

```json
{
  "id": 1,
  "type": "CONVERSION",
  "state": "COMPLETED",
  "balancesAfter": [
    {
      "value": 10000594.71,
      "currency": "GBP"
    },
    {
      "value": 9998887.01,
      "currency": "EUR"
    }
  ],
  "creationTime": "2017-11-21T09:55:49.275Z",
  "steps": [
    {
      "id": 369588,
      "type": "CONVERSION",
      "creationTime": "2017-11-21T09:55:49.276Z",
      "balancesAfter": [
        {
          "value": 9998887.01,
          "currency": "EUR"
        },
        {
          "value": 10000594.71,
          "currency": "GBP"
        }
      ],
      "channelName": null,
      "channelReferenceId": null,
      "tracingReferenceCode": null,
      "sourceAmount": {
        "value": 113.48,
        "currency": "EUR"
      },
      "targetAmount": {
        "value": 100,
        "currency": "GBP"
      },
      "fee": {
        "value": 0.56,
        "currency": "EUR"
      },
      "rate": 0.88558
    }
  ],
  "sourceAmount": {
    "value": 113.48,
    "currency": "EUR"
  },
  "targetAmount": {
    "value": 100,
    "currency": "GBP"
  },
  "rate": 0.88558,
  "feeAmounts": [
    {
      "value": 0.56,
      "currency": "EUR"
    }
  ]
}
```

## Get status of a one time token 

**`GET /v1/identity/one-time-token/status`**

This endpoint is deprecated. Please check new [the new endpoint](/api-reference/one-time-token#get-status) to get status of a one time token.

Retrieve necessary information to clear a OTT.

Request header**One-Time-Token** text

Text value of a OTT.

### Response

**oneTimeTokenProperties** OneTimeToken

Properties of OneTimeToken

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/identity/one-time-token/status \
  -H 'Authorization: Bearer <your api token>'
  -H 'One-Time-Token: <one time token>'
```

```json
{
  "oneTimeTokenProperties": {
    "oneTimeToken": "9f5f5812-2609-4e48-8418-b64437c0c7cd",
    "challenges": [
      {
        "primaryChallenge": {
          "type": "PIN",
          "viewData": {
            "attributes": {
              "userId": 6146956
            }
          }
        },
        "alternatives": [],
        "required": true,
        "passed": false
      }
    ],
    "validity": 3600,
    "actionType": "BALANCE__GET_STATEMENT",
    "userId": 6146956
  }
}
```

## Update phone number associated with a card 

**`PUT /v3/spend/profiles/{{profileId}}/cards/{{cardToken}}/phone-number `**

Please note that this endpoint is deprecated for partners onboarded after 1/3/2025. For new and existing partners, your clientId will need to be whitelisted by our support team to make the `phoneNumber` optional in the [create card order endpoint](/api-reference/card-order#create-card-order). See [3ds for more information](/guides/product/issue-cards/3ds). This endpoint updates the phone number of a card. The new phone number must be a valid phone number.

### Request

**phoneNumber** text

Must be a valid phone number prefixed with + and country code. An example of a valid phone number would be `+6588888888`

### Response

**token** text

The card token that you are modifying

**profileId** text

The profile ID that the card is linked to

**phoneNumber** text

The new phone number associated with the card

```bash
curl -X PUT \
  'https://api.wise-sandbox.com/v3/spend/profiles/{{profileId}}/cards/{{cardToken}}/phone-number' \
  --H 'Authorization: Bearer <your api token>' \
  --d '{
    "phoneNumber": "+6588888888
  }'
```

```json
{
  "token": "12345-12345-12345-12345",
  "profileId": 30000000,
  "phoneNumber": "+6588888888"
}
```

## Submitting dispute file V3 

**`POST /v3/spend/profiles/{{profileId}}/dispute-form/file`**

Submit a file for disputes

This API is not available for sandbox testing.

Use the `fileId` in the response for the dispute submission.

**NB:** A dispute referencing the returned `fileId` must be submitted no later than **two hours** after the file submission, otherwise the file will expire and must be re-submitted.

### Response

**[fileName: text]** text

ID of the file to be used on dispute submission.

```bash
curl -X POST \
  'https://api.wise.com/v3/spend/profiles/{{profileId}}/dispute-form/file' \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: multipart/form-data'  \
  -F 'receipt=@"<your file>"'
```

```json
{
  "receipt" : "ab4f5d2"
}
```