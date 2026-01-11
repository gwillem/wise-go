<!-- Source: https://docs.wise.com/api-reference/currencies -->

# Currencies
Operations 

[GET/v1/currencies](/api-reference/currencies#get)

Get all currencies allowed for transfers

## The Currencies resource 

**code** text

Currency code (ISO 4217 Alphabetic Code)

**symbol** text

The symbol of this currency

**name** text

Display name of this currency

**countryKeywords** list of strings

Keywords associated with this currency

**supportsDecimals** boolean

Whether this currency supports decimal values or not

```json

[
  {
    "code": "AUD",
    "symbol": "A$",
    "name": "Australian dollar",
    "countryKeywords": [
      "AUD",
      "AU",
      "Australia",
      "aus"
    ],
    "supportsDecimals": true
  },
  {
    "code": "JPY",
    "symbol": "¥",
    "name": "Japanese yen",
    "countryKeywords": [
      "JPY",
      "JP",
      "Japan",
      "jpn"
    ],
    "supportsDecimals": false
  },
  ...
]
```

## Get all currencies allowed for transfers 

**`GET /v1/currencies`**

Get the list of allowed currencies that you can use when setting up your transfers.

#### Response

Returns a [currencies object](/api-reference/currencies#object).

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/currencies \
  -H 'Authorization: Bearer <your api token>'
```