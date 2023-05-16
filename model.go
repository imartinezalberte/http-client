package httpclient

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type (
	// Config struct is the main entry point to configure your client. Here
	// we can configure the host, the name of the integration, timeout when fething resources,
	// retry stuff, log level, etc.
	Config struct {
		Integration string        `json:"integration"         yaml:"integration"`
		Host        string        `json:"host"                yaml:"host"`
		Timeout     time.Duration `json:"timeout,omitempty"   yaml:"timeout"`
		Retry       Retry         `json:"retry,omitempty"     yaml:"retry"`
		LogLevel    zapcore.Level `json:"log_level,omitempty" yaml:"log_level"`
		Ofuscate    Ofuscate      `json:"ofuscate,omitempty"  yaml:"ofuscate"`
	}

	// Retry struct allows us to specify how many times we want to retry to execute
	// the request, how many time are we able to wait between retries and the maximum
	// allowed.
	Retry struct {
		Count       int           `json:"count,omitempty"         yaml:"count"`
		WaitTime    time.Duration `json:"wait_time,omitempty"     yaml:"wait_time"`
		MaxWaitTime time.Duration `json:"max_wait_time,omitempty" yaml:"max_wait_time"`
	}

	// Ofuscate struct is a really interesting one, because it allows us to ofuscate
	// all the parameters that we desire:
	//	- Query
	//	- Headers
	//	- Request body
	//	- Response body
	Ofuscate struct {
		QueryParams []string `json:"query_params,omitempty" yaml:"query_params"`
		Headers     []string `json:"headers,omitempty"      yaml:"headers"`
		Request     []string `json:"request,omitempty"      yaml:"request"`
		Response    []string `json:"response,omitempty"     yaml:"response"`
	}
)
