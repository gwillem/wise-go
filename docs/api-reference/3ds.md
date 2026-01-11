<!-- Source: https://docs.wise.com/api-reference/3ds -->

# 3D Secure Authentication
To manage certain aspects of the 3D Secure (3DS) authentication, you will need to integrate with the following APIs.

Operations 

[POST/v3/spend/profiles/{{profileId}}/3dsecure/challenge-result](/api-reference/3ds#inform-challenge-result)

Inform challenge result

## Inform Challenge Result

**`POST /v3/spend/profiles/{{profileId}}/3dsecure/challenge-result`**

Once the customer has accepted or rejected the push notification for a 3DS challenge, you can use this endpoint to notify us of the result.

You must call this endpoint before the expiration time, otherwise it will return a 400 error. You can find the expiration information from the 

[webhook event](/guides/developer/webhooks/event-types#3ds-challenge)

Only the first call to this endpoint will be processed, any subsequent duplicate requests will be ignored, although you will still receive a success response.

This endpoint is SCA protected when it applies. If your profile is registered within the UK and/or EEA, SCA most likely applies to you. For more information please read more [implementing SCA](/guides/developer/auth-and-security/sca-and-2fa).

### Request

**transactionReference** text

Transaction reference as received in the [webhook event](/guides/developer/webhooks/event-types#3ds-challenge)

**challengeStatus** text

One of `APPROVED` or `REJECTED`

```bash
curl -X POST \
  https://api.wise-sandbox.com/v3/spend/profiles/{{profile_id}}/3dsecure/challenge-result \
  -H 'Authorization: Bearer <your api token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "transactionReference": "148579538",
    "challengeStatus": "APPROVED"
  }'
```

#### Response

204 - No Content