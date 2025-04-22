package main
import (
        "context"
        "fmt"
        "net/http"
        "crypto/tls"
        "github.com/jean1/terraform-provider-netbox-dns/client"
)

/*
Example API client to NetBox DNS Plugin
*/

func doPlainReq(ctx context.Context, req *http.Request, c *client.Client) (*http.Response, error) {
	req = req.WithContext(ctx)
	for _, e := range c.RequestEditors {
		if err := e(ctx, req); err != nil {
			return nil, err
		}
	}

	return c.Client.Do(req)
}

func apiKeyAuth(token string) client.RequestEditorFn {
        return func(ctx context.Context, req *http.Request) error {
                req.Header.Set("Authorization", "Token "+token)
                return nil
        }
}

func main () {

	serverUrl := "http://nb/"
	token := "a339fbe313c1c183e7896490a9778a4981c90202"

	ctx := context.Background()
	// init api client
	opts := []client.ClientOption{
		client.WithRequestEditorFn(apiKeyAuth(token)), // auth
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	opts = append(opts, client.WithHTTPClient(httpClient))
	myclient, err := client.NewClient(serverUrl, opts...)
	if err != nil {
		fmt.Println("failed to create client: %v", err.Error())
		return
	}

	// do request
	var params client.PluginsNetboxDnsZonesListParams
	zonenames := []string{"u-strasbg.fr"}
	params.Name = &zonenames
	viewnames := []string{"_default_"}
	params.View = &viewnames
	
	nextHTTPReq, err := client.NewPluginsNetboxDnsZonesListRequest(myclient.Server, &params)
	if err != nil {
		fmt.Println("Client Error", fmt.Sprintf("failed to create zone list request: %s", err))
		return
	}
	var httpRes *http.Response
	httpRes, err = doPlainReq(ctx, nextHTTPReq, myclient)
	if err != nil {
		fmt.Println("Client Error", fmt.Sprintf("failed to retrieve Zones: %s", err))
		return
	}
	var res *client.PluginsNetboxDnsZonesListResponse
	res, err = client.ParsePluginsNetboxDnsZonesListResponse(httpRes)
	if err != nil {
		fmt.Println("Client Error", fmt.Sprintf("failed to parse Zones: %s", err))
		return
	}
	// fmt.Println("res.Body: %v", res.Body)
	//fmt.Println("res.httpresp: %v", res.HTTPResponse)
	if res.JSON200 != nil {
		var found = res.JSON200.Results
		if len(found) >= 1 {
			fmt.Printf("zone id=%d\n", *(found[0].Id))
		} else {
			fmt.Println("zone not found")
		}
	}

	// fmt.Println("r: %v", r)
}
