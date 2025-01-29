package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	ctlresource "github.com/conduktor/ctl/resource"
	ctlschema "github.com/conduktor/ctl/schema"
	gateway "github.com/conduktor/terraform-provider-conduktor/internal/model/gateway"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

// Enum used to perform different actions based on provider mode.
// Setting as string so it can be used for log messages.
type Mode string

const (
	CONSOLE Mode = "Console"
	GATEWAY Mode = "Gateway"
)

type Client struct {
	baseUrl string
	client  *resty.Client
}

type ApiParameter struct {
	ApiKey        string
	BaseUrl       string
	CdkUser       string
	CdkPassword   string
	TLSParameters TLSParameters
}

type TLSParameters struct {
	Key      string
	Cert     string
	Cacert   string
	Insecure bool
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type ApplyResult struct {
	UpsertResult string      `json:"upsertResult"`
	Resource     interface{} `json:"resource"`
}

func Make(ctx context.Context, mode Mode, apiParameter ApiParameter, providerVersion string) (*Client, error) {
	restyClient := resty.New().SetHeader("X-CDK-CLIENT", "TF/"+providerVersion)
	if mode == CONSOLE {
		apiParameter.BaseUrl = uniformizeBaseUrl(apiParameter.BaseUrl)
	}
	var err error

	// Enable http client debug logs when provider log is set to TRACE
	restyClient.SetDebug(TraceLogEnabled())

	restyClient, err = ConfigureTLS(ctx, restyClient, apiParameter.TLSParameters)
	if err != nil {
		return nil, err
	}

	restyClient, err = ConfigureAuth(ctx, mode, restyClient, apiParameter)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseUrl: apiParameter.BaseUrl,
		client:  restyClient,
	}, nil
}

func ConfigureAuth(ctx context.Context, mode Mode, restyClient *resty.Client, apiParameter ApiParameter) (*resty.Client, error) {
	var err error
	switch mode {
	case CONSOLE:
		{
			apiKey := apiParameter.ApiKey
			if apiKey == "" {
				// Only Login with username and password if no apiKey has been provided.
				apiKey, err = Login(apiParameter, restyClient)
				if err != nil {
					return nil, fmt.Errorf("Could not login: %s", err)
				}
			}

			restyClient = restyClient.SetAuthScheme("Bearer")
			restyClient = restyClient.SetAuthToken(apiKey)
		}
	case GATEWAY:
		{
			restyClient.SetBasicAuth(apiParameter.CdkUser, apiParameter.CdkPassword)

			// Testing authentication parameters against /metrics API.
			// Returning error after 3 retries.
			testUrl := apiParameter.BaseUrl + "/metrics"
			resp, err := restyClient.SetRetryCount(3).SetRetryWaitTime(1 * time.Second).R().Get(testUrl)
			if err != nil {
				return nil, err
			} else if resp.StatusCode() != 200 {
				return nil, fmt.Errorf("Invalid username or password")
			}
		}
	}

	return restyClient, nil
}

