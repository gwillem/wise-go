<!-- Source: https://docs.wise.com/api-reference/claim-account -->

# Claim Account
Operations 

[POST/v1/user/claim-account](/api-reference/claim-account#create)

Generate a claim account code

## Generate a claim account code 

**`POST /v1/user/claim-account`**

The `claim_account_code` is meant to be used as a request parameter when redirecting a new customer to Wise, effectively allowing the customer to claim the account in question as their own.

Use the `registration_code` used for 

[creating the user](/api-reference/user#createuser-regcode)

### Request

**registrationCode** text

The `registration_code` belonging to the user

### Response

**claimAccountCode** text

The `claim_account_code` to be used in the redirect to Wise

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/user/claim-account \
  -H 'Authorization: Bearer <your user token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "registrationCode": <registration code>
  }'
```

```json
{
  "claimAccountCode": "<claim_account_code>"
}
```