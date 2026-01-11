<!-- Source: https://docs.wise.com/api-reference/facetec -->

# FaceTec
Wise leverages [FaceTec's](https://www.facetec.com/) facial biometric technology for authentication, offering a seamless integration experience through these APIs.

Operations 

[GET/v1/facetec/public-key](/api-reference/facetec#public-key)

Get Public Key

## Get Public Key 

**`GET /v1/facetec/public-key`**

Retrieve Wise's FaceTec public key to be used when 

[exporting 3D Facemap](https://dev.facetec.com/api-guide#export-3d-facemap) from your FaceTec host to Wise.

The exported FaceMap can be used to [Enrol FaceMap](/api-reference/user-security#enrol-facemap).

#### Response

Plain text containing public key.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/facetec/public-key \
  -H 'Authorization: Bearer <your api token>' 
```

```
-----BEGIN PUBLIC KEY-----
Public Key Content
-----END PUBLIC KEY----- 
```