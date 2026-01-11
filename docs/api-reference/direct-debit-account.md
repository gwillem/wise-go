<!-- Source: https://docs.wise.com/api-reference/direct-debit-account -->

# Direct Debit Account
Operations 

[POST/v1/profiles/{{profileId}}/direct-debit-accounts](/api-reference/direct-debit-account#create)

Create a direct debit account

[GET/v1/profiles/{{profileId}}/direct-debit-accounts?type={{type}}&currency={{currencyCode}}](/api-reference/direct-debit-account#get)

Retrieve all direct debit accounts

## The Direct Debit Account resource 

A *direct debit account* represents details of your external bank account. It is used for funding batches via direct debit.

Properties**id** number

Direct debit account ID

**currency** text

Currency of a bank account. 3 character currency code. Should be "USD" for ACH direct debits.

**type** text

Payment type for which account is used

**details** object

Currency specific fields

**details.routingNumber** text

A valid ABA routing number. For CAD direct debits, the routing number is a combination of an institution number + transit number.

**details.accountNumber** text

A valid bank account number

**details.accountType** text

Bank account type

### Bank Account Type

| Name | Description | Currencies |
| --- | --- | --- |
| CHECKING | Checking account | USD, CAD |
| SAVINGS | Savings account | USD, CAD |

### Direct Debit Account Payment Type

| Name | Description | Currencies |
| --- | --- | --- |
| ACH | US Direct debit | USD |
| EFT | CA Direct debit | CAD |

```json
{
  "id": 1,
  "currency": "USD",
  "type": "ACH",
  "details": {
    "routingNumber": "000000000",
    "accountNumber": "0000000000",
    "accountType": "CHECKING"
  }
}
```

## Create a direct debit account 

**`POST /v1/profiles/{{profileId}}/direct-debit-accounts`**

Create a new direct debit account.

### Request

**currency** text

Currency of a bank account. 3 character currency code. Should be "USD" for ACH direct debits.

**type** text

"ACH" for USD direct debits; "EFT" for CAD direct debits.

**details** object

Currency specific fields

**details.routingNumber** text

A valid ABA routing number. For CAD direct debits, the routing number is a combination of an institution number + transit number.

**details.accountNumber** text

A valid bank account number

**details.accountType** text

Bank account type (Either CHECKING or SAVINGS)

#### Response

Returns a 

[direct debit account object](/api-reference/direct-debit-account#object).

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/direct-debit-accounts \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "currency": "USD",
    "type": "ACH",
    "details": {
      "routingNumber": "000000000",
      "accountNumber": "0000000000",
      "accountType": "CHECKING"
    }
  }'
```

## Retrieve all direct debit accounts 

**`GET /v1/profiles/{{profileId}}/direct-debit-accounts?type={{type}}&currency={{currencyCode}}`**

Get a list of your direct debit accounts. Use "type" and "currency" query parameters to filter accounts by type and currency.

### Request

**type** text

Required. A payment type for which account is used.

**currency** text

Required. Currency of a bank account. 3 character currency code. Should be "USD" for ACH direct debits; "CAD" for EFT direct debits.

#### Response

Returns a [direct debit account object](/api-reference/direct-debit-account#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/direct-debit-accounts?type={{type}}&currency={{currencyCode}} \
  -H 'Authorization: Bearer <your api token>'
```