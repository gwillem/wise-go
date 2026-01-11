<!-- Source: https://docs.wise.com/api-reference/contact -->

# Contact
A contact is a representation of a recipient or beneficiary.

Operations 

[POST/v2/profiles/{{profileId}}/contacts?isDirectIdentifierCreation=true](/api-reference/contact#create)

Create a contact from an identifier

## Create a contact from an identifier 

**`POST /v2/profiles/{{profileId}}/contacts?isDirectIdentifierCreation=true`**

This endpoint finds an existing **discoverable** Wise profile and adds it to your recipient list. The recipient is found just by an identifier, without having to know their bank details. The contactId from the response can be used to create a transfer.

### Request

**identifier** text

A Wise profile’s identifier. It can be a Wisetag, email or phone number.

**targetCurrency** text

3 character currency code

### Response

**contactId** text

contactId in the format of UUID

**name** text

full name of the contact

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/contacts?isDirectIdentifierCreation=true \
  -H 'Authorization: Bearer <client credential token>' \
  -H 'Content-type: application/json' \
  -d '{
    "identifier":"@JonathanP",
    "targetCurrency":"EUR"
  }'

```

### Creating a transfer with a contact

1. Create a 

[Quote](/api-reference/quote#create-authenticated) where the **targetAccount** field is replaced with "contactId": `<contactId>`
2. Extract the **targetAccount** value from the Quote
3. Create a [Transfer](/api-reference/transfer#create)