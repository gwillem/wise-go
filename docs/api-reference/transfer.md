<!-- Source: https://docs.wise.com/api-reference/transfer -->

# Transfer
A transfer is a payment order to [recipient account](/api-reference/recipient) based on a [quote](/api-reference/quote).

Once created, a transfer needs to be funded within the next fourteen days. Otherwise, it will be automatically canceled.

Operations 

[POST/v1/transfers](/api-reference/transfer#create)

Create a transfer

[POST/v2/profiles/{{profileId}}/third-party-transfers](/api-reference/transfer#create-third-party)

Create a third party transfer

[POST/v1/profiles/{{profileId}}/partner-licence-transfers](/api-reference/transfer#create-partner-license)

Create a partner license transfer

[POST/v1/transfer-requirements](/api-reference/transfer#transfer-requirements)

Validate transfer requirements

[GET/v1/transfers/{{transferId}}](/api-reference/transfer#get-by-id)

Get a transfer by ID

[GET/v2/profiles/{{profileId}}/third-party-transfers/{{transferId}}](/api-reference/transfer#get-third-party-transfer-by-id)

Get a third party transfer by ID

[POST/v3/profiles/{{profileId}}/transfers/{{transferId}}/payments](/api-reference/transfer#fund)

Fund a transfer

[GET/v1/transfers](/api-reference/transfer#list)

List transfers

[PUT/v1/transfers/{{transferId}}/cancel](/api-reference/transfer#cancel)

Cancel a transfer

[GET/v1/transfers/{{transferId}}/receipt.pdf](/api-reference/transfer#receipt-pdf)

Get a transfer's receipt by ID

[GET/v1/transfers/{{transferId}}/us-combined-receipt.pdf](/api-reference/transfer#us-combined-disclosure-receipt-pdf)

Get a transfer's US combined disclosure receipt by ID

[GET/v1/transfers/{{transferId}}/documents/noc](/api-reference/transfer#noc-inr)

Get a transfer's non objection certificate (NOC) by ID (INR)

[GET/v1/transfers/{{transferId}}/payments](/api-reference/transfer#payments)

List of completed payments

[GET/v2/transfers/{{transferId}}/invoices/bankingpartner](/api-reference/transfer#banking-reference-number)

Get details on how the transfer was paid out

## The Transfer resource 

There are 2 types of transfer resources: standard transfer resource and originator transfer resource. Please see below for differences between the two.

### Standard Transfer

**id** integer, (int64)

Transfer ID

**user** integer, (int64)

Your user ID

**targetAccount** integer, (int64)

Recipient account ID

**sourceAccount** integer, (int64)

Refund recipient account ID

**quote** integer, (int64)

v1 quote ID

**quoteUuid** text

v2 quote ID

**status** text

Transfer current status 

[For information about transfer statuses](/guides/product/send-money/tracking-transfers)

**reference** text

Deprecated, use details.reference instead

**rate** decimal

Exchange rate value

**created** timestamp

Timestamp when transfer was created

**business** integer, (int64)

Your business profile ID

**transferRequest** integer

Deprecated

**details.reference** text

Payment reference text

**hasActiveIssues** boolean

Are there any pending issues which stop executing the transfer?

**sourceCurrency** text

Source currency code

**sourceValue** decimal

Transfer amount in source currency

**targetCurrency** text

Target currency code

**targetValue** decimal

Transfer amount in target currency

**customerTransactionId** text

Unique identifier randomly generated per transfer request by the calling client

**payinSessionId** text

ID of the Payin Session generated for the transfer, which can be used for certain payin methods when funding the transfer.

**originator** group

Data block to capture payment originator details

**originator.legalEntityType** text

Payment originator legal type

**originator.reference** text

Unique customer ID in your system

**originator.name.givenName** text

Payment originator first name

**originator.name.middleNames** text array

Payment originator middle name(s)

**originator.name.familyName** text

Payment originator family name

**originator.name.patronymicName** text

Payment originator patronymic name

**originator.name.fullName** text

Payment originator full legal name

**originator.dateOfBirth** yyyy-mm-dd

Payment originator date of birth

**originator.businessRegistrationCode** text

Payment originator business registry or incorporation number

**originator.address.firstLine** text

Payment originator address first line

**originator.address.city** text

Payment originator address city

**originator.address.stateCode** text

Payment originator address state code

**originator.address.countryCode** text

Payment originator address first line

**originator.address.postCode** text

Payment originator address zip code

```json
{
  "id": 16521632,
  "user": 4342275,
  "targetAccount": 8692237,
  "sourceAccount": null,
  "quote": null,
  "quoteUuid": "8fa9be20-ba43-4b15-abbb-9424e1481050",
  "status": "cancelled",
  "reference": "reference text",
  "rate": 0.89,
  "created": "2017-11-24 10:47:49",
  "business": null,
  "transferRequest": null,
  "details": {
    "reference": "Testing"
  },
  "hasActiveIssues": false,
  "sourceCurrency": "EUR",
  "sourceValue": 0,
  "targetCurrency": "GBP",
  "targetValue": 150,
  "customerTransactionId": "54a6bc09-cef9-49a8-9041-f1f0c654cd88",
  "payinSessionId": "23330542-8e9e-419f-95eb-312b880f92ad"
}
```

### Originator transfer

**id** integer, (int64)

Transfer ID

**user** integer, (int64)

Your user ID

**targetAccount** integer, (int64)

Recipient account ID

**sourceAccount** integer, (int64)

Refund recipient account ID

**quote** text

quote ID

**status** text

Transfer current status [For information about transfer statuses](/guides/product/send-money/tracking-transfers)

**reference** text

Deprecated, use details.reference instead

**rate** decimal

Exchange rate value

**created** yyyy-mm-ddThh:mm:ss.sssZ

Timestamp when transfer was created

**business** integer, (int64)

Your business profile ID

**transferRequest** integer

Deprecated

**details.reference** text

Payment reference text

**hasActiveIssues** boolean

Are there any pending issues which stop executing the transfer?

**sourceCurrency** text

Source currency code

**sourceValue** decimal

Transfer amount in source currency

**targetCurrency** text

Target currency code

**targetValue** decimal

Transfer amount in target currency

**originalTransferId** text

Unique identifier randomly generated per transfer request by the calling client

**payinSessionId** text

ID of the Payin Session generated for the transfer, which can be used for certain payin methods when funding the transfer.

**originator** group

Data block to capture payment originator details

**originator.legalEntityType** text

Payment originator legal type

**originator.reference** text

Unique customer ID in your system

**originator.name.givenName** text

Payment originator first name

**originator.name.middleNames** text array

Payment originator middle name(s)

**originator.name.familyName** text

Payment originator family name

**originator.name.patronymicName** text

Payment originator patronymic name

**originator.name.fullName** text

Payment originator full legal name

**originator.dateOfBirth** yyyy-mm-dd

Payment originator date of birth

**originator.businessRegistrationCode** text

Payment originator business registry or incorporation number

**originator.address.firstLine** text

Payment originator address first line

**originator.address.city** text

Payment originator address city

**originator.address.stateCode** text

Payment originator address state code

**originator.address.countryCode** text

Payment originator address first line

**originator.address.postCode** text

Payment originator address zip code

```json
{
  "id": 16521632,
  "user": 4342275,
  "targetAccount": 8692237,
  "sourceAccount": null,
  "quote": "8fa9be20-ba43-4b15-abbb-9424e1481050",
  "status": "cancelled",
  "reference": "reference text",
  "rate": 0.89,
  "created": "2025-10-29T12:28:16.000Z"
  "business": null,
  "transferRequest": null,
  "details": {
    "reference": "Testing"
  },
  "originator": {
    "name": {
      "givenName": "John",
      "middleNames": [
        "Ryan"
      ],
      "familyName": "Godspeed",
      "patronymicName": null,
      "fullName": "John Godspeed"
    },
    "dateOfBirth": "1977-07-01",
    "reference": "CST-2991992",
    "legalEntityType": "PRIVATE",
    "businessRegistrationCode": null,
    "address": {
      "firstLine": "Salu tee 14",
      "city": "Tallinn",
      "stateCode": null,
      "countryCode": "EE",
      "postCode": "12112"
    },
    "accountDetails": "23456789"
  },
  "hasActiveIssues": false,
  "sourceCurrency": "EUR",
  "sourceValue": 0,
  "targetCurrency": "GBP",
  "targetValue": 150,
  "originalTransferId": "54a6bc09-cef9-49a8-9041-f1f0c654cd88",
  "payinSessionId": "23330542-8e9e-419f-95eb-312b880f92ad"
}
```

## Create a transfer 

**`POST /v1/transfers`**

### Request

**sourceAccount (optional)** integer, (int64)

Refund recipient account ID

**targetAccount** integer, (int64)

Recipient account ID. You can create multiple transfers to same recipient account.

**quoteUuid** text

V2 quote ID. You can only create one transfer per one quote. You cannot use same quote ID to create multiple transfers.

**customerTransactionId** uuid

This is required to perform idempotency check to avoid duplicate transfers in case of network failures or timeouts.

**details.reference** text

Recipient will see this reference text in their bank statement. Maximum allowed characters depends on the currency route. [Business Payments Tips](https://wise.com/help/articles/2932870/tips-for-paying-invoices) article has a full list.

**details.transferPurpose (conditionally required)** text

For example when target currency is THB. See more about conditions at [Transfers.Requirements](/api-reference/transfer#transfer-requirements)

**details.transferPurposeSubTransferPurpose (conditionally required)** text

For example when target currency is CNY. See more about conditions at [Transfers.Requirements](/api-reference/transfer#transfer-requirements)

**details.transferPurposeInvoiceNumber (conditionally required)** text

For example when target currency is INR. See more about conditions at [Transfers.Requirements](/api-reference/transfer#transfer-requirements)

**details.sourceOfFunds (conditionally required)** text

For example when target currency is USD and transfer amount exceeds 80k. See more about conditions at [Transfers.Requirements](/api-reference/transfer#transfer-requirements)

There are options to deal with conditionally required fields:

- Always call `transfer-requirements` endpoint and submit values only if indicated so.
- Always provide values for these fields based on a dynamically retrieved list (transfer-requirements endpoint). It is possible these fields change over time, so take this into account if hard coding the options.

Contact api@wise.com if you have questions about this property.

#### Response

Returns a [standard transfer object](/api-reference/transfer#object).

### Avoiding duplicate transfers

We use **customerTransactionId** field to avoid duplicate transfer requests. If your initial call to create a transfer fails (error or timeout) then you should retry the call using the same value in the **customerTransactionId** field that you used in the original call. This way we can treat subsequent retry messages as **repeat messages** and will not create duplicate transfers to your account should one have succeeded before. You should not retry indefinitely but use a sensible limit, perhaps with a back-off approach.

### Payment Approvals

Business Payment Approvals created on your wise.com settings page are not compatible with creating transfers over the API.

If you use personal tokens and do not use client credentials, and if your business account has payment approvals, your application will run in to this error when attempting to create a transfer

`Quote cannot be accepted with this request due to missing approval.`

Consider removing the payment rule if you are going to use the API to create transfers.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/transfers \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "sourceAccount": <refund recipient account ID>,
    "targetAccount": <recipient account ID>,
    "quoteUuid": <quote ID>,
    "customerTransactionId": "<the unique identifier you generated for the transfer attempt>",
    "details" : {
      "reference" : "to my friend",
      "transferPurpose": "verification.transfers.purpose.pay.bills",
      "transferPurposeSubTransferPurpose": "verification.sub.transfers.purpose.pay.interpretation.service",
      "sourceOfFunds": "verification.source.of.funds.other"
    }
  }'
```

## Create a third party transfer 

**`POST /v2/profiles/{{profileId}}/third-party-transfers`**

When creating a transfer on behalf of a third party, you must take note that:

- The `originator` datablock is **required**. This details the ultimate sender of funds in the transfer.
- Depending on the legal entity type of the originator (PRIVATE or BUSINESS), the required fields vary. Please refer the two sample request examples on the right.
- `OriginalTransferId` field must be used. This is your own ID for the transfer.
### Request

**targetAccount** integer, (int64)

Recipient account ID. You can create multiple transfers to same recipient account.

**quote** text

V2 quote ID. You can only create one transfer per one quote. 
You cannot use same quote ID to create multiple transfers.

**originalTransferId** text

Unique transfer ID in your system. We use this field also to perform idempotency check to avoid duplicate transfers in case of network failures or timeouts. You can only submit one transfer with same originalTransferId.

**details.reference (optional)** text

Recipient will see this reference text in their bank statement. Maximum allowed characters depends on the currency route. [Business Payments Tips](https://wise.com/help/articles/2932870/tips-for-paying-invoices) article has a full list.

**originator** group

Data block to capture payment originator details.

**originator.legalEntityType** text

PRIVATE or BUSINESS. Payment originator legal type.

**originator.reference** text

Unique customer ID in your system. This allows us to uniquely identify each originator. Required.

**originator.name** text

Data block to capture the originator name details. 
 Depends on the type of legal entity (PRIVATE or BUSINESS), the required fields and inputs are different.

**originator.name.givenName** text

Payment originator first name. Required if `legalEntityType = PRIVATE`.

**originator.name.middleNames (optional)** text array

Payment originator middle name(s). Used only if `legalEntityType = PRIVATE`.

**originator.name.familyName** text

Payment originator family name. Required if `legalEntityType = PRIVATE`.

**originator.name.patronymicName (optional)** text

Payment originator patronymic name. Used only if `legalEntityType = PRIVATE`.

**originator.name.fullName** text

Payment originator full legal name. Required if `legalEntityType = BUSINESS`.

**originator.dateOfBirth** yyyy-mm-dd

Payment originator date of birth. Required if `legalEntityType = PRIVATE`.

**originator.businessRegistrationCode** text

Payment originator business registry number / incorporation number. Required if `legalEntityType = BUSINESS`.

**originator.address.firstLine** text

Payment originator address first line. Required

**originator.address.city** text

Payment originator address city. Required

**originator.address.stateCode** text

Payment originator address state code. Required if address country code in (US, CA, BR, AU).

**originator.address.countryCode** text

Payment originator address first line. Required

**originator.address.postCode (optional)** text

Originator address zip code.

**originator.accountDetails (optional)** text

Originator account number.

#### Response

Returns an [originator transfer object](/api-reference/transfer#originator-transfer).

You need to save the transfer ID for tracking its status later via webhooks.

### Avoiding duplicate transfers

We use **originalTransferId** field to avoid duplicate transfer requests. When your first call fails (error or timeout) then you should use the same value in **originalTransferId** field that you used in the original call when you are submitting a retry message. This way we can treat subsequent retry messages as **repeat messages** and will not create duplicate transfers to your account.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/third-party-transfers \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "targetAccount": <recipient account ID>,
    "quote": "<quote ID>",
    "originalTransferId": "<unique transfer ID in your system>",
    "details" : {
      "reference" : "Ski trip"
    },
    "originator" : {
      "legalEntityType" : "PRIVATE",
      "reference" : "<unique customer ID in your system>",
      "name" : {
        "givenName": "John",
        "middleNames": ["Ryan"],
        "familyName": "Godspeed"
      },
      "dateOfBirth": "1977-07-01",
      "address" : {
        "firstLine": "Salu tee 100, Apt 4B",
        "city": "Tallinn",
        "countryCode": "EE",
        "postCode": "12112"
      }
    }
  }'
```

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/third-party-transfers \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "targetAccount": <recipient account ID>,
    "quote": "<quote ID>",
    "originalTransferId": "<unique transfer ID in your system>",
    "details" : {
        "reference" : "Payment for invoice 22092"
    },
    "originator" : {
      "legalEntityType" : "BUSINESS",
      "reference" : "<originator customer ID in your system>",
      "name" : {
        "fullName": "Hot Air Balloon Services Ltd"
      },
      "businessRegistrationCode": "1999212",
      "address" : {
        "firstLine": "Aiandi tee 1431",
        "city": "Tallinn",
        "countryCode": "EE",
        "postCode": "12112"
      }
    }
  }'
```

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/third-party-transfers \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "targetAccount": <recipient account ID>,
    "quote": "<quote ID>",
    "originalTransferId": "<unique transfer ID in your system>",
    "details" : {
      "reference" : "Ski trip"
    },
    "originator" : {
      "legalEntityType" : "PRIVATE",
      "reference" : "<unique customer ID in your system>",
      "name" : {
        "givenName": "John",
        "middleNames": ["Ryan"],
        "familyName": "Godspeed"
      },
      "dateOfBirth": "1977-07-01",
      "address" : {
        "firstLine": "102-3393 Capilano Road",
        "city": "North Vancouver",
        "countryCode": "CA",
        "postCode": "V7R 4W7"
      },
      "accountDetails": "<the unique account number of the customer>"
    }
  }'
```

## Create a partner license transfer 

**`POST /v1/profiles/{{profileId}}/partner-licence-transfers`**

This is very similar to **Create transfers** endpoint, but please note these differences:

- originator datablock is additionally required
### Request

**sourceAccount (optional)** integer, (int64)

Refund recipient account ID.

**targetAccount** integer, (int64)

Recipient account ID. You can create multiple transfers to same recipient account.

**quote** text

V2 quote ID. You can only create one transfer per one quote. 
You cannot use same quote ID to create multiple transfers.

**customerTransactionId** text

Unique transfer ID in your system. We use this field also to perform idempotency check to avoid duplicate transfers in case of network failures or timeouts. You can only submit one transfer with same customerTransactionId.

**details.reference (optional)** text

Recipient will see this reference text in their bank statement. Maximum allowed characters depends on the currency route. [Business Payments Tips](https://wise.com/help/articles/2932870/tips-for-paying-invoices) article has a full list.

**originator** group

Data block to capture payment originator details.

**originator.legalEntityType** text

PRIVATE or BUSINESS. Payment originator legal type.

**originator.externalId** text

Unique customer ID in your system. This allows us to uniquely identify each originator. Required.

**originator.name.givenName** text

Payment originator first name. Required if legalEntityType = PRIVATE.

**originator.name.middleNames** text array

Payment originator middle name(s). Used only if legalEntityType = PRIVATE. Optional

**originator.name.familyName** text

Payment originator family name. Required if legalEntityType = PRIVATE.

**originator.name.patronymicName** text

Payment originator patronymic name. Used only if legalEntityType = PRIVATE. Optional

**originator.name.fullName** text

Payment originator full legal name. Required if legalEntityType = BUSINESS.

**originator.dateOfBirth** yyyy-mm-dd

Payment originator date of birth. Required if legalEntityType = PRIVATE.

**originator.businessRegistrationCode** text

Payment originator business registry number / incorporation number. Required if legalEntityType = BUSINESS.

**originator.address.firstLine** text

Payment originator address first line. Required

**originator.address.city** text

Payment originator address city. Required

**originator.address.stateCode** text

Payment originator address state code. Required if address country code in (US, CA, BR, AU).

**originator.address.countryCode** text

Payment originator address first line. Required

**originator.address.postCode** text

Originator address zip code. Optional

#### Response

Returns an [originator transfer object](/api-reference/transfer#originator-transfer).

You need to save the transfer ID for tracking its status later via webhooks.

### Avoiding duplicate transfers

We use **customerTransactionId** field to avoid duplicate transfer requests. When your first call fails (error or timeout) then you should use the same value in **customerTransactionId** field that you used in the original call when you are submitting a retry message. This way we can treat subsequent retry messages as **repeat messages** and will not create duplicate transfers to your account.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/partner-licence-transfers \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "sourceAccount": <refund recipient account ID>,
    "targetAccount": <recipient account ID>,
    "quote": "<quote ID>",
    "customerTransactionId": "<unique transfer ID in your system>",
    "details" : {
      "reference" : "Ski trip"
    },
    "originator" : {
      "legalEntityType" : "PRIVATE",
      "reference" : "<unique customer ID in your system>",
      "name" : {
        "givenName": "John",
        "middleNames": ["Ryan"],
        "familyName": "Godspeed"
      },
      "dateOfBirth": "1977-07-01",
      "address" : {
        "firstLine": "Salu tee 100, Apt 4B",
        "city": "Tallinn",
        "countryCode": "EE",
        "postCode": "12112"
      },
      "accountDetails": "<the unique account number of the customer>"
    }
  }'
```

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/profiles/{{profileId}}/partner-licence-transfers \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "sourceAccount": <refund recipient account ID>,
    "targetAccount": <recipient account ID>,
    "quote": "<quote ID>",
    "customerTransactionId": "<unique transfer ID in your system>",
    "details" : {
      "reference" : "Payment for invoice 22092"
    },
    "originator" : {
      "legalEntityType" : "BUSINESS",
      "reference" : "<originator customer ID in your system>",
      "name" : {
        "fullName": "Hot Air Balloon Services Ltd"
      },
      "businessRegistrationCode": "1999212",
      "address" : {
        "firstLine": "Aiandi tee 1431",
        "city": "Tallinn",
        "countryCode": "EE",
        "postCode": "12112"
      },
      "accountDetails": "<the unique account number of the customer>"
    }
  }'
```

## Requesting transfer requirements 

**`POST /v1/transfer-requirements`**

Almost every region has its own specific nuances when it comes to the nitty gritty details of domestic payment systems and money transfer regulations. The maximum allowed length of reference text is a good example. The US payment system, ACH, supports 10 characters only, but transfers within Mexico allow up to 100 characters, and so on.

The same is true for requirements arising from Anti Money Laundering regulations adopted in different countries. Depending on the chosen currencies and the amount to be transferred, either in one go or cumulatively over time, Wise may require more details about the customer's source of funds or transfer purpose, for example.

The endpoint `/transfer-requirements` exposes all these specific requirements based on the sender, the specific quote and selected target recipient account.

To make sure that processing of your customer's transfers does not get delayed because of missing details, we highly recommend to verify the transfer requirements before submitting any transfer and collecting the data we request from the user using the returned dynamic form.

Transfer Requirement Options**Transfer Reference** string

`reference`- General reference or memo field for a transfer. Mandatory for some transfers, `required` will be true. It is important to use the `minLength`, `maxLength` and `validationRegexp` in your interface to minimize errors.

**Originator Legal Entity Type (optional)** text

`originatorLegalEntityType`- Optional feature flag to identify the type of sender. Allowed values are `PRIVATE` or `BUSINESS`. 
Note this field is **required** from March 2026.

**Transfer Purpose** select

`transferPurpose` - A select list of values to help understand the purpose of the transfer. Often required when transfers are over a certain size or through a particular route. As an option is selected, you must refresh requirements, as additional fields could become required.

**Sub Transfer Purpose** select

`transferPurposeSubTransferPurpose` - A secondary select list of values to further help understand the purpose of the transfer. As an option is selected, you must refresh requirements, as additional fields could become required.

**Source of Funds** select

`sourceOfFunds` - A select list of values to help understand the source of the funds. As an option is selected, you must refresh requirements, as additional fields could become required.

**Transfer Nature** select

`transferNature` - A select list of values to help understand the reason for the transfer, only for routes to and from BRL. Note that transfer nature needs to be added to a quote through creation or update in order to correctly calculate the IOF tax.

**Transfer Purpose Invoice Number** text

`transferPurposeInvoiceNumber`- Mandatory for trade related transfers, `required` will be shown as true for these cases. It is important to use `minLength`, `maxLength` and `validationRegexp` in your interface to minimize errors.

As transfer requirements are part of a dynamic form, the list of options can change. The list above are the most common requirements and when they apply, but others may be included. Always design your integration to work dynamically with this endpoint.

Note: The `originatorLegalEntityType` field is currently optional, but will become **required** from March 2026. Please update your integration to include this field as soon as possible.

Transfers to-BRL or from-BRL require a transfer nature (see: `transferNature`). Transfer nature has some special handling:

- Be aware when rendering this field that labels might be present more than once under the same key.
- The chosen value must be added when creating a quote or via an update of an existing quote.
- Failure to include a transfer nature as instructed will prevent the transfer from being sent.

#### Request

1. 

Prepare the request body to create transfer object first. Now post this request body to the `transfer-requirements` endpoint to figure out if there are any other required fields.

2. 

Analyze the returned list of fields. Our example includes reference, sourceOfFunds and transferPurpose fields. Field 'reference' is optional in this example. Fields 'sourceOfFunds' and 'transferPurpose' are required and both have refreshRequirementsOnChange=true which indicates that there could be additional fields required depending on the selected value.

3. 

In our example you will have to POST request to/v1/transfer-requirements` second time as well with values set for 'transferPurpose' and 'sourceOfFunds'. So in case you set sourceOfFunds = 'verification.source.of.funds.other' then another text field called "sourceOfFundsOther" is also required where you need to specify the details in free format.

4. 

Once you get to the point where you have provided values for all fields which have refreshRequirementsOnChange=true then you have complete set of fields to compose a valid request to create a transfer object. For example this is a valid request to create a transfer.

### Response

**type** text

"transfer"

**fields[n].name** text

Field description

**fields[n].group[n].key** text

Key is name of the field you should include in the JSON

**fields[n].group[n].name** text

Field description

**fields[n].group[n].type** text

Display type of field (e.g. text, select, etc)

**fields[n].group[n].refreshRequirementsOnChange** boolean

Tells you whether you should call POST transfer-requirements once the field value is set to discover required lower level fields.

**fields[n].group[n].required** boolean

Indicates if the field is mandatory or not

**fields[n].group[n].displayFormat** text

Display format pattern.

**fields[n].group[n].example** text

Example value.

**fields[n].group[n].minLength** integer

Min valid length of field value.

**fields[n].group[n].maxLength** integer

Max valid length of field value.

**fields[n].group[n].validationRegexp** text

Regexp validation pattern.

**fields[n].group[n].validationAsync** text

Deprecated. This validation will instead be performed when submitting the request.

**fields[n].group[n].valuesAllowed[n].key** text

List of allowed values. Value key

**fields[n].group[n].valuesAllowed[n].name** text

List of allowed values. Value name.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/transfer-requirements \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{ 
    "targetAccount": <recipient account ID>,
    "quoteUuid": <quote ID>,
    "originatorLegalEntityType": "PRIVATE",
    "details": {
      "reference": "good times",
      "sourceOfFunds": "verification.source.of.funds.other",
      "sourceOfFundsOther": "Trust funds"
    },
    "customerTransactionId": "6D9188CF-FA59-44C3-87A2-4506CE9C1EA3"
  }'
```

```json
[
  {
    "type": "transfer",
    "fields": [
      {
        "name": "Transfer reference",
        "group": [
          {
            "key": "reference",
            "name": "Transfer reference",
            "type": "text",
            "refreshRequirementsOnChange": false,
            "required": false,
            "displayFormat": null,
            "example": null,
            "minLength": null,
            "maxLength": 10,
            "validationRegexp": "[a-zA-Z0-9- ]*",
            "validationAsync": null,
            "valuesAllowed": null
          }
        ]
      },
      {
        "name": "Transfer purpose",
        "group": [
          {
            "key": "transferPurpose",
            "name": "Transfer purpose",
            "type": "select",
            "refreshRequirementsOnChange": true,
            "required": true,
            "displayFormat": null,
            "example": null,
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": [
              {
                "key": "verification.transfers.purpose.purchase.property",
                "name": "Buying property abroad"
              },
              {
                "key": "verification.transfers.purpose.pay.bills",
                "name": "Rent or other property expenses"
              },
              {
                "key": "verification.transfers.purpose.mortgage",
                "name": "Mortgage payment"
              },
              {
                "key": "verification.transfers.purpose.pay.tuition",
                "name": "Tuition fees or studying expenses"
              },
              {
                "key": "verification.transfers.purpose.send.to.family",
                "name": "Sending money home to family"
              },
              {
                "key": "verification.transfers.purpose.living.expenses",
                "name": "General monthly living expenses"
              },
              {
                "key": "verification.transfers.purpose.other",
                "name": "Other"
              }
            ]
          },
          {
            "key": "transferPurposeSubTransferPurpose",
            "name": "Please select a specific reason for your transfer",
            "type": "select",
            "refreshRequirementsOnChange": true,
            "required": true,
            "displayFormat": null,
            "example": null,
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": [
              {
                  "key": "INTERPRETATION_SERVICE",
                  "name": "Interpretation service"
              },
              {
                  "key": "TRANSLATION_SERVICE",
                  "name": "Translation service"
              },
              {
                  "key": "HUMAN_RESOURCE_SERVICE",
                  "name": "Human resource service"
              },
              {
                  "key": "ESTATE_AGENCY_SERVICE",
                  "name": "Estate agency service"
              },
              {
                  "key": "SOFTWARE_DEVELOPMENT_SERVICE",
                  "name": "Software development service"
              },
              {
                  "key": "WEB_DESIGN_OR_DEVELOPMENT_SERVICE",
                  "name": "Web design or development service"
              },
              {
                  "key": "DRAFTING_LEGAL_SERVICE",
                  "name": "Drafting legal service"
              },
              {
                  "key": "LEGAL_RELATED_CERTIFICATION_SERVICE",
                  "name": "Legal related certification service"
              },
              {
                  "key": "ACCOUNTING_SERVICE",
                  "name": "Accounting service"
              },
              {
                  "key": "TAX_SERVICE",
                  "name": "Tax service"
              },
              {
                  "key": "ARCHITECTURAL_DECORATION_DESIGN_SERVICE",
                  "name": "Architectural decoration design service"
              },
              {
                  "key": "ADVERTISING_SERVICE",
                  "name": "Advertising service"
              },
              {
                  "key": "MARKET_RESEARCH_SERVICE",
                  "name": "Market research service"
              },
              {
                  "key": "EXHIBITION_BOOTH_SERVICE",
                  "name": "Exhibition booth service"
              }
            ]
          }
        ]
      },
      {
        "name": "Source of funds",
        "group": [
          {
            "key": "sourceOfFunds",
            "name": "Source of funds",
            "type": "select",
            "refreshRequirementsOnChange": true,
            "required": true,
            "displayFormat": null,
            "example": null,
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": [
              {
                "key": "verification.source.of.funds.salary",
                "name": "Salary"
              },
              {
                "key": "verification.source.of.funds.investment",
                "name": "Investments (stocks, properties, etc.)"
              },
              {
                "key": "verification.source.of.funds.inheritance",
                "name": "Inheritance"
              },
              {
                "key": "verification.source.of.funds.loan",
                "name": "Loan"
              },
              {
                "key": "verification.source.of.funds.other",
                "name": "Other"
              }
            ]
          }
        ]
      },
      {
        "name": "Brazilian regulation requires you to provide a reason behind your transaction.",
        "group": [
          {
            "key": "transferNature",
            "name": "Please select an option that best describes the reason for your transfer",
            "type": "select",
            "refreshRequirementsOnChange": false,
            "required": true,
            "displayFormat": null,
            "example": null,
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": [
              {
                "key": "CHARITABLE_DONATIONS",
                "name": "Donations to friends or family"
              },
              {
                "key": "MOVING_MONEY_BETWEEN_OWN_ACCOUNTS",
                "name": "Moving money between own accounts"
              },
              {
                "key": "INTERNATIONAL_TRAVEL",
                "name": "Tourism service"
              },
              {
                "key": "BUY_OR_SELL_COMPUTER_AND_INFORMATION_SERVICE",
                "name": "Technology service"
              },
              {
                "key": "BUY_OR_SELL_OTHER_SERVICE",
                "name": "Other service"
              },
              {
                "key": "BUY_OR_SELL_MERCHANDISE",
                "name": "Purchase of goods"
              },
              {
                "key": "BUY_OR_SELL_OTHER_SERVICE",
                "name": "Property rental"
              },
              {
                "key": "OTHER",
                "name": "Property purchase or sale"
              },
              {
                "key": "CHARITABLE_DONATIONS",
                "name": "Transfer without compensation"
              }
            ]
          }
        ]
      }
    ]
  }
]
```

## Get a transfer by ID 

**`GET /v1/transfers/{{transferId}}`**

Get transfer info by ID. To receive dynamic updates as the state of the transfer changes, please see our documentation on webhooks.

#### Response

Returns a [transfer object](/api-reference/transfer#object), with or without an `originator` block depending on the type of transfer.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/transfers/{{transferId}} \
  -H 'Authorization: Bearer <your api token>'
```

## Get a third party transfer by ID v2

**`GET /v2/profiles/{{profileId}}/third-party-transfers/{{transferId}}`**

Get transfer info by ID. To receive dynamic updates as the state of the transfer changes, please see [our documentation on webhooks](/api-reference/webhook).

#### Response

Returns an [originator transfer object](/api-reference/transfer#originator-transfer).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/third-party-transfers/{{transferId}} \
  -H 'Authorization: Bearer <your api token>'
```

## Fund a transfer 

**`POST /v3/profiles/{{profileId}}/transfers/{{transferId}}/payments`**

This endpoint is SCA protected when it applies. If your profile is registered within the UK and/or EEA, SCA most likely applies to you. Please read more about [implementing SCA](/guides/developer/auth-and-security/sca-and-2fa). If you use mTLS and client credentials to make API requests, SCA does not apply.

This API call is the final step for executing payouts when using a balance with Wise. Upon calling the endpoint, Wise will begin the processing of the transfer, depending on the status of funds.

When using the transfer by transfer settlement model, the following funding type(s) must be used:

- **BALANCE** - Funds are pulled from a multi-currency account held with Wise.
- **BANK_TRANSFER** - Manually send funds from your business bank account to pay for any transfers. Only applicable when using the [Batch Group API](/api-reference/batch-group/).

When funding through the Bulk Settlement model, the following funding type(s) must be used:

- **TRUSTED_PRE_FUND_BULK** - Funds for the transfer will be settled through a bulk payment at a later date. This method is not applicable to first party Payouts.

If funding from `BALANCE`, and your multi-currency account does not have the required funds to complete the action, then this call will fail with an "insufficient funds" error. Once funds are added and available, you must call this endpoint again.

{{profileId}} refers to the profile that created the transfer. It can be either your personal profile ID, or your business profile ID.

### Request

**type** text

"BALANCE" 
This indicates the type of funding you would like to apply to the transfer.

**partnerReference (conditionally required)** string

The transaction/payment identifier in your system, uniquely identifies the transfer in your platform. This is required for the `Cross Currency Bulk Settlement` model.

### Response

**type** text

"BALANCE" 
This indicates the type of funding you would like to apply to the transfer.

**status** text

"COMPLETED" or "REJECTED"

**errorCode** text

Failure reason. For example "balance.payment-option-unavailable".

```bash
curl -X POST \
  https://api.wise-sandbox.com/v3/profiles/{{profileId}}/transfers/{{transferId}}/payments \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "type": "BALANCE"
  }'
```

```json
{
  "type": "BALANCE",
  "status": "COMPLETED",
  "errorCode": null
}
```

```json
{
  "type": "BALANCE",
  "status": "REJECTED",
  "errorCode": "transfer.insufficient_funds"
}
```

## List transfers for a profile 

**`GET /v1/transfers`**

Get the list of transfers for given user's profile (defaults to user's personal profile).

You can add query parameters to specify user's profile (personal or business), time period and/or payment status.

For example, you can query:

- all failed payments created since last week
- all completed payments created since yesterday
### Request Parameters

**profile** integer

User profile ID. If parameter is omitted, defaults to user's personal profile.

**status** text

Comma separated list of one or more status codes to filter transfers. See [Track transfer status](/guides/product/send-money/tracking-transfers) for complete list of statuses.

**sourceCurrency** text

Source currency code

**targetCurrency** text

Target currency code

**createdDateStart** yyyy-mm-ddThh:mm:ss.sssZ

Starting date to filter transfers, inclusive of the provided date.

**createdDateEnd** yyyy-mm-ddThh:mm:ss.sssZ

Ending date to filter transfers, inclusive of the provided date.

**limit** integer

Maximum number of records to be returned in response

**offset** integer

Starting record number

#### Response

Returns an array of [transfer objects](/api-reference/transfer#object), with or without an `originator` block depending on the type of each transfer.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/transfers
    ?profile={{profileId}}
    &status=funds_refunded
    &offset=0
    &limit=100
    &createdDateStart=2018-12-15
    &createdDateEnd=2018-12-30  \
  -H 'Authorization: Bearer <your api token>'
```

## Cancel a transfer 

**`PUT /v1/transfers/{{transferId}}/cancel`**

Transfers may be cancelled up until the Transfer is sent. Cancellation is final - it can not be undone.

#### Response

Returns a [transfer object](/api-reference/transfer#object), with or without an `originator` block depending on the type of the transfer.

```bash
curl -X PUT \
  https://api.wise-sandbox.com/v1/transfers/{{transferId}}/cancel \
  -H 'Authorization: Bearer <your api token>'
```

## Get a transfer receipt in PDF format 

**`GET /v1/transfers/{{transferId}}/receipt.pdf`**

Download transfer confirmation receipt in PDF format for transfers that are in status **outgoing_payment_sent**. There's also [the transfer state change webhook](/guides/developer/webhooks/event-types#transfer-state-change).

If you service US retail consumers you must use us-combined-receipt instead of this endpoint

#### Response

Transfer confirmation receipt in Wise branded PDF format.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/transfers/{{transferId}}/receipt.pdf \
  -H 'Authorization: Bearer <your api token>' 
```

## Get a transfer US combined receipt in PDF format 

**`GET /v1/transfers/{{transferId}}/us-combined-receipt.pdf`**

Download US combined receipt in PDF format for transfers that are in status **incoming_payment_initiated** and again when the status is updated to **outgoing_payment_sent**.

#### Response

The US Combined Receipt in Wise branded PDF format.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/transfers/{{transferId}}/us-combined-receipt.pdf \
  -H 'Authorization: Bearer <your api token>' 
```

## Get a transfer non objection certificate (INR) 

**`GET /v1/transfers/{{transferId}}/documents/noc`**

Download transfer non objection certificate in PDF format for transfers that are in status **outgoing_payment_sent**.

This document can be used to obtain an Foreign Inward Remittance Certificate (FIRC) from the bank that paid out the transfer.

This is only applicable to INR payments with either a business sender or recipient.

#### Response

Non objection certificate in PDF format.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/transfers/{{transferId}}/documents/noc \
  -H 'Authorization: Bearer <your api token>'  \
  -H 'Accept: application/pdf'
```

## List of completed payments 

**`GET /v1/transfers/{{transferId}}/payments`**

Fetch completed payments used to fund the transfer.

In most cases there should only be a single payment associated with the transfer. There are rare occasions that a transfer can be funded with multiple payment methods and when this occurs the first completed payment method would be used to calculate the fees provided on the quote.

### Response

**id** integer, (int64)

Transfer ID

**method** string

The payment method used to pay for the transfer

**pricingVariant (optional)** string

The qualifier that allows to apply different pricing policy within the same payment method. Often the value might be the payment method itself

**amount** decimal

Transfer source amount

**currency** string

Transfer source currency

**timeCreated** timestamp

Date of when payment was created

**timeUpdated** timestamp

Date of when payment was updated

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/transfers/{{transferId}}/payments \
  -H 'Authorization: Bearer <your api token>' 
```

```json
[
  {
    "id": 50000000,
    "method": "BANK_TRANSFER",
    "pricingVariant": null,
    "amount": 1000.00,
    "currency": "GBP",
    "timeCreated": "2020-01-01T12:30:15.000Z",
    "timeUpdated": "2020-01-01T12:30:15.000Z"
  }
]
```

## Get payout information 

**`GET /v2/transfers/{{transferId}}/invoices/bankingpartner`**

To utilize this endpoint, you will need to replace **transferID** with the specific transfer's unique identifier. The transfer endpoint will return the details of the transfer, including the processorName, deliveryMode, bankingPartnerReference, bankingPartnerName, and mt103. This information will enable your recipients to track the transfer with their bank. It may take up to 3 days to get the correct information through this endpoint, as some partners don't share the information until 3 days later.

Fetch banking reference information for transfers that are in **outgoing_payment_sent** status, enabling you to track transfers with the transfer recipient’s bank.

### Request

**transferID** text

The unique identifier of the transfer.

### Response

**processorName** string

The legal entity that processed the transfer on behalf of the customer.

**deliveryMode** string

The delivery mode for the payment (e.g., Swift).

**bankingPartnerReference** string

The reference used by the partner bank to identify and track the transfer.

**bankingPartnerName** string

The name of the sending bank to the recipient's bank.

**mt103** string

The MT103 of the transfer, if available.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v2/transfers/{{transferId}}/invoices/bankingpartner \
  -H 'Authorization: Bearer <your api token>' 
```

```json
{
  "processorName": "Acme Bank Ltd.",
  "deliveryMode": "SWIFT",
  "bankingPartnerReference": "ABCD1234",
  "bankingPartnerName": "Global Bank Corp."
  "mt103": "{1:F01XXXXGBXXAXXX0000000000}{2:I103XXXXGBXXXXXXN}{3:{108:1234567}{111:001}{121:00000000-0000-0000-0000-000000000000}}{4:\n:20:1234567\n:23B:CRED\n:32A:221212USD12345,\n:33B:USD12345,\n:50K:/11111111\nSOME COMPANY INC.\n1 SOME STREET MIAMI 33132 US\n:59:/GB00000000000000\nCOMPANY NAME LTD\nUK LONDON 1234 GB\n:70:REFERENCE\n:71A:OUR\n:71G:USD11,\n-}\n"
}
```