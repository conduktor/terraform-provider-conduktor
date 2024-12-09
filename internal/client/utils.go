package client

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
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

// hack to guess if current provider logs are in trace level.
func TraceLogEnabled() bool {
	selfLevel := strings.ToUpper(os.Getenv("TF_LOG_PROVIDER_CONDUKTOR"))
	providersLevel := strings.ToUpper(os.Getenv("TF_LOG_PROVIDER"))
	terraformLevel := strings.ToUpper(os.Getenv("TF_LOG"))

	return terraformLevel == "TRACE" || providersLevel == "TRACE" || selfLevel == "TRACE"
}
