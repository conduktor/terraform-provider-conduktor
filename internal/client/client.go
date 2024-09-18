package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	ctlresource "github.com/conduktor/ctl/resource"
	ctlschema "github.com/conduktor/ctl/schema"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsoniter "github.com/json-iterator/go"
)

type Client struct {
	apiKey  string
	baseUrl string
	client  *resty.Client
}

type ApiParameter struct {
	ApiKey      string
	BaseUrl     string
	Key         string
	Cert        string
	Cacert      string
	CdkUser     string
	CdkPassword string
	Insecure    bool
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type ApplyResult struct {
	UpsertResult string
}

func Make(ctx context.Context, apiParameter ApiParameter, providerVersion string) (*Client, error) {
	restyClient := resty.New().SetHeader("X-CDK-CLIENT", "TF/"+providerVersion)

	// Enable http client debug logs when provider log is set to TRACE
	restyClient.SetDebug(TraceLogEnabled())

	if apiParameter.BaseUrl == "" {
		return nil, fmt.Errorf("Please set api base url")
	}

	if (apiParameter.Key == "" && apiParameter.Cert != "") || (apiParameter.Key != "" && apiParameter.Cert == "") {
		return nil, fmt.Errorf("Certificate and Key must be provided together")
	} else if apiParameter.Key != "" && apiParameter.Cert != "" {
		tflog.Debug(ctx, fmt.Sprintf("Loading certificate and key from files %s and %s", apiParameter.Cert, apiParameter.Key))
		certificate, err := tls.LoadX509KeyPair(apiParameter.Cert, apiParameter.Key)
		restyClient.SetCertificates(certificate)
		if err != nil {
			return nil, err
		}
	}

	if (apiParameter.CdkUser != "" && apiParameter.CdkPassword == "") || (apiParameter.CdkUser == "" && apiParameter.CdkPassword != "") {
		return nil, fmt.Errorf("CDK_USER and CDK_PASSWORD must be provided together")
	}
	if apiParameter.CdkUser != "" && apiParameter.ApiKey != "" {
		return nil, fmt.Errorf("Can't set both CDK_USER and CDK_API_KEY")
	}

	if apiParameter.Cacert != "" {
		tflog.Debug(ctx, fmt.Sprintf("Loading cacert from file %s", apiParameter.Cacert))
		restyClient.SetRootCertificate(apiParameter.Cacert)
	}

	if apiParameter.Insecure {
		tflog.Debug(ctx, "Insecure mode enabled (skipping TLS verification)")
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	result := &Client{
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
		if err != nil {
			return nil, fmt.Errorf("Could not login: %s", err)
		}

		result.apiKey = tokens.AccessToken
	}

	if result.apiKey != "" {
		result.client = result.client.SetAuthScheme("Bearer")
		result.client = result.client.SetAuthToken(result.apiKey)
	}

	return result, nil
}

func (client *Client) login(username, password string) (LoginResult, error) {
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

func (client *Client) Apply(ctx context.Context, path string, resource interface{}) (string, error) {
	url := client.baseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return "", fmt.Errorf("Error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("PUT on %s body : %s", path, string(jsonData)))
	builder := client.client.R().SetBody(jsonData)
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

func (client *Client) Delete(ctx context.Context, path string) error {
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
