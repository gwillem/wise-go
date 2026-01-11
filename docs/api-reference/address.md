<!-- Source: https://docs.wise.com/api-reference/address -->

# Address
The Address resource represents physical addresses for a profile. Once an address is created, it can be added to a profile.

Operations 

[POST/v1/addresses](/api-reference/address#create-or-update-an-address)

Create or update an address

[GET/v1/addresses/{{addressId}}](/api-reference/address#retrieve-an-address-by-id)

Retrieve an address by ID

[GET/v1/addresses](/api-reference/address#list-addresses-for-a-profile)

List addresses for a profile

[POST/v1/address-requirements](/api-reference/address#retrieve-address-requirements-dynamically)

Retrieve address requirements dynamically

## The Address resource

Address Object**id** integer

Address ID

**profile** integer

User profile ID

**details.country** text

Country code (ISO 3166-2 Country Code)

**details.firstLine** text

Address line: street, house, apartment

**details.postCode** text

Postal / zip code

**details.city** text

City name

**details.state** text

State code | Required if country is US, CA, BR, AU

**details.occupations

[n].code** text

User occupation - any value permitted

**details.occupations[n].format** text

Occupation type - always `FREE_FORM`

```json
{
  "id": 10000001,
  "profile": {{profile ID}}, 
  "details": {
    "country" : "US",
    "state": "AZ",
    "city": "Phoenix",
    "postCode": "10025",
    "firstLine": "50 Sunflower Ave",
    "occupations": [
      {
        "code": "Software Engineer",
        "format": "FREE_FORM"
      }
    ]
  }
}
```

## Create or update an address

**`POST /v1/addresses`**

Adds address info to user profile.

List of required fields are different for different countries.

For updating personal profiles, consider using [the personal profile update endpoint](/api-reference/profile#update-a-personal-profile) instead. It allows to submit the address information alongside the other profile data.

State field is required for US, CA, BR and AU addresses.

Occupations is required for CA, IN, JP, ID, IL, MX and within the US for the state NM.

### Request

**profile** integer

User profile ID

**details.country** text

Country code (ISO 3166-2 Country Code)

**details.firstLine** text

Address line: street, house, apartment

**details.postCode** text (max 30 chars)

Zip code

**details.city** text

City name

**details.state** text

State code | Required if country is US, CA, BR, AU

**details.occupations[n].code** text

User occupation - any value permitted

**details.occupations[n].format** text

Occupation type - always `FREE_FORM`

#### Response

Returns an [address object](#the-address-resource).

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/addresses \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "profile": {{profile ID}},
    "details": {
      "country" : "US",
      "state": "AZ",
      "city": "Phoenix",
      "postCode": "10025",
      "firstLine": "50 Sunflower Ave",
      "occupations": [
        {
          "code": "Software Engineer",
          "format": "FREE_FORM"
        }
      ]
    }
  }'
```

## Retrieve an address by ID

**`GET /v1/addresses/{{addressId}}`**

Get address info by ID.

Returns an [address object](#the-address-resource).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/addresses/{{addressId}} \
  -H 'Authorization: Bearer <your api token>'
```

## List addresses for a profile

**`GET /v1/addresses?profile={{profileId}}`**

List of addresses belonging to user profile.

Returns an array of [address objects](#the-address-resource).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/addresses?profile={{profileId}} \
  -H 'Authorization: Bearer <your api token>'
```

## Retrieve address requirements dynamically

**`GET /v1/address-requirements`** 
 **`POST /v1/address-requirements`** 

GET and POST address-requirements endpoints help you to figure out which fields are required to create a valid address for different countries. You could even build a dynamic user interface on top of these endpoints. This is a step-by-step guide on how these endpoints work.

1. Call `GET /v1/address-requirements` to get list of fields you need to fill with values in "details" section for creating a valid address. Response contains 4 required top level fields:

- country (select field with list of values)
- city (text field)
- postCode (text field)
- firstLine (text field)

1. 

Analyze the list of fields. Because `refreshRequirementsOnChange` for field 'country' is marked as `true` then this indicates that there are additional fields required depending on the selected value.

2. 

Call `POST /v1/address-requirements` with selected country value to expose sub-fields. 
For example posting `{"details": {"country" : "US"}}` will also add "state" to list of fields.
 But posting `{"details": {"country" : "GB"}}` will not.

3. 

If you choose "US" as country you will notice that "state" field also has `refreshRequirementsOnChange=true`. This means you would need to make another POST call to `/v1/address-requirements` with a specific state value.
 For example posting `{"details": { "country" : "US", "state": "NM" }}` will also add "occupation" to list of fields.
 But posting `{"details": { "country" : "US", "state": "AL" }}` will not.

4. 

So once you get to the point where you have provided values for all fields which have refreshRequirementsOnChange=true then you have complete set of fields to compose a valid request to create an address object.

### Response

**type** text

"address"

**fields[n].name** text

Field description

**fields[n].group[n].key** text

Key is name of the field you should include in the JSON

**fields[n].group[n].type** text

Display type of field (e.g. text, select, etc)

**fields[n].group[n].refreshRequirementsOnChange** boolean

Tells you whether you should call POST address-requirements once the field value is set to discover required lower level fields

**fields[n].group[n].required** boolean

Indicates if the field is mandatory or not

**fields[n].group[n].displayFormat** text

Display format pattern

**fields[n].group[n].example** text

Example value to help customers understand what to input

**fields[n].group[n].minLength** integer

Min valid length of field value

**fields[n].group[n].maxLength** integer

Max valid length of field value

**fields[n].group[n].validationRegexp** text

Regexp validation pattern

**fields[n].group[n].validationAsync** text

Deprecated. This validation will instead be performed when submitting the request.

**fields[n].group[n].valuesAllowed[n].key** text

List of allowed value's key

**fields[n].group[n].valuesAllowed[n].name** text

List of allowed value's name

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/address-requirements \
  -H 'Authorization: Bearer <your api token>'
     
curl -X POST \
  https://api.wise-sandbox.com/v1/address-requirements \
  -H 'Authorization: Bearer <your api token>' \
  -d '{
    "profile" : {{profile ID}},
    "details": {
      "country" : "US",
      "state": "AZ",
      "city": "Phoenix",
      "postCode": "10025",
      "firstLine": "50 Sunflower Ave",
      "occupations": [
        {
          "code": "Software Engineer",
          "format": "FREE_FORM"
        }
      ]
    }
  }' 
```

```json
[
  {
    "type": "address",
    "fields": [
      {
        "name": "Country",
        "group": [
          {
            "key": "country",
            "type": "select",
            "refreshRequirementsOnChange": true,
            "required": true,
            "displayFormat": null,
            "example": "Germany",
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": [
              {
                "key": "AX",
                "name": "Åland Islands"
              },
              ...
              {
                "key": "ZM",
                "name": "Zambia"
              }
            ]
          }
        ]
      },
      {
        "name": "City",
        "group": [
          {
            "key": "city",
            "type": "text",
            "refreshRequirementsOnChange": false,
            "required": true,
            "displayFormat": null,
            "example": "London",
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": null
          }
        ]
      },
      {
        "name": "Postal code",
        "group": [
          {
            "key": "postCode",
            "type": "text",
            "refreshRequirementsOnChange": false,
            "required": true,
            "displayFormat": null,
            "example": "10025",
            "minLength": null,
            "maxLength": null,
            "validationRegexp": null,
            "validationAsync": null,
            "valuesAllowed": null
          }
        ]
      }
      ...
    ]
  }
]
```