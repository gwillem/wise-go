<!-- Source: https://docs.wise.com/api-reference/balance -->

# Balance
Balance accounts are included as part of the Wise Multi-Currency Account.

Operations 

[POST/v4/profiles/{{profileId}}/balances](/api-reference/balance#create)

Create a balance account

[GET/v4/profiles/{{profileId}}/balances?types=STANDARD](/api-reference/balance#get)

Retrieve a balance by ID

[GET/v4/profiles/{{profileId}}/balances/{{balanceId}}](/api-reference/balance#list)

List balances for a profile

[DELETE/v4/profiles/{{profileId}}/balances/{{balanceId}}](/api-reference/balance#delete)

Remove a balance account

[POST/v2/profiles/{{profileId}}/balance-movements](/api-reference/balance#convert)

Convert across balance accounts

[GET/v1/profiles/{{profileId}}/balance-capacity](/api-reference/balance#capacity)

Returns regulatory deposit limit of an account

[POST/v1/profiles/{{profileId}}/excess-money-account](/api-reference/balance#excess-money-account)

Add an excess money account to a profile

[GET/v1/profiles/{{profileId}}/total-funds/{{currency}}](/api-reference/balance#total-funds)

Total worth of all balances in a profile converted to the specified currency.

## The Balance resource 

**id** integer, (int64)

Balance ID

**currency** text

Currency code (ISO 4217 Alphabetic Code)

**type** text

Type (`STANDARD`, `SAVINGS`)

**name** text

Name of the balance

**icon** text

Emoji string

**investmentState** text

Investment state of balance (`NOT_INVESTED`, `INVESTED`, `DIVESTING`, `UNKNOWN`)

**amount.value** decimal

Available balance value that can be used to fund transfers

**amount.currency** text

Currency code (ISO 4217 Alphabetic Code)

**reservedAmount.value** decimal

Amount reserved for transactions

**reservedAmount.currency** text

Currency code (ISO 4217 Alphabetic Code)

**cashAmount.value** decimal

Cash amount in the account

**cashAmount.currency** text

Currency code (ISO 4217 Alphabetic Code)

**totalWorth.value** decimal

Current total worth

**totalWorth.currency** text

Currency code (ISO 4217 Alphabetic Code)

**creationTime** timestamp

Date of when the balance was created

**modificationTime** timestamp

Date of when the balance was last modified

**visible** boolean

Balance is visible for the user or not

```json
{
  "id": 200001,
  "currency": "EUR",
  "type": "STANDARD",
  "name": null,
  "icon": null,
  "investmentState": "NOT_INVESTED",
  "amount": {
    "value": 0,
    "currency": "EUR"
  },
  "reservedAmount": {
    "value": 0,
    "currency": "EUR"
  },
  "cashAmount": {
    "value": 0,
    "currency": "EUR"
  },
  "totalWorth": {
    "value": 0,
    "currency": "EUR"
  },
  "creationTime": "2020-05-20T14:43:16.658Z",
  "modificationTime": "2020-05-20T14:43:16.658Z",
  "visible": true
}
```

## Create a balance account 

**`POST /v4/profiles/{{profileId}}/balances`**

This endpoint opens a balance within the specified profile, in the currency and type specified in the request.

For `STANDARD` balances, only one can be created for a currency. For `SAVINGS` balances, multiples in the same currency can be opened.

When sending the request, the `currency` and `type` are required. If creating a `SAVINGS` type balance, a `name` is also required.

### Request

**currency** text

Currency code (ISO 4217 Alphabetic Code)

**type** text

Type (`STANDARD`, `SAVINGS`)

**name** text

Name of the balance

#### Response

Returns a 

[balance object](#object).

```bash
curl -X POST \
  https://api.wise-sandbox.com/v4/profiles/{{profileId}}/balances \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -H 'X-idempotence-uuid: <generated uuid>' \
  -d '{
    "currency": "EUR",
    "type": "STANDARD"
  }'
```

## Retrieve a balance by ID 

**`GET /v4/profiles/{{profileId}}/balances/{{balanceId}}`**

This endpoint returns a balance based on the specified balance ID.

Returns a [balance object](#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v4/profiles/{{profileId}}/balances/{{balanceId}} \
  -H 'Authorization: Bearer <your api token>'
```

## List balances for a profile 

**`GET /v4/profiles/{{profileId}}/balances?types=STANDARD`**

Retrieve the user's multi-currency account balance accounts. It returns all balance accounts the profile has in the types specified.

A parameter of `type` must be passed and include at least a single type. To return more than one type, comma separate the values. Acceptable values are `STANDARD` and `SAVINGS`.

Returns an array of [balance objects](#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v4/profiles/{{profileId}}/balances?types=STANDARD \
  -H 'Authorization: Bearer <your api token>'
```

## Remove a balance account 

**`DELETE /v4/profiles/{{profileId}}/balances/{{balanceId}}`**

Close a balance account for the users profile. Balance accounts must have a zero balance in order for it to be closed. Bank account details for a balance account will also be deactivated and may not be restored in the future.

Returns a [balance object](#object).

```bash
curl -X DELETE \
  https://api.wise-sandbox.com/v4/profiles/{{profileId}}/balances/{{balanceId}} \
  -H 'Authorization: Bearer <your api token>'
```

## Convert across balance accounts 

**`POST /v2/profiles/{{profileId}}/balance-movements`**

This endpoint allows the conversion of funds between two `STANDARD` balance accounts in different currencies. You first must generate a quote with `"payOut": "BALANCE"`.

### Request

**profileId** integer

Authenticated profile ID

**quoteId** uuid

Quote ID - quote must be created with `"payOut": "BALANCE"`.

**X-idempotence-uuid** uuid

Unique identifier assigned by you. Used for idempotency check purposes. Should your call fail for technical reasons then you can use the same value again for making retry call.

Please also note that this call needs an extra field in the header called "X-idempotence-uuid". This should be generated and used for any subsequent retry call in the case that the initial `POST` fails.

Interested in moving money in/out of Jars?
 [Move money between balances](#move)

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/balance-movements \
  -H 'Authorization: Bearer <your api token>'  \
  -H 'Content-Type: application/json' \
  -H 'X-idempotence-uuid: <generated uuid>' \
  -d '{
    "quoteId": {{quote ID}}
  }'
```

```json
{
  "id": 30000001,
  "type": "CONVERSION",
  "state": "COMPLETED",
  "balancesAfter": [
    {
      "id": 1,
      "value": 10000594.71,
      "currency": "GBP"
    },
    {
      "id": 2,
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

## Move money between balances 

**`POST /v2/profiles/{{profileId}}/balance-movements`**

This endpoint allows the following money movements:

- Add money to same-currency jar (i.e *move* money from `STANDARD` to `SAVINGS` balance without conversion, `amount` is provided as request parameter);
- Add money to another-currency jar (i.e. *convert* money, `amount` is determined by provided `quoteId`);
- Withdraw money from jar (i.e *move* money from `SAVINGS` to `STANDARD` balance without conversion, `amount` is provided as request parameter).
### Request

**profileId** integer

Authenticated profile ID

**sourceBalanceId** integer

Source balance ID. If present, `targetBalanceId` is required.

**targetBalanceId** integer

Target balance ID. If present, `sourceBalanceId` is required.

**quoteId** uuid

Quote ID - quote must be created with `"payOut": "BALANCE"`. Required for cross-currency movements.

**amount** money

Either `amount` or `quoteId` are required.

**X-idempotence-uuid** uuid

Unique identifier assigned by you. Used for idempotency check purposes. If your call fails for technical reasons, you can use the same value again for making a retry call.

Interested in converting money from one balance to another?
 [Convert across balance accounts](#convert)

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/balance-movements \
  -H 'Authorization: Bearer <your api token>'  \
  -H 'Content-Type: application/json' \
  -H 'X-idempotence-uuid: <generated uuid>' \
  -d '{
    "amount": {
      "value": <decimal>,
      "currency": <currency>
    }
    "sourceBalanceId": <source ID>,
    "targetBalanceId": <target balance ID>,
    "quoteId": {{quote ID}}
  }'
```

```json
{
  "id": <conversion transaction ID>,
  "type": "CONVERSION",
  "state": "COMPLETED",
  "balancesAfter": [
    {
      "id": 1,
      "value": 10000594.71,
      "currency": "GBP"
    },
    {
      "id": 2,
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

## Retrieve the deposit limits of an account 

**`GET /v1/profiles/{{profileId}}/balance-capacity`**

This endpoint allows a user to check how much money they are able to deposit in their balance based on regulatory limits.

It is useful for personal profiles located in countries that have hold limits.

We advise to call this API before depositing money into an account if the profile is located in Singapore or Malaysia.

### Request

**currency** text

Currency code (ISO 4217 Alphabetic Code) - The deposit limit will be returned in this currency

### Response

**hasLimit** boolean

True if there is a regulatory hold limit for that profile's country

**depositLimit.amount** decimal

Amount of money that can be added to the account

**depositLimit.currency** text

Currency code (ISO 4217 Alphabetic Code) of money that can be added to the account. Specified in the request

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/balance-capacity?currency=SGD \
  -H 'Authorization: Bearer <your api token>'
```

```json
{
  "hasLimit": true,
  "depositLimit": {
    "amount": 2000.00,
    "currency": "SGD"
  }
}
```

## Add an excess money account to a profile 

**`POST /v1/profiles/{profileId}/excess-money-account`**

If a balance goes over the regulatory hold limit, we need to automatically move the excess amount to another account at the end of the day.

To do this, we require a recipient where the money will be moved to be added to the profile.

This API is primarily used for SG and MY customers.

### Request

**recipientId** integer

ID of the recipient for excess money transfers ([Recipient](/api-reference/recipient))

### Response

### Response

**userProfileId** integer

ID of the profile

**recipientId** integer

ID of the recipient for excess money transfers

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/excess-money-account \
  -H 'Authorization: Bearer <your api token> \'
  -d '{
    "recipientId": 148393305
  }'
```

```json
{
  "userProfileId": 12321323,
  "recipientId": 148393305
}
```

## Total funds 

**`GET /v1/profiles/{profileId}/total-funds/{currency}`**

This endpoint provides an overview of your account's total valuation and available liquidity across all balances.

#### Example (Assuming GBP and USD has 1:1 exchange rate)

| Scenario | GBP balance | USD balance | Total Worth | Total Available | Overdraft Usage | Overdraft Limit |
| --- | --- | --- | --- | --- | --- | --- |
| Positive account value with no overdraft | 2000 | 0 | 2000 | 2000 | 0 | 0 |
| Positive account value with overdraft | 2000 | -100 | 1900 | 2400 | 100 | 500 |
| Negative account value with overdraft | 0 | -100 | -100 | 400 | 100 | 500 |

### Request

**currency** text

A three-letter currency code that specifies the currency in which the returned values should be converted to.

### Response

**totalWorth** text

The total worth of the account, which includes the cash ledger balance and the valuation of any assets portfolio if the balance is invested.

**totalAvailable** text

The total available balance, which is the sum of the cash ledger balance and any approved overdraft limit.

**overdraft.limit** text

The maximum overdraft available to the customer through an overdraft program. The value is greater than 0 if there is an approved overdraft limit, else it will be 0.

**overdraft.usage** text

The portion of the approved overdraft limit that is currently being utilized by the account holder.

**overdraft.available** text

The amount of overdraft that is currently available. `overdraft.available`= `overdraft.limit`- `overdraft.usage`.

**overdraft.availableByCurrency** object

A list that shows the available overdraft amounts converted to each currency in which the customer has a balance.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/total-funds/EUR \
  -H 'Authorization: Bearer <your api token>'
```

```json
{
  "totalWorth": {
    "value": 2000.00,
    "currency": "EUR"
  },
  "totalAvailable": {
    "value": 2500.00,
    "currency": "EUR"
  },
  "overdraft": {
    "available": {
      "value": 500.00,
      "currency": "EUR"
    },
    "limit": {
      "value": 500.00,
      "currency": "EUR"
    },
    "used": {
      "value": 0.00,
      "currency": "EUR"
    },
    "availableByCurrency": [
      {
        "value": 500.00,
        "currency": "EUR"
      },
      {
        "value": 600.00,
        "currency": "USD"
      },
      {
        "value": 400.00,
        "currency": "GBP"
      }
    ]
  }
}
```