<!-- Source: https://docs.wise.com/api-reference/user -->

# User
In our API, a User serves as the primary entity and can possess multiple Profiles to represent different contexts or settings. Specifically, a User can have one personal Profile and multiple business Profiles. Each [Profile](/api-reference/profile) - whether personal or business - can have its own [multi-currency account](/api-reference/multi-currency-account), enabling transactions across various currencies. This hierarchical structure allows for flexible management of user settings and financial operations, accommodating both personal and business needs.

UserPersonal ProfileBusiness ProfileMulti Currency AccountMulti Currency Account

Operations 

[GET/v1/me](/api-reference/user#retrieve-by-token)

Retrieve current user by token

[GET/v1/users/{{userId}}](/api-reference/user#retrieve-by-id)

Retrieve a user by ID

[POST/v1/user/signup/registration_code](/api-reference/user#createuser-regcode)

Create a user with a registration code

[POST/v1/users/exists](/api-reference/user#checkuserexist)

Check to see if customer exist with queried email address

[PUT/v1/users/{{userId}}/contact-email](/api-reference/user#set-contact-email)

Set the contact email address

[GET/v1/users/{{userId}}/contact-email](/api-reference/user#get-contact-email)

Retrieve the contact email address

## The User resource 

**id** integer

userId

**name** text

User's full name

**email** text

User's email

**active** boolean

If user is active or not

**details.firstName** text

User's first name

**details.lastName** text

User's lastname

**details.phoneNumber** text

Phone number

**details.dateOfBirth** YYYY-MM-DD

Date of birth

**details.occupation** text

Person's occupation

**details.avatar** text

Link to person avatar image.

**details.primaryAddress** integer

Address object ID to use in addresses endpoints

**details.address.countryCode** text

Address country code in 2 digits. "US" for example

**details.address.firstLine** text

Address first line

**details.address.postCode** text

Address post code

**details.address.city** text

Address city name

**details.address.state** text

Address state code State code. Required if country is US, CA, AU, BR.

**details.address.occupation** text

User occupation. Required for US, CA, JP

```json
{
  "id": 101,
  "name": "Example Person",
  "email": "person@example.com",
  "active": true,
  "details": {
    "firstName": "Example",
    "lastName": "Person",
    "phoneNumber": "+37111111111",
    "occupation": "",
    "address": {
      "city": "Tallinn",
      "countryCode": "EE",
      "postCode": "11111",
      "state": "",
      "firstLine": "Road 123"
    },
    "dateOfBirth": "1977-01-01",
    "avatar": "https://lh6.googleusercontent.com/photo.jpg",
    "primaryAddress": 111
  }
}
```

## Retrieve current user by token

**`GET /v1/me`**

Get authenticated user details for the user's token submitted in the Authorization header. Response includes also personal user's profile info.

#### Response

Returns a 

[user object](/api-reference/user#object)

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/me \
  -H 'Authorization: Bearer <your api token>' 
```

## Retrieve a user by Id

**`GET /v1/users/{{userId}}`**

Get authenticated user details by user ID. Response includes also personal user's profile info.

#### Response

Returns a [user object](/api-reference/user#object)

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/users/{{userId}} \
  -H 'Authorization: Bearer <your api token>' 
```

## Create a user with a registration code

**`POST /v1/user/signup/registration_code`**

Wise uses email address as unique identifier for users. If email is new (there is no active user already) then new user will be created.

When you are submitting an email which already exists amongst our users then you will get a warning that "You’re already a member. Please login". If user already exists then you need to redirect to "Get user authorization" webpage.

### Request

**email** email

New user's email address

**registrationCode** text, min length is 32 chars

Randomly generated registration code that is unique to this user and request. At least 32 characters long. You need to store registration code to obtain access token on behalf of this newly created user in next step. Please apply the same security standards to handling registration code as if it was a password.

**language (Optional)** text, 2 chars

User default language for UI and email communication. Allowed values EN, US, PT, ES, FR, DE, IT, JA, RU, PL, HU, TR, RO, NL, HK. Default value EN.

### Response

**id** integer

userId

**name** text

User full name. Empty.

**email** text

Customer email

**active** boolean

true

**details** group

User details. Empty.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/user/signup/registration_code \
  -H 'Authorization: Bearer <client credentials token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "user@email.com",
    "registrationCode": <unique guid for the user>,
    "language": "EN"
  }'
```

```json
{
  "id": 12345,
  "name": null,
  "email": "new.user@domain.com",
  "active": true,
  "details": null
}
```

```json
{
  "errors": [
    {
      "code": "NOT_UNIQUE",
      "message": "You’re already a member. Please login",
      "path": "email",
      "arguments": [
        "email",
        "class com.transferwise.fx.api.ApiRegisterCommand",
        "existing.user@domain.com"
      ]
    }
  ]
}
```

## Check User Exists 

**`POST /v1/users/exists`**

Wise uses email address as unique identifier for users. If email has already been used by a user, it cannot be reused to create a new user.

Note that this uses a `client-credentials-token` and not a `user access_token` for authentication.

### Request

**email** email

User's email address

### Response

**exists** boolean

Email has already exist

```bash
curl -X POST \
  https://api.wise-sandbox.com/v1/users/exists \
  -H 'Authorization: Bearer <client credentials token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "test@wise.com"
  }'
```

```json
{
  "exists": true
}
```

## Set a contact email address 

**`PUT /v1/users/{{userId}}/contact-email`**

Sets a contact email address. The contact email address is used to send notifications to users who have been registered with a dummy email address.

### Request

**email** email

Contact email address

### Response

**email** email

Contact email address

```bash
curl -X PUT \
  https://api.wise-sandbox.com/v1/users/{{userId}}/contact-email \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "new-user@email.com"
  }'
```

```json
{
  "email": "new-user@email.com"
}
```

## Retrieve a contact email address 

**`GET /v1/users/{{userId}}/contact-email`**

Retrieves a contact email address.

#### Response

Returns a contact email object.

### Response

**email** email

Contact email address

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/users/{{userId}}/contact-email \
  -H 'Authorization: Bearer <your api token>' 
```

```json
{
  "email": "new-user@email.com"
}
```