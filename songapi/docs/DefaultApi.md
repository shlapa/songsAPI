# \DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InfoGet**](DefaultApi.md#InfoGet) | **Get** /info | 



## InfoGet

> SongDetail InfoGet(ctx).Group(group).Song(song).Execute()



### Example

```go
package main

import (
	openapiclient "./openapi"
	"context"
	"fmt"
	"os"
)

func main() {
	group := "group_example" // string | 
	song := "song_example"   // string | 

	configuration := openapiclient.NewConfiguration()
	api_client := openapiclient.NewAPIClient(configuration)
	resp, r, err := api_client.DefaultApi.InfoGet(context.Background()).Group(group).Song(song).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.InfoGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InfoGet`: SongDetail
	fmt.Fprintf(os.Stdout, "Response from `DefaultApi.InfoGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInfoGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **group** | **string** |  | 
 **song** | **string** |  | 

### Return type

[**SongDetail**](SongDetail.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

