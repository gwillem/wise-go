<!-- Source: https://docs.wise.com/api-reference/multi-currency-account -->

# Multi Currency Account
The Wise multi-currency account (MCA) enables users to hold, convert, and fund transfers (single or batches) with balances in up to 56 currencies. Please note that of the 50+ currency balances that we support, 10+ of these come with local account details.

Please refer to our [multi-currency account guide](/guides/product/accounts/managing-accounts) for more information on the general use of the endpoints included below.

Please see the [Balances APIs](/api-reference/balance) for more details.

Operations 

[GET/v2/borderless-accounts-configuration/profiles/{{profileId}}/available-currencies](/api-reference/multi-currency-account#available-currencies)

Retrieve available currencies

[GET/v2/borderless-accounts-configuration/profiles/{{profileId}}/payin-currencies](/api-reference/multi-currency-account#payin-currencies)

Retrieve payin currencies

[GET/v4/profiles/{{profileId}}/multi-currency-account](/api-reference/multi-currency-account#get-account)

Retrieve multi currency account for a profile

[GET/v4/multi-currency-account/eligibility?profileId={{profileId}}](/api-reference/multi-currency-account#get-eligibility)

Retrieve multi currency account eligibility for a profile

[GET/v4/multi-currency-account/eligibility?country={{country-ISO-Code}}&state={{state-ISO-code}}](/api-reference/multi-currency-account#get-eligibility)

Retrieve multi currency account eligibility for a country and state

## The Multi Currency Account resource 

**id** integer

Multi currency account ID

**profileId** integer

Profile ID the multi currency account is attributed to.

**recipientId** integer

Recipient ID of the multi currency account, to be used for transfer recipient

**creationTime** datetime

Datetime when multi currency account was created

**modificationTime** datetime

Datetime when multi currency account was last modified

**active** bool

Whether multi currency account is active or not

**eligible** bool

Whether multi currency account is eligible or not

```json
{
  "id": 1,
  "profileId": 33333333,
  "recipientId": 12345678,
  "creationTime": "2020-05-20T14:43:16.658Z",
  "modificationTime": "2020-05-20T14:43:16.658Z",
  "active": true,
  "eligible": true
}
```

## Retrieve available currencies 

Two endpoints exist to retrieve all the currencies available for balance accounts. You can use this list to create a balance account for the currency included.

#### Available Currencies

**`GET /v2/borderless-accounts-configuration/profiles/{{profileId}}/available-currencies`**

Lists all currencies that are available for balance accounts. This includes those that can have funds added from external accounts, as well as those that a balance can be held in.

#### Payin Currencies

**`GET /v2/borderless-accounts-configuration/profiles/{{profileId}}/payin-currencies`**

Lists all the currencies available for balance accounts that are also available to have bank account details. You can use this list to create a balance account for the currency included and then subsequently create bank account details.

#### Response

Returns a list of currencies supported for that option.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v2/borderless-accounts-configuration/profiles/{{profileId}}/available-currencies \
  -H 'Authorization: Bearer <your api token>' 
```

```bash
curl -X GET \
  https://api.wise-sandbox.com/v2/borderless-accounts-configuration/profiles/{{profileId}}/payin-currencies \
  -H 'Authorization: Bearer <your api token>' 
```

```json

[
  "EUR",
  "GBP",
  "USD",
  ...
]
```

## Retrieve multi currency account for a profile 

**`GET /v4/profiles/{{profileId}}/multi-currency-account`**

This endpoint returns the multi-currency account details for the specified profileId. If the user does not yet have a multi-currency account, a `404 Not Found` will be returned.

#### Response

Returns a [multi currency account object](/api-reference/multi-currency-account#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v4/profiles/{{profileId}}/multi-currency-account \
  -H 'Authorization: Bearer <your api token>'
```

## Retrieve multi currency account eligibility 

**`GET /v4/multi-currency-account/eligibility`**

This endpoint checks eligibility for a multi-currency account for either a specific profile or for a location. Customers in some countries and states/provinces may not be eligible for a multi currency account.

To check a profile, the `profileId` should be passed as a parameter.

To check a specific location, the country the user is in should be passed as `country` using 2-letter ISO 3166 codes. If the country is `US`, a valid 2 letter `state` parameter must also be passed.

- 

Ex 1: France: `/v4/multi-currency-account/eligibility?country=FR`

- 

Ex 2: USA, California: `/v4/multi-currency-account/eligibility?country=US&state=CA`

### Response

**eligible** boolean

Profile is eligible for MCA Account

**eligibilityCode**"eligible", "invalid.profile.type", "invalid.country", or "invalid.state"

Reason for the ineligibility

**accountType**"ineligible", "receive_only" or "full". generally this will be returned as "full".

Account type available

**ineligibilityReason** string

Reason the profile is not eligible

```bash
curl -X GET \
  https://api.wise-sandbox.com/v4/multi-currency-account/eligibility?profileId={{profileId}} \
  -H 'Authorization: Bearer <your api token>'
```

```bash
curl -X GET \
  https://api.wise-sandbox.com/v4/multi-currency-account/eligibility?country={{country-ISO-Code}}&state={{state-ISO-code}} \
  -H 'Authorization: Bearer <your api token>'
```

```json
{
  "eligible": true,
  "eligibilityCode": "eligible",
  "accountType": "FULL",
  "ineligibilityReason": null
}
```