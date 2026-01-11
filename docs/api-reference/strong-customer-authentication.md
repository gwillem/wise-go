<!-- Source: https://docs.wise.com/api-reference/strong-customer-authentication -->

# Strong Customer Authentication
Strong Customer Authentication (SCA) is a regulatory requirement introduced under the Second Payment Services Directive (PSD2) in the European Economic Area (EEA). It aims to enhance the security and reduce fraud risks.

These endpoints allow your customers to comply with SCA, ensuring seamless integration over API while meeting security requirements.

Please read this [guide](/guides/developer/auth-and-security/sca-over-api) to understand how SCA integration works.

SCA is required for high-risk and low-risk operations. On low-risk operations, clearing SCA once allows a 5 minutes window where SCA won't be requested again.

Operations 

[GET/v1/one-time-token/status](/api-reference/strong-customer-authentication#one-time-token-status)

Get One-Time Token status

[POST/v2/profiles/{{profileId}}/pin](/api-reference/strong-customer-authentication#create-pin)

Create a PIN

[POST/v2/profiles/{{profileId}}/pin/verify](/api-reference/strong-customer-authentication#verify-pin)

Verify a PIN

[DELETE/v2/profiles/{{profileId}}/pin](/api-reference/strong-customer-authentication#delete-pin)

Delete a PIN

[POST/v2/profiles/{{profileId}}/device-fingerprints](/api-reference/strong-customer-authentication#create-device-fingerprint)

Create a device fingerprint

[POST/v2/profiles/{{profileId}}/device-fingerprints/verify](/api-reference/strong-customer-authentication#verify-device-fingerprint)

Verify a device fingerprint

[DELETE/v2/profiles/{{profileId}}/device-fingerprints/{{deviceFingerprintId}}](/api-reference/strong-customer-authentication#delete-device-fingerprint)

Delete a device fingerprint

[POST/v2/profiles/{{profileId}}/facemaps](/api-reference/strong-customer-authentication#create-facemap)

Create a facemap

[POST/v2/profiles/{{profileId}}/facemaps/verify](/api-reference/strong-customer-authentication#verify-facemap)

Verify a facemap

[DELETE/v2/profiles/{{profileId}}/facemaps](/api-reference/strong-customer-authentication#delete-facemap)

Delete a facemap

[POST/v2/profiles/{{profileId}}/sca-sessions/authorise](/api-reference/strong-customer-authentication#create-sca-session)

Create a SCA session

## The One-Time Token resource 

A One-Time Token is generated when accessing an endpoint secured by SCA. This token includes a list of all available challenges to complete. To view the challenges and their statuses linked to the token, please use the 

[status endpoint](/api-reference/strong-customer-authentication#one-time-token-status).

Alternatively, you can use [create SCA session](/api-reference/strong-customer-authentication#create-sca-session) to manually trigger SCA and return a One-Time Token.

At least two challenges must be completed to access an SCA-secured endpoint. For low-risk operations, access to these endpoints is valid for 5 minutes before the One-Time Token expires.

Please reach out to your Implementation Manager to determine which challenges can be enabledFields**oneTimeToken** text

A One-Time Token unique identifier.

**challenges** object[]

An array of challenges.

**challenges.primaryChallenge.type** text

A type of challenge.

**challenges.passed** boolean

The status of a challenge.

**validity** number

The One-Time Token expiration in seconds.

```json
{
  "oneTimeToken": "5932d5b5-ec13-452f-8688-308feade7834",
  "challenges": [
    {
      "primaryChallenge": {
        "type": "PIN",
      },
      "passed": false
    }
  ],
  "validity": 3600
}
```

## Get One-Time Token status 

### Request

Set the `One-Time Token` property in the request headers to return challenges associated to it.

### Response

Returns a [One-Time Token](/api-reference/strong-customer-authentication#one-time-token-object)

```bash
curl -X GET \
  'https://api.wise-sandbox.com/v1/one-time-token/status' \
  -H 'Authorization: Bearer <your api token>' \
  -H 'One-Time-Token: <one time token>'
```

## Create PIN 

Creates a new PIN factor used to resolve a SCA `knowledge challenge` type.

Note that the request and response are encrypted using JOSE framework. Please refer to this [guide](/guides/developer/auth-and-security/sca-over-api) to understand how encryption and decryption work.

### Request

**pin** text

A four digits string.

### Response

**200 - OK**

The PIN has been successfully created.

**409 - Conflict**

A PIN has already been created for this profile.

```bash
curl -X POST \
  'https://api.wise-sandbox.com/v2/profiles/{{profileId}}/pin' \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Accept: application/jose+json' \
  -H 'Accept-Encoding: identity' \
  -H 'Content-Type: application/jose+json' \
  -H 'Content-Encoding: identity' \
  -H 'X-tw-jose-method: jwe' \
  -d '
  {
    "pin" : "1234"
  }'
```

## Verify PIN 

Verifies a PIN challenge when calling a SCA-secured endpoint. Please make sure to [create a PIN](/api-reference/strong-customer-authentication#create-pin) before using this endpoint.

Note that the request and response are encrypted using JOSE framework. Please refer to this [guide](/guides/developer/auth-and-security/sca-over-api) to understand how encryption and decryption work.

### Request

Header**One-Time-Token** text

A One-Time Token unique identifier.

Body**pin** text

A 4 digits string.

### Response

Returns an [One-Time Token resource](/api-reference/strong-customer-authentication#one-time-token-object)

**200 - OK**

The PIN has been successfully verified.

**400 - Bad request**

The PIN verification failed.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/pin/verify \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Accept: application/jose+json' \
  -H 'Accept-Encoding: identity' \
  -H 'Content-Type: application/jose+json' \
  -H 'Content-Encoding: identity' \
  -H 'X-tw-jose-method: jwe' \
  -H 'One-Time-Token: <one time token>' \
  -d '
  {
    "pin" : "1234"
  }'
```

## Delete PIN 

Deletes a PIN associated to a profile.

To update a PIN for a profile, please use this endpoint followed by [create a PIN](/api-reference/strong-customer-authentication#create-pin).

Please be aware that this operation is irreversible.

### Response

**204 - No content**

The PIN has been deleted.

**404 - Not found**

No PIN has been setup for this profile.

```bash
curl -X DELETE \
  'https://api.wise-sandbox.com/v2/profiles/{{profileId}}/pin' \
  -H 'Authorization: Bearer <your api token>'
```

## Create Device Fingerprint 

Creates a new device fingerprint factor used to resolve a SCA `possession challenge` type.

Note that the request and response are encrypted using JOSE framework. Please refer to this [guide](/guides/developer/auth-and-security/sca-over-api) to understand how encryption and decryption work.

### Request

**deviceFingerprint** text

A string value used as a device fingerprint.

### Response

Body**deviceFingerprintId** UUID

The identifier of the device fingerprint

**createdAt** date

The device fingerprint creation timestamp

HTTP code**200 - HTTP OK**

The device fingerprint has been successfully created.

**409 - Conflict**

The device fingerprint has already been created.

**400 - Bad Request**

Maximum number of device fingerprints reached (default is 3).

```bash
curl -X POST \
  'https://api.wise-sandbox.com/v2/profiles/{{profileId}}/device-fingerprints' \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Accept: application/jose+json' \
  -H 'Accept-Encoding: identity' \
  -H 'Content-Type: application/jose+json' \
  -H 'Content-Encoding: identity' \
  -H 'X-tw-jose-method: jwe' \
  -d '
  {
    "deviceFingerprint": "3207da22-a0d3-4b6b-a591-6297e646fe32" 
  }'
```

```json
{
  "deviceFingerprintId": "636a5514-aa86-4719-8700-e9a9a0ae7ea7",
  "createdAt": "2025-05-24T07:27:58.273205554Z"
}
```

## Verify device fingerprint 

Verifies a device fingerprint challenge when calling a SCA-secured endpoint. Please make sure to [create a device fingerprint](/api-reference/strong-customer-authentication#create-device-fingerprint) before using this endpoint.

Note that the request and response are encrypted using JOSE framework. Please refer to this [guide](/guides/developer/auth-and-security/sca-over-api) to understand how encryption and decryption work.

### Request

Header**One-Time-Token** text

A One-Time Token unique identifier.

Body**deviceFingerprint** text

A device fingerprint value.

### Response

Returns an [One-Time Token resource](/api-reference/strong-customer-authentication#one-time-token-object)

**200 - OK**

The device fingerprint has been successfully verified.

**400 - Bad request**

The device fingerprint verification failed.

```bash
curl -X POST \
  'https://api.wise-sandbox.com/v2/profiles/{{profileId}}/device-fingerprints/verify' \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Accept: application/jose+json' \
  -H 'Accept-Encoding: identity' \
  -H 'Content-Type: application/jose+json' \
  -H 'Content-Encoding: identity' \
  -H 'X-tw-jose-method: jwe' \
  -H 'One-Time-Token: <one time token>' \
  -d '
  {
    "deviceFingerprint": "3207da22-a0d3-4b6b-a591-6297e646fe32" 
  }'
```

## Delete device fingerprint 

Deletes a device fingerprint associated to a profile. Include the `deviceFingerprintId` in the URL to delete a device fingerprint. This ID is provided in the response when the device fingerprint is created.

Please be aware that this operation is irreversible.

### Response

**204 - No content**

The device fingerprint has been deleted.

**404 - Not found**

The `deviceFingerprintId` does not exist.

```bash
curl -X DELETE \
  'https://api.wise-sandbox.com/v2/profiles/{{profileId}}/device-fingerprints/{{deviceFingerprintId}}' \
  -H 'Authorization: Bearer <your api token>' \
```

## Create facemap 

Creates a new facemap factor used to resolve a SCA `inherence challenge` type.

A facemap should be exported from your FaceTec server using the SDK's [export API](https://dev.facetec.com/api-guide#export-3d-facemap). Please use Wise's FaceTec [public key](/api-reference/facetec#public-key) to encrypt the facemap during the export process.

### Request

**faceMap** text

A base64 encoded string.

### Response

**204 - No content**

The facemap has been successfully created.

**409 - Conflict**

A facemap has already been created for this profile.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/facemaps \
  -H 'Authorization: Bearer <your api token>' \
  -d '{
    "faceMap": "<base64 encrypted facemap>"
  }'
```

## Verify facemap 

Verifies a facemap challenge when calling a SCA-secured endpoint. Please make sure to [create a facemap](/api-reference/strong-customer-authentication#create-facemap) before using this endpoint.

A facemap should be exported from your FaceTec server using the SDK's [export API](https://dev.facetec.com/api-guide#export-3d-facemap). Please use Wise's FaceTec [public key](/api-reference/facetec#public-key) to encrypt a facemap during the export process.

### Request

Header**One-Time-Token** text

A One-Time Token unique identifier.

Body**faceMap** text

A base64 encoded string.

### Response

Returns an [One-Time Token resource](/api-reference/strong-customer-authentication#one-time-token-object)

**200 - OK**

The facemap has been successfully verified.

**400 - Bad request**

The facemap verification failed.

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/facemaps/verify \
  -H 'Authorization: Bearer <your api token>' \
  -H 'One-Time-Token: <one time token>' \
  -d '{
    "faceMap": "<base64 encoded string>"
  }'
```

## Delete facemap 

Deletes a facemap associated to a profile.

To update a facemap for a profile, please use this endpoint followed by [create a facemap](/api-reference/strong-customer-authentication#create-facemap).

Please be aware that this operation is irreversible.

### Response

**204 - No content**

The facemap has been deleted.

**404 - Not found**

No facemap has been setup for this profile.

```bash
curl -X DELETE \
  'https://api.wise-sandbox.com/v2/profiles/{{profileId}}/facemaps' \
  -H 'Authorization: Bearer <your api token>'
```

## Create SCA session 

SCA can be triggered manually allowing more control when integrating with our APIs. The endpoint returns a [One-Time Token](/api-reference/strong-customer-authentication#one-time-token-object) along with a list of associated challenges. These challenges can be cleared with verify endpoints.

### Response

Returns a list of challenges to clear SCA.

**oneTimeTokenProperties** OneTimeToken

Properties of [OneTimeToken](/api-reference/strong-customer-authentication#one-time-token-object)

```bash
curl -X POST \
  https://api.wise-sandbox.com/v2/profiles/{{profileId}}/sca-sessions/authorise \
  -H 'Authorization: Bearer <your api token>'
```

```json
{
  "oneTimeTokenProperties": {
    "oneTimeToken": "9f5f5812-2609-4e48-8418-b64437c0c7cd",
    "challenges": [
      {
        "primaryChallenge": {
          "type": "PIN",
        },
        "passed": false
      }
    ],
    "validity": 3600,
  }
}
```