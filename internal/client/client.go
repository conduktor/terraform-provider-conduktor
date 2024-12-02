package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	ctlresource "github.com/conduktor/ctl/resource"
	ctlschema "github.com/conduktor/ctl/schema"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

type ConsoleClient struct {
	apiKey  string
	baseUrl string
	client  *resty.Client
}

type GatewayClient struct {
	GatewayUser     string
	GatewayPassword string
	baseUrl         string
	client          *resty.Client
}

type ApiParameter struct {
	ApiKey        string
	BaseUrl       string
	CdkUser       string
	CdkPassword   string
	TLSParameters TLSParameters
}

type GatewayApiParameters struct {
	BaseUrl         string
	GatewayUser     string
	GatewayPassword string
	TLSParameters   TLSParameters
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

func Make(ctx context.Context, apiParameter ApiParameter, providerVersion string) (*ConsoleClient, error) {
	restyClient := resty.New().SetHeader("X-CDK-CLIENT", "TF/"+providerVersion)

	// Enable http client debug logs when provider log is set to TRACE
	restyClient.SetDebug(TraceLogEnabled())

	if (apiParameter.CdkUser != "" && apiParameter.CdkPassword == "") || (apiParameter.CdkUser == "" && apiParameter.CdkPassword != "") {
		return nil, fmt.Errorf("CDK_USER and CDK_PASSWORD must be provided together")
	}
	if apiParameter.CdkUser != "" && apiParameter.ApiKey != "" {
		return nil, fmt.Errorf("Can't set both CDK_USER and CDK_API_KEY")
	}

	restyClient, err := ConfigureTLS(ctx, restyClient, apiParameter.TLSParameters)
	if err != nil {
		return nil, err
	}

	result := &ConsoleClient{
		apiKey:  apiParameter.ApiKey,
		baseUrl: uniformizeBaseUrl(apiParameter.BaseUrl),
		client:  restyClient,
	}

	if apiParameter.CdkUser != "" {
		retry3Time := retry(3, 1*time.Second, ctx)
		loginResult, err := retry3Time(
			func(err error) bool {
				return err.Error() == "Invalid username or password"
			}, func() (interface{}, error) {
				return result.login(apiParameter.CdkUser, apiParameter.CdkPassword)
			},
		)

		if err != nil {
			return nil, fmt.Errorf("Could not login: %s", err)
		}
		tokens, _ := loginResult.(LoginResult)

		result.apiKey = tokens.AccessToken
	}

	if result.apiKey != "" {
		result.client = result.client.SetAuthScheme("Bearer")
		result.client = result.client.SetAuthToken(result.apiKey)
	}

	return result, nil
}

func MakeGateway(ctx context.Context, apiParameter GatewayApiParameters, providerVersion string) (*GatewayClient, error) {
	restyClient := resty.New().SetHeader("X-CDK-CLIENT", "TF/"+providerVersion)

	// Enable http client debug logs when provider log is set to TRACE
	restyClient.SetDebug(TraceLogEnabled())

	restyClient, err := ConfigureTLS(ctx, restyClient, apiParameter.TLSParameters)
	if err != nil {
		return nil, err
	}

	restyClient.SetBasicAuth(apiParameter.GatewayUser, apiParameter.GatewayPassword)

	// Testing authentication parameters against /metrics API
	// returning error after 3 retries
	testUrl := apiParameter.BaseUrl + "/metrics"
	resp, err := restyClient.SetRetryCount(3).SetRetryWaitTime(1 * time.Second).R().Get(testUrl)
	if err != nil {
		return &GatewayClient{}, err
	} else if resp.StatusCode() != 200 {
		return &GatewayClient{}, fmt.Errorf("Invalid username or password")
	}

	return &GatewayClient{
		GatewayUser:     apiParameter.GatewayUser,
		GatewayPassword: apiParameter.GatewayPassword,
		baseUrl:         apiParameter.BaseUrl,
		client:          restyClient,
	}, nil
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

func (client *ConsoleClient) login(username, password string) (LoginResult, error) {
	url := client.baseUrl + "/login"
	resp, err := client.client.R().SetBody(map[string]string{"username": username, "password": password}).Post(url)
	if err != nil {
		return LoginResult{}, err
	} else if resp.IsError() {
		if resp.StatusCode() == 401 {
			return LoginResult{}, fmt.Errorf("Invalid username or password")
		} else {
			return LoginResult{}, fmt.Errorf("%s", extractApiError(resp))
		}
	}
	result := LoginResult{}
	err = jsoniter.Unmarshal(resp.Body(), &result)
	if err != nil {
		return LoginResult{}, err
	}
	return result, nil
}

func (client *ConsoleClient) ApplyGeneric(ctx context.Context, cliResource ctlresource.Resource) (string, error) {
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

func (client *ConsoleClient) Apply(ctx context.Context, path string, resource interface{}) (ApplyResult, error) {
	url := client.baseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("Error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("PUT %s request body : %s", path, string(jsonData)))
	builder := client.client.R().SetBody(jsonData)
	resp, err := builder.Put(url)
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

func (client *ConsoleClient) Describe(ctx context.Context, path string) ([]byte, error) {
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

func (client *ConsoleClient) Delete(ctx context.Context, path string) error {
	url := client.baseUrl + path
	tflog.Trace(ctx, fmt.Sprintf("DELETE %s", path))
	resp, err := client.client.R().Delete(url)
	if err != nil {
		return err
	} else if resp.IsError() {
		return fmt.Errorf("%s", extractApiError(resp))
	}

	return nil
}

// This is a temporary workaround - will be revisited with the future client works
func (client *GatewayClient) Apply(ctx context.Context, path string, resource interface{}) (ApplyResult, error) {
	url := client.baseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("Error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("PUT %s request body : %s", path, string(jsonData)))
	builder := client.client.R().SetBody(jsonData)
	resp, err := builder.Put(url)
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

func (client *GatewayClient) Describe(ctx context.Context, path string) ([]byte, error) {
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

func (client *GatewayClient) Delete(ctx context.Context, path string, resource interface{}) error {
	url := client.baseUrl + path

	jsonData, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("Error marshalling resource: %s", err)
	}
	tflog.Debug(ctx, string(jsonData))

	tflog.Trace(ctx, fmt.Sprintf("PUT %s request body : %s", path, string(jsonData)))
	builder := client.client.R().SetBody(string(jsonData))
	tflog.Trace(ctx, fmt.Sprintf("DELETE %s", path))
	resp, err := builder.Delete(url)
	if err != nil {
		return err
	} else if resp.IsError() {
		return fmt.Errorf("%s", extractApiError(resp))
	}

	return nil
}
