<!-- Source: https://docs.wise.com/api-reference/client-credentials-token -->

# Client Credentials Token
Operations 

[POST/oauth/token](/api-reference/client-credentials-token#retrieve-client-creds)

Retrieve a client credentials token

## The Client Credentials Token resource 

**access_token** text

Access token to be used when calling "create user" endpoint. Valid for 12 hours.

**token_type** text

Type (`bearer`)

**expires_in** integer

Expiry time in seconds

**scope** text

Scope (`transfers`)

```json
{
  "access_token": "ba8k1234-00f2-475a-60d8-6g45377b4062",
  "token_type":"bearer",
  "expires_in": 43199,
  "expires_at": "2025-04-11T03:43:28.148Z",
  "scope":"transfers"
}
```

## Retrieve a client credentials token

**`POST /oauth/token`**

Obtain access_token based on your API client credentials.

Request

Use Basic Authentication with your api-client-id/api-client-secret as the username/password. The body of the request must be sent as x-www-form-urlencoded.

**grant_type** text

"client_credentials"

#### Response

Returns a 

[client credentials token object](/api-reference/client-credentials-token#object)

```bash
curl \
  https://api.wise-sandbox.com/oauth/token \
  -u '<client id>:<client secret>' \
  -d 'grant_type=client_credentials'
```