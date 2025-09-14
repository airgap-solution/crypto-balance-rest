# DefaultApi

All URIs are relative to *http://restartfu.com:8082*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**balanceGet**](#balanceget) | **GET** /balance | Get XPUB balance in BTC and fiat|

# **balanceGet**
> BalanceGet200Response balanceGet()

Derives addresses from an XPUB, queries Electrum for balances, and converts to the requested fiat currency using the crypto rate API. 

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from '@airgap-solution/crypto-balance-rest-client';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let xpub: string; //Extended public key (XPUB) to derive addresses from. (default to undefined)
let currency: string; //Crypto currency symbol (must be supported by the rate API). (default to undefined)
let fiat: string; //Fiat currency symbol. (default to undefined)

const { status, data } = await apiInstance.balanceGet(
    xpub,
    currency,
    fiat
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **xpub** | [**string**] | Extended public key (XPUB) to derive addresses from. | defaults to undefined|
| **currency** | [**string**] | Crypto currency symbol (must be supported by the rate API). | defaults to undefined|
| **fiat** | [**string**] | Fiat currency symbol. | defaults to undefined|


### Return type

**BalanceGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Balance successfully retrieved |  -  |
|**400** | Invalid request (bad xpub or parameters) |  -  |
|**500** | Internal server error (Electrum or rate API issue) |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

