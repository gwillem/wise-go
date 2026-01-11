<!-- Source: https://docs.wise.com/api-reference/user-tokens -->

# User Tokens
User tokens allow your system to make API calls on behalf of a Wise user.

*Access tokens* are short-lived API tokens used to access Wise customer API resources.

*Refresh tokens* are long-lived API tokens that are used to generate access tokens.

Operations 

[POST/oauth/token (registration code)](/api-reference/user-tokens#regcode)

Retrieve user tokens with registration code

[POST/oauth/token (authorization code)](/api-reference/user-tokens#authzcode)

Retrieve user tokens with authorization code

[POST/oauth/token (refresh token)](/api-reference/user-tokens#refresh)

Retrieve user tokens with refresh token

## User Tokens resource 

**access_token** text

Access token to be used when calling API endpoints on behalf of user. Valid for 12 hours.

**token_type** text

"bearer"

**refresh_token** text

Refresh token which you need to use in order to request new access_token. The lifetime of refresh tokens is 20 years.

**expires_in** integer

Access Token expiry time in seconds

**expires_at** text

Access Token expiration timestamp (UTC)

**refresh_token_expires_in** integer

Refresh Token expiry time in seconds

**refresh_token_expires_at** text

Refresh Token expiration timestamp (UTC)

**scope** text

"transfers"

**created_at** datetime

Creation time in ISO 8601 format

```json
{
  "access_token": "01234567-89ab-cdef-0123-456789abcdef",
  "token_type": "bearer",
  "refresh_token": "01234567-89ab-cdef-0123-456789abcdef",
  "expires_in": 43199,
  "expires_at": "2025-04-11T03:43:28.148Z",
  "refresh_token_expires_in": 628639555,
  "refresh_token_expires_at": "2045-03-12T13:49:23.552Z",
  "scope": "transfers",
  "created_at": "2020-01-01T12:33:33.12345Z"
}
```

## Retrieve user tokens with registration code

**`POST /oauth/token`**

You can now use registration code to obtain user access token and refresh token.

### Request

**grant_type** text

"registration_code"

**email** email

New user's email address

**client_id** text

Your API client_id

**registration_code** text

registrationCode

#### Response

Returns a 

[user tokens object](/api-reference/user-tokens#object)

```bash
curl \
  https://api.wise-sandbox.com/oauth/token \
  -u '<client id>:<client secret>' \
  -d 'grant_type=registration_code' \
  -d 'client_id=<client id>' \
  -d 'email=<user email>' \
  -d 'registration_code=<registration code used to create user>'
```

```json
{
  "error": "invalid_grant",
  "error_description": "Invalid user credentials."
}
```

## Retrieve user tokens with authorization code

**`POST /oauth/token`**

You can now use authorization code to obtain user access token and refresh token.

### Request

**grant_type** text

"authorization_code"

**client_id** text

Your API client_id

**code** text

Authorization code provided to you upon redirect back from the authorization flow.

**redirect_uri** text

Redirect URL associated with your API client credentials.

#### Response

Returns a [user tokens object](/api-reference/user-tokens#object)

```bash
curl \
  https://api.wise-sandbox.com/oauth/token \
  -u '<client id>:<client secret>' \
  -d 'grant_type=authorization_code' \
  -d 'client_id=<client id>' \
  -d 'code=<code from redirect uri>' \
  -d 'redirect_uri=https://www.yourapp.com'
```

```json
{
  "access_token": {access-token},
  "token_type": "bearer",
  "refresh_token": {refresh-token},
  "expires_in": 43199,
  "expires_at": "2025-04-11T03:43:28.148Z",
  "refresh_token_expires_in": 628639555,
  "refresh_token_expires_at": "2045-03-12T13:49:23.552Z",
  "scope": "transfers",
  "created_at": "2023-12-06T18:28:14.206824830Z"
}
```

```json
{
  "error": "invalid_request",
  "error_description": "Missing grant type"
}
```

## Retrieve user tokens with refresh token

**`POST /oauth/token`**

Access tokens are valid for 12 hours, so upon expiry you need to use the refresh token to generate a new access token.

In order to maintain an uninterrupted connection, you can request a new access token whenever it’s close to expiring. There is no need to wait for the actual expiration to happen first.

Depending on how your application uses the Wise Platform API, you may find that requesting a new access token before attempting a series of API calls on behalf of an individual user will avoid issues with expired access tokens.

### Request

**grant_type** text

"refresh_token"

**refresh_token** text

User's refresh token obtained from creating or linking to a Wise user.

### Response

Returns a [user tokens object](/api-reference/user-tokens#object)

```bash
curl \
  https://api.wise-sandbox.com/oauth/token \
  -u '<client id>:<client secret>' \
  -d 'grant_type=refresh_token' \
  -d 'refresh_token=<user refresh token>'
```