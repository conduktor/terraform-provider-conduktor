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

// Enum used to perform different actions based on provider mode.
// Setting as string so it can be used for log messages.
type Mode string

const (
	CONSOLE Mode = "Console"
	GATEWAY Mode = "Gateway"
)

type Client struct {
	BaseUrl string
	Client  *resty.Client
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type ApplyResult struct {
	UpsertResult string `json:"upsertResult"`
	Resource     any    `json:"resource"`
}

func Make(ctx context.Context, mode Mode, apiParameter ApiParameter, providerVersion string) (*Client, error) {
	restyClient := resty.New().SetHeader("X-CDK-CLIENT", "TF/"+providerVersion)
	if mode == CONSOLE {
		apiParameter.BaseUrl = uniformizeBaseUrl(apiParameter.BaseUrl)
	}
	var err error

	// Enable http client debug logs when TF_LOG_PROVIDER_CONDUKTOR_INIT is set to trace
	restyClient.SetDebug(InitTraceEnabled())

	restyClient, err = ConfigureTLS(ctx, restyClient, apiParameter.TLSParameters)
	if err != nil {
		return nil, err
	}

	restyClient, err = ConfigureAuth(mode, restyClient, apiParameter)
	if err != nil {
		return nil, err
	}

	// Enable http client debug logs when provider log is set to TRACE
	restyClient.SetDebug(TraceLogEnabled())

	return &Client{
		BaseUrl: apiParameter.BaseUrl,
		Client:  restyClient,
	}, nil
}

func ConfigureAuth(mode Mode, restyClient *resty.Client, apiParameter ApiParameter) (*resty.Client, error) {
	var err error
	switch mode {
	case CONSOLE:
		{
			apiKey := apiParameter.ApiKey
			if apiKey == "" {
				// Only Login with username and password if no apiKey has been provided.
				apiKey, err = Login(apiParameter, restyClient)
				if err != nil {
					return nil, fmt.Errorf("could not login: %s", err)
				}
			}

			restyClient = restyClient.SetAuthScheme("Bearer")
			restyClient = restyClient.SetAuthToken(apiKey)
		}
	case GATEWAY:
		{
			restyClient.SetDisableWarn(true)
			restyClient.SetBasicAuth(apiParameter.CdkUser, apiParameter.CdkPassword)

			// Testing authentication parameters against /metrics API.
			// Returning error after 3 retries.
			testUrl := apiParameter.BaseUrl + "/metrics"
			resp, err := restyClient.SetRetryCount(3).SetRetryWaitTime(1 * time.Second).R().Get(testUrl)
			if err != nil {
				return nil, err
			} else if resp.IsError() {
				switch resp.StatusCode() {
				case 401:
					return nil, fmt.Errorf("invalid username or password")
				case 403:
					return nil, fmt.Errorf("forbidden: You do not have permission to access admin API. Please check your user permissions")
				case 500:
					return nil, fmt.Errorf("internal Server Error: Please check the server logs for more details")
				default:
					return nil, fmt.Errorf("unexpected response (%d): %s", resp.StatusCode(), resp.String())
				}
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
			return "", fmt.Errorf("invalid username or password")
		} else {
			return "", fmt.Errorf("%s", ExtractApiError(resp))
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
		return nil, fmt.Errorf("certificate and Key must be provided together")
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
	catalog := ctlschema.ConsoleDefaultCatalog() // TODO support gateway kind and client too
	kinds := catalog.Kind
	kindName := cliResource.Kind
	kind, ok := kinds[kindName]
	if !ok {
		return "", fmt.Errorf("apply kind %s not found", cliResource.Kind)
	}

	applyPath, err := kind.ApplyPath(&cliResource)
	if err != nil {
		return "", err
	}

	url := client.BaseUrl + applyPath.Path

	tflog.Trace(ctx, fmt.Sprintf("PUT on %s body : %s", applyPath, string(cliResource.Json)))
	builder := client.Client.R().SetBody(cliResource.Json)
	resp, err := builder.Put(url)
	if err != nil {
		return "", err
	} else if resp.IsError() {
		return "", fmt.Errorf("%s", ExtractApiError(resp))
	}
	bodyBytes := resp.Body()
	var upsertResponse ApplyResult
	err = jsoniter.Unmarshal(bodyBytes, &upsertResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %s", err)
	}
	return upsertResponse.UpsertResult, nil
}

func (client *Client) Apply(ctx context.Context, path string, resource any) (ApplyResult, error) {
	url := client.BaseUrl + path
	jsonData, err := jsoniter.Marshal(resource)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("error marshalling resource: %s", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("PUT %s request body : %s", path, string(jsonData)))

	resp, err := client.Client.R().SetBody(jsonData).Put(url)
	if err != nil {
		return ApplyResult{}, err
	} else if resp.IsError() {
		return ApplyResult{}, fmt.Errorf("%s", ExtractApiError(resp))
	}

	bodyBytes := resp.Body()
	tflog.Trace(ctx, fmt.Sprintf("PUT %s response body : %s", path, string(bodyBytes)))

	var upsertResponse ApplyResult
	err = jsoniter.Unmarshal(bodyBytes, &upsertResponse)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("error unmarshalling response: %s", err)
	}
	return upsertResponse, nil
}

func (client *Client) Describe(ctx context.Context, path string) ([]byte, error) {
	url := client.BaseUrl + path
	resp, err := client.Client.R().Get(url)
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

func (client *Client) Delete(ctx context.Context, mode Mode, path string, resource any) error {
	var req *resty.Request
	url := client.BaseUrl + path
	tflog.Trace(ctx, fmt.Sprintf("DELETE %s", path))

	switch mode {
	case CONSOLE:
		req = client.Client.R()
	case GATEWAY:
		// Gateway API handles deletion in a different way
		// It needs information about the resource in the body of the request
		// as opposed to Console API who needs them in the URL
		jsonData, err := json.Marshal(resource)
		if err != nil {
			return fmt.Errorf("error marshalling resource: %s", err)
		}
		tflog.Debug(ctx, string(jsonData))
		tflog.Trace(ctx, fmt.Sprintf("DELETE %s request body : %s", path, string(jsonData)))

		req = client.Client.R().SetBody(string(jsonData))
	}

	resp, err := req.Delete(url)
	if err != nil {
		return err
	} else if resp.IsError() {
		return fmt.Errorf("%s", ExtractApiError(resp))
	}

	return nil
}

func (client *Client) GetAPIVersion(ctx context.Context, mode Mode) (string, error) {
	var path string
	if mode == CONSOLE {
		path = "/versions"
	}
	if mode == GATEWAY {
		path = "/health"
	}

	url := client.BaseUrl + path
	resp, err := client.Client.R().Get(url)
	if err != nil {
		return "", err
	} else if resp.IsError() {
		return "", fmt.Errorf("%s", ExtractApiError(resp))
	}
	tflog.Trace(ctx, fmt.Sprintf("GET %s response : %s", path, string(resp.Body())))

	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	var v any
	var ok bool
	if mode == CONSOLE {
		v, ok = result["platform"]
		if !ok {
			return "", fmt.Errorf("no version found in response")
		}
	}
	if mode == GATEWAY {
		// This is a temporary workaround for the Gateway API till a dedicated endpoint is available.
		checks, ok := result["checks"].([]any)
		if !ok || len(checks) == 0 {
			return "", fmt.Errorf("no checks found in response")
		}
		for _, check := range checks {
			// Need to assert check to map[string]any to access its fields.
			id, ok := check.(map[string]any)["id"]
			if !ok {
				return "", fmt.Errorf("error parsing check ID")
			}
			if id == "buildInfo" {
				data, ok := check.(map[string]any)["data"]
				if !ok {
					return "", fmt.Errorf("no data found in checks response")
				}
				v, ok = data.(map[string]any)["version"]
				if !ok {
					return "", fmt.Errorf("no version found in data response")
				}
			}
		}
	}

	version, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("error parsing version to string")
	}

	// Go module semver requires version to start with v.
	return "v" + version, nil
}