// Helper function for Console Auth flow to retrieve access token.
func Login(apiParameter ApiParameter, client *resty.Client) (string, error) {
	url := apiParameter.BaseUrl + "/login"
	body := map[string]string{
		"username": apiParameter.CdkUser,
		"password": apiParameter.CdkPassword,
	}

	resp, err := client.SetRetryCount(3).SetRetryWaitTime(1 * time.Second).R().SetBody(body).Post(url)
	if err != nil {
		return "", err
	} else if resp.IsError() {
		if resp.StatusCode() == 401 {
			return "", fmt.Errorf("Invalid username or password")
		} else {
			return "", fmt.Errorf("%s", extractApiError(resp))
		}
	}
	result := LoginResult{}
	err = jsoniter.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func ConfigureTLS(ctx context.Context, restyClient *resty.Client, tlsParameter TLSParameters) (*resty.Client, error) {
	if (tlsParameter.Key == "" && tlsParameter.Cert != "") || (tlsParameter.Key != "" && tlsParameter.Cert == "") {
		return nil, fmt.Errorf("Certificate and Key must be provided together")
	} else if tlsParameter.Key != "" && tlsParameter.Cert != "" {
		tflog.Debug(ctx, fmt.Sprintf("Loading certificate and key from files %s and %s", tlsParameter.Cert, tlsParameter.Key))
		certificate, err := tls.LoadX509KeyPair(tlsParameter.Cert, tlsParameter.Key)
		restyClient.SetCertificates(certificate)
		if err != nil {
			return nil, err
		}
	}

	if tlsParameter.Cacert != "" {
		tflog.Debug(ctx, fmt.Sprintf("Loading cacert from file %s", tlsParameter.Cacert))
		restyClient.SetRootCertificate(tlsParameter.Cacert)
	}

	if tlsParameter.Insecure {
		tflog.Debug(ctx, "Insecure mode enabled (skipping TLS verification)")
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return restyClient, nil
}

func (client *Client) ApplyGeneric(ctx context.Context, cliResource ctlresource.Resource) (string, error) {
	kinds := ctlschema.ConsoleDefaultKind() // TODO support gateway kind and client too

	kindName := cliResource.Kind
	kind, ok := kinds[kindName]
	if !ok {
		return "", fmt.Errorf("Apply kind %s not found", cliResource.Kind)
	}

	applyPath, err := kind.ApplyPath(&cliResource)
	if err != nil {
		return "", err
	}

	url := client.baseUrl + applyPath

	tflog.Trace(ctx, fmt.Sprintf("PUT on %s body : %s", applyPath, string(cliResource.Json)))
	builder := client.client.R().SetBody(cliResource.Json)
	resp, err := builder.Put(url)
	if err != nil {
		return "", err
	} else if resp.IsError() {
		return "", fmt.Errorf("%s", extractApiError(resp))
	}
	bodyBytes := resp.Body()
	var upsertResponse ApplyResult
	err = jsoniter.Unmarshal(bodyBytes, &upsertResponse)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling response: %s", err)
	}
	return upsertResponse.UpsertResult, nil
}

func (client *Client) Apply(ctx context.Context, path string, resource interface{}) (ApplyResult, error) {
	url := client.baseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("Error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("PUT %s request body : %s", path, string(jsonData)))

	resp, err := client.client.R().SetBody(jsonData).Put(url)
	if err != nil {
		return ApplyResult{}, err
	} else if resp.IsError() {
		return ApplyResult{}, fmt.Errorf("%s", extractApiError(resp))
	}

	bodyBytes := resp.Body()
	tflog.Trace(ctx, fmt.Sprintf("PUT %s response body : %s", path, string(bodyBytes)))

	var upsertResponse ApplyResult
	err = jsoniter.Unmarshal(bodyBytes, &upsertResponse)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("Error unmarshalling response: %s", err)
	}
	return upsertResponse, nil
}

func (client *Client) Describe(ctx context.Context, path string) ([]byte, error) {
	url := client.baseUrl + path
	resp, err := client.client.R().Get(url)
	if err != nil {
		return []byte{}, err
	} else if resp.IsError() {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return []byte{}, fmt.Errorf("error describing resources %s, got status code: %d:\n %s", path, resp.StatusCode(), string(resp.Body()))
	}
	tflog.Trace(ctx, fmt.Sprintf("GET %s response : %s", path, string(resp.Body())))
	return resp.Body(), nil
}

func (client *Client) Delete(ctx context.Context, mode Mode, path string, resource interface{}) error {
	var req *resty.Request
	url := client.baseUrl + path
	tflog.Trace(ctx, fmt.Sprintf("DELETE %s", path))

	if mode == CONSOLE {
		req = client.client.R()
	} else if mode == GATEWAY {
		// Gateway API handles deletion in a different way
		// It needs information about the resource in the body of the request
		// as opposed to Console API who needs them in the URL
		jsonData, err := json.Marshal(resource)
		if err != nil {
			return fmt.Errorf("Error marshalling resource: %s", err)
		}
		tflog.Debug(ctx, string(jsonData))
		tflog.Trace(ctx, fmt.Sprintf("DELETE %s request body : %s", path, string(jsonData)))

		req = client.client.R().SetBody(string(jsonData))
	}

	resp, err := req.Delete(url)
	if err != nil {
		return err
	} else if resp.IsError() {
		return fmt.Errorf("%s", extractApiError(resp))
	}

	return nil
}

func (client *Client) ApplyGatewayToken(ctx context.Context, path string, resource interface{}) (ApplyResult, error) {
	url := client.baseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("Error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("POST %s request body : %s", path, string(jsonData)))

	resp, err := client.client.R().SetBody(jsonData).Post(url)
	if err != nil {
		return ApplyResult{}, err
	} else if resp.IsError() {
		return ApplyResult{}, fmt.Errorf("%s", extractApiError(resp))
	}

	bodyBytes := resp.Body()
	tflog.Trace(ctx, fmt.Sprintf("POST %s response body : %s", path, string(bodyBytes)))

	var upsertResponse gateway.GatewayTokenResource
	err = jsoniter.Unmarshal(bodyBytes, &upsertResponse)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("Error unmarshalling response: %s", err)
	}
	return ApplyResult{Resource: upsertResponse}, nil
}
