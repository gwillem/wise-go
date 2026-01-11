<!-- Source: https://docs.wise.com/api-reference/bank-account-details -->

# Bank Account Details
Operations 

[GET/v1/profiles/{{profileId}}/account-details](/api-reference/bank-account-details#get)

Retrieve bank account details for a profile

[POST/v1/profiles/{{profileId}}/account-details-orders](/api-reference/bank-account-details#create)

Create a bank account details order

[POST/v3/profiles/{{profileId}}/bank-details](/api-reference/bank-account-details#create-multiple-details)

Create multiple bank account details

[GET/v1/profiles/{{profileId}}/account-details-orders](/api-reference/bank-account-details#list)

List bank account detail orders

[POST/v1/profiles/{{profileId}}/account-details/payments/{{paymentId}}/returns](/api-reference/bank-account-details#create-payment-return)

Create return for payment received to bank account details

## The Bank Account Details resource 

**id** integer

Bank account details ID

**currency.code** text

Currency code (ISO 4217 Alphabetic Code)

**currency.name** text

Currency name

**title** text

Account title

**subtitle** text

Account subtitle

**status** text

- `AVAILABLE`: Account details do not exist for the user but may be created
- `ACTIVE`: Account details are ready to be used by this user
**deprecated** boolean

*Important!* When the value is `true`, Wise issued another account details for the same currency. The users should not use the deprecated account anymore, but they still need to be informed as they might have external references pointing to the old one.

**receiveOptions** array

This array will contain the available receive options for the given currency.

- `LOCAL`: local bank details to receive money in the account currency
- `INTERNATIONAL`: Swift bank details to receive money internationally in the account currency.
**receiveOptions.details** array

It returns all the details that need to be displayed to users to make transfers to this account. This array will vary from currency to currency as each one has its requirements. For example, EUR will return `ACCOUNT_HOLDER`, `SWIFT_CODE`, `IBAN` and `ADDRESS`, whereas `GBP` will return `ACCOUNT_HOLDER`, `BANK_CODE`, `ACCOUNT_NUMBER`, `IBAN` and `ADDRESS`.

Each element in this array will contain the following attributes:

- `type`: String informing the account detail type
- `title`: Label to display in the UI
- `body`: Value to display in the UI
- `description`: When present, it has the content to show a tooltip/popup hint
- `hidden`: If the field should be displayed or not in the UI
**bankFeatures** array

Features enabled on the account

```json
{
  "id": 14000001,
  "currency": {
    "code": "EUR",
    "name": "Euro"
  },
  "title": "Your EUR account details",
  "subtitle": "IBAN, SWIFT/BIC",
  "status": "ACTIVE",
  "deprecated": false,
  "receiveOptions": 

[
  ],
  "bankFeatures": [
    {
      "key": "LOCAL_RECEIVE",
      "title": "Receive locally",
      "supported": true
    },
    {
      "key": "SWIFT",
      "title": "Receive internationally (Swift)",
      "supported": true
    },
    {
      "key": "DIRECT_DEBITS",
      "title": "Set up Direct Debits",
      "supported": true
    },
    {
      "key": "PLATFORM_RECEIVE",
      "title": "Receive from PayPal and Stripe",
      "supported": true
    }
  ]
}
```

## Retrieve bank account details for a profile 

**`GET /v1/profiles/{{profileId}}/account-details`**

This endpoint returns a list with all the `AVAILABLE` and `ACTIVE` account details for the given profile, including examples. Account receive options can also include local and international details to receive money on the currency balance.

Example bank account details are returned for any currency where bank account details have not been requested and issued. Examples will always include an `id` of `null`.

#### Response

Returns a list of [bank account details objects](#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/account-details \
  -H 'Authorization: Bearer <your api token>'
```

## Create a bank account details order 

**`POST /v1/profiles/{{profileId}}/account-details-orders`**

Creates an order which will issue account details. It should use the same currency as the balance previously created. Fulfilling all the requirements will complete the order. That means reaching status `DONE`.

The possible values for a requirement are:

- `PENDING_USER`: The requirement has some pending action from the user.
- `PENDING_TW`: The requirement has some pending action from Wise.
- `DONE`: The requirement is completed

The more common requirements are:

- `VERIFICATION`: The user needs to be fully verified before complete this requirement.
- `TOP_UP`: A fee will be charged and must be paid using through wise.com before complete this requirement.
### Request

**currency** text

balance currency (ISO 4217 Alphabetic Code)

### Response

**status** text

Status (`PENDING_USER`, `PENDING_TW`, `DONE`)

**currency** text

currency (ISO 4217 Alphabetic Code)

**requirements.type** text

Requirement type (`VERIFICATION`, `TOP_UP`)

**requirement.status** text

Requirement status (`PENDING_USER`, `PENDING_TW`, `DONE`)

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/account-details-orders \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "currency": "EUR"
  }'
```

```json
{
  "status": "PENDING_USER",
  "currency": "EUR",
  "requirements": [
    {
      "type": "VERIFICATION",
      "status": "PENDING_USER"
    }
  ]
}
```

## Create multiple bank account details 

**`POST /v3/profiles/{{profileId}}/bank-details`**

Please reach out to our Support Team for access to this endpoint

Creates and assigns a pair of local account details and international account details (where available) that are linked to the target balance specified in the request.

### Request

**targetAccountId** integer

ID of the currency balance to create account details on

### Response

**id** string

Account detail ID

**currency** string

Account detail currency (ISO 4217 Alphabetic Code)

**active** boolean

Status of the account detail

**localDetails.bankName** string

Name of the bank

**localDetails.bankAddress** string

Address of the bank

**localDetails.sortCode** string

Sort code of the bank

**localDetails.accountNumber** string

Bank account number

**localDetails.type** string

Type of account detail

**internationalDetails.bankName** string

Name of the bank

**internationalDetails.bankAddress** string

Address of the bank

**internationalDetails.swiftCode** string

Bank SWIFT code

**internationalDetails.iban** string

IBAN

**internationalDetails.type** string

Type of account detail

```bash
curl -X POST \
  https://api.wise-sandbox.com/v3/profiles/{{profileId}}/bank-details \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "targetAccountId": "123456"
  }'
```

```json
{
    "id": "2",
    "currency": "GBP",
    "active": true,
    "localDetails": {
        "bankName": "Wise",
        "bankAddress": "TEA BUILDING, FLOOR 6, SHOREDITCH HIGH STREET",
        "sortCode": "231370",
        "accountNumber": "00000001",
        "type": "UK_ACCOUNT"
    },
    "internationalDetails": { // optional, only if SWIFT payments are supported for that currency
        "bankName": "Wise",
        "bankAddress": "TEA BUILDING, FLOOR 6, SHOREDITCH HIGH STREET",
        "swiftCode": "TRWIGB22XXX",
        "iban": "GB123450000000001",
        "type": "IBAN"
    }
}
```

## List bank account detail orders v3

**`GET /v3/profiles/{{profileId}}/account-details-orders?currency=GBP`**

This endpoint returns the bank account assignment requests for a profile and multi-currency account. The response includes bank‑details orders in the following statuses: (`PENDING_USER`, `PENDING_TW`, `REQUIREMENTS_FULFILLED`, `DONE`)

**currency** text

Currency code (ISO 4217 Alphabetic Code)

```bash
curl -X GET \
  https://api.wise-sandbox.com/v3/profiles/{{profileId}}/account-details-orders?currency=EUR \
  -H 'Authorization: Bearer <your api token>'
```

```json
[
  {
    "status": "DONE",
    "currency": "EUR",
    "requirements": [
      {
        "type": "TOP_UP",
        "status": "DONE",
      },
      {
        "type": "VERIFICATION",
        "status": "DONE",
      }
    ]
  },
  ...
]
```

## Create return for payment to bank account details 

**`POST /v1/profiles/{{profileId}}/account-details/payments/{{paymentId}}/returns`**

Creates a return for a payment received to bank account details.

When you create a return, you must provide the ID of the payment you wish to return as well as the ID of the profile that received the payment. In addition, you can provide a `reason` for the return in the request body. When returning Swift payments, `reason` is a required field.

### Request

**reason** text

Reason for the return. Only required to return Swift payments.
Reason(`INCORRECT_ACCOUNT_NUMBER`, `CLOSED_ACCOUNT`, `BLOCKED_ACCOUNT`, `CANCELLATION_REQUEST`, `REGULATORY`, `CUSTOMER_REQUEST`)
*optional*

### Response

**id** uuid

ID of the return created.

*required*

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/account-details/payments/{{paymentId}}/returns \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "reason": "CLOSED_ACCOUNT"
  }'
```

```json
{
  "id": "4cc39f2b-3513-453d-8792-9ccc22e513c3"
}
```

```json
{
  "errors": [
    "code": "error.payment.not-found",
    "message": "No payment found with id [123456789] for profile id [987654321]",
    "arguments": [
      987654321,
      123456789,
      null
    ]
  ]
}
```

```json
{
  "errors": [
    "code": "error.payment-return.invalid",
    "message": "A valid reason is required to return a Swift payment",
    "arguments": [
      987654321,
      123456789,
      null
    ]
  ]
}
```