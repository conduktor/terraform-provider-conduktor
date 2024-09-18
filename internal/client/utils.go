package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ApiError struct {
	Title string
	Msg   string
}

func (e *ApiError) String() string {
	if e.Msg == "" {
		return e.Title
	} else {
		return e.Msg
	}
}

func extractApiError(resp *resty.Response) string {
	var apiError ApiError
	jsonError := json.Unmarshal(resp.Body(), &apiError)
	if jsonError != nil {
		return resp.String()
	} else {
		return apiError.String()
	}
}

func uniformizeBaseUrl(baseUrl string) string {
	regex := regexp.MustCompile(`(/api)?/?$`)
	return regex.ReplaceAllString(baseUrl, "/api")
}

func retry(attempts int, delay time.Duration, ctx context.Context) func(filter func(error) bool, f func() (interface{}, error)) (interface{}, error) {
	return func(filter func(error) bool, f func() (interface{}, error)) (interface{}, error) {
		var result interface{}
		var err error
		for i := 0; i < attempts; i++ {
			result, err = f()
			if err != nil {
				if filter != nil && !filter(err) {
					return nil, err
				}
				if i < attempts-1 {
					tflog.Warn(ctx, fmt.Sprintf("Retrying after error: %v", err))
					time.Sleep(delay)
					continue
				} else {
					return nil, err
				}
			}
			break
		}
		return result, nil
	}
}

// hack to guess if current provider logs are in trace level.
func TraceLogEnabled() bool {
	selfLevel := strings.ToUpper(os.Getenv("TF_LOG_PROVIDER_CONDUKTOR"))
	providersLevel := strings.ToUpper(os.Getenv("TF_LOG_PROVIDER"))
	terraformLevel := strings.ToUpper(os.Getenv("TF_LOG"))

	return terraformLevel == "TRACE" || providersLevel == "TRACE" || selfLevel == "TRACE"
}
