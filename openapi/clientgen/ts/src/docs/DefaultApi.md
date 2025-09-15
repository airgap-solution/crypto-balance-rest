# DefaultApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**balanceGet**](#balanceget) | **GET** /balance | Get account balance in crypto and fiat|

# **balanceGet**
> BalanceGet200Response balanceGet()

Retrieves the balance for an extended public key (XPUB or equivalent), calculates its value in the requested fiat currency, and returns the current balance along with the 24-hour change in fiat value. 

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from '@airgap-solution/crypto-balance-rest-client';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let xpub: string; //Extended public key to derive addresses from. (default to undefined)
let currency: string; //Cryptocurrency symbol. (default to undefined)
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
| **xpub** | [**string**] | Extended public key to derive addresses from. | defaults to undefined|
| **currency** | [**string**] | Cryptocurrency symbol. | defaults to undefined|
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
|**400** | Invalid request (bad XPUB or parameters) |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

