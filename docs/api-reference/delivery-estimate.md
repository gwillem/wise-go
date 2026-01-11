<!-- Source: https://docs.wise.com/api-reference/delivery-estimate -->

# Delivery Estimate
Operations 

[GET/v1/delivery-estimates/{{transferId}}](/api-reference/delivery-estimate#get)Get the delivery estimate for a transfer

## Get the delivery estimate for a transfer 

**`GET /v1/delivery-estimates/{{transferId}}?timezone=Asia/Singapore`**

Get the live delivery estimate for a transfer by the transfer ID.

The delivery estimate is the time at which we currently expect the transfer to arrive in the beneficiary's bank account.

This is not a guaranteed time but we are working hard to make these estimates as accurate as possible.

### Request

**timezone** text

Timezone ID for the formatted text. Example: `UTC`, `Asia/Singapore`. Defaults to `UTC` if not provided.

### Response

**estimatedDeliveryDate** timestamp

Estimated time when funds will arrive to recipient's bank account

**formattedEstimatedDeliveryDate** text

A string to display to users for the estimated delivery date.

```bash
curl -X GET \
  https://api.wise-sandbox.com/v1/delivery-estimates/{{transferId}}?timezone=Asia/Singapore \
  -H 'Authorization: Bearer <your api token>' 
```

```json
{
  "estimatedDeliveryDate" : "2018-01-10T12:15:00.000+0000",
  "formattedEstimatedDeliveryDate" : "in seconds"
}
```