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

func ExtractApiError(resp *resty.Response) string {
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

func isLogLevelEnabled(envVars []string, level string) bool {
	for _, envVar := range envVars {
		if strings.ToUpper(os.Getenv(envVar)) == level {
			return true
		}
	}
	return false
}

// hack to guess if current provider logs are in trace level.
func TraceLogEnabled() bool {
	return isLogLevelEnabled([]string{"TF_LOG_PROVIDER_CONDUKTOR", "TF_LOG_PROVIDER", "TF_LOG"}, "TRACE")
}

// hack to guess if current provider logs are in debug or trace level.
func DebugLogEnabled() bool {
	return TraceLogEnabled() || isLogLevelEnabled([]string{"TF_LOG_PROVIDER_CONDUKTOR", "TF_LOG_PROVIDER", "TF_LOG"}, "DEBUG")
}

func InitTraceEnabled() bool {
	initLevel := strings.ToUpper(os.Getenv("TF_LOG_PROVIDER_CONDUKTOR_INIT"))
	return initLevel == "TRACE"
}
