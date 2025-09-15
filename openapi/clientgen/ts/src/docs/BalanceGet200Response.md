# BalanceGet200Response


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**value** | **number** | Balance converted to the requested fiat currency. | [optional] [default to undefined]
**balance** | **number** | Balance in the requested cryptocurrency. | [optional] [default to undefined]
**change24h** | **number** | Change in fiat value over the past 24 hours (positive &#x3D; gain, negative &#x3D; loss). | [optional] [default to undefined]

## Example

```typescript
import { BalanceGet200Response } from '@airgap-solution/crypto-balance-rest-client';

const instance: BalanceGet200Response = {
    value,
    balance,
    change24h,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
