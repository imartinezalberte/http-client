package httpclient

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultTimeout = 10 * time.Second
	minTimeout     = 5 * time.Second
	maxTimeout     = 60 * time.Second

	defaultRetryCount = 3
	minRetryCount     = 0
	maxRetryCount     = 5

	defaultRetryWaitTime = 1 * time.Second
	minRetryWaitTime     = 100 * time.Millisecond
	maxRetryWaitTime     = 5 * time.Second

	defaultRetryMaxWaitTime = 2 * time.Second
	minRetryMaxWaitTime     = 1 * time.Second
	maxRetryMaxWaitTime     = 10 * time.Second

	defaultLogLevel = zapcore.InfoLevel
)

var (
	checkTimeout          = CheckValueInRange(minTimeout, maxTimeout, defaultTimeout)
	checkRetryCount       = CheckValueInRange(minRetryCount, maxRetryCount, defaultRetryCount)
	checkRetryWaitTime    = CheckValueInRange(minRetryWaitTime, maxRetryWaitTime, defaultRetryWaitTime)
	checkRetryMaxWaitTime = CheckValueInRange(minRetryMaxWaitTime, maxRetryMaxWaitTime, defaultRetryMaxWaitTime)
)

var (
	ErrHostNeeded        = errors.New("host attribute is needed")
	ErrIntegrationNeeded = errors.New("integration attribute is needed")
)

func (c *Config) UnmarshalJSON(b []byte) error {
	type Tmp Config
	if err := json.Unmarshal(b, (*Tmp)(c)); err != nil {
		return err
	}

	c.SetDefaultTimeout().
		SetDefaultRetryCount().
		SetDefaultRetryWaitTime().
		SetDefaultRetryMaxWaitTime()

	return c.Check()
}

func (c *Config) Check() error {
	if strings.TrimSpace(c.Host) == "" {
		return ErrHostNeeded
	}

	if strings.TrimSpace(c.Integration) == "" {
		return ErrIntegrationNeeded
	}

	return nil
}

// SetHost method allows you to modify or set the host value programmatically instead of using a configuration file.
//
//		config.SetHost("http://www.example.com")
func (c *Config) SetHost(host string) *Config {
	c.Host = host
	return c
}

// SetTimeout method allows you to modify or set the timeout for the request.
//
//		config.SetTimeout(2 * time.Second)
func (c *Config) SetTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// SetDefaultTimeout method allows you to set the default timeout for the request.
// If the actual value of timeout is less than `minTimeout` or greater than `maxTimeout`,
// then the default is set.
// See `defaultTimeout`
//
//		config.SetDefaultTimeout()
func (c *Config) SetDefaultTimeout() *Config {
	c.Timeout = checkTimeout(c.Timeout)
	return c
}

// SetRetryCount method allows you to set the amount of retries allowed per request.
//
//		config.SetRetryCount(3)
func (c *Config) SetRetryCount(retryCount int) *Config {
	c.Retry.Count = retryCount
	return c
}

// SetDefaultRetryCount method allows you to set the default retry count for each request.
// If the actual value of retryCount if less than `minRetryCount` or greater than `maxRetryCount`,
// then the default is set.
// See `defaultRetryCount`
//
//		config.SetDefaultRetryCount()
func (c *Config) SetDefaultRetryCount() *Config {
	c.Retry.Count = checkRetryCount(c.Retry.Count)
	return c
}

// SetRetryWaitTime method allows you to set the wait time between each retry in each request.
//
//		config.SetRetryWaitTime(1 * time.Second)
func (c *Config) SetRetryWaitTime(retryWaitTime time.Duration) *Config {
	c.Retry.WaitTime = retryWaitTime
	return c
}

// SetDefaultRetryWaitTime method allows you to set the default retry wait time between each request.
// If the actual value of retryWaitTime is less than `minRetryWaitTime` or greater than `maxRetryWaitTime`,
// then the default is set.
// See `defaultRetryWaitTime`
//
//		config.SetDefaultRetryWaitTime()
func (c *Config) SetDefaultRetryWaitTime() *Config {
	c.Retry.WaitTime = checkRetryWaitTime(c.Retry.WaitTime)
	return c
}

// SetRetryMaxWaitTime method allows you to set the max wait time between each request when retrying.
//
//		config.SetRetryMaxWaitTime(2 * time.Second)
func (c *Config) SetRetryMaxWaitTime(retryMaxWaitTime time.Duration) *Config {
	c.Retry.MaxWaitTime = retryMaxWaitTime
	return c
}

// SetDefaultRetryMaxWaitTime method allows you to set the default retry max wait time between each request.
// If the actual value of retryMaxWaitTime is less than `minRetryMaxWaitTime` or greater than `maxRetryMaxWaitTime`,
// then the default is set.
// See `defaultRetryMaxWaitTime`
//
//		config.SetDefaultRetryMaxWaitTime()
func (c *Config) SetDefaultRetryMaxWaitTime() *Config {
	c.Retry.MaxWaitTime = checkRetryMaxWaitTime(c.Retry.MaxWaitTime)
	return c
}

// SetLogLevel method allows us to set the log level that we want to our client.
//
//		config.SetLogLevel(zapcore.DebugLevel)
func (c *Config) SetLogLevel(logLevel zapcore.Level) *Config {
	c.LogLevel = logLevel
	return c
}

func (c *Config) SetOfuscateQueryParam(queryParam string) *Config {
	c.Ofuscate.QueryParams = append(c.Ofuscate.QueryParams, queryParam)
	return c
}

func (c *Config) SetOfuscateQueryParams(queryParams ...string) *Config {
	c.Ofuscate.QueryParams = append(c.Ofuscate.QueryParams, queryParams...)
	return c
}

func (c Config) OfuscateQueryParams() OfuscateParamsFunc {
	queryParams := ArrToMap(c.Ofuscate.QueryParams)
	return ofuscate(queryParams)
}

func (c *Config) SetOfuscateHeader(header string) *Config {
	c.Ofuscate.Headers = append(c.Ofuscate.Headers, header)
	return c
}

func (c *Config) SetOfuscateHeaders(headers ...string) *Config {
	c.Ofuscate.Headers = append(c.Ofuscate.Headers, headers...)
	return c
}

func (c Config) OfuscateHeaders() OfuscateParamsFunc {
	headers := ArrToMap(c.Ofuscate.Headers)
	return ofuscate(headers)
}

func (c *Config) SetOfuscateRequest(request string) *Config {
	c.Ofuscate.Request = append(c.Ofuscate.Request, request)
	return c
}

func (c *Config) SetOfuscateRequests(requests ...string) *Config {
	c.Ofuscate.Request = append(c.Ofuscate.Request, requests...)
	return c
}

func (c Config) OfuscateRequests() OfuscateBodyFunc {
	return ofuscateBody(c.Ofuscate.Request)
}

func (c *Config) SetOfuscateResponse(response string) *Config {
	c.Ofuscate.Response = append(c.Ofuscate.Response, response)
	return c
}

func (c *Config) SetOfuscateResponses(responses ...string) *Config {
	c.Ofuscate.Response = append(c.Ofuscate.Response, responses...)
	return c
}

func (c Config) OfuscateResponses() OfuscateBodyFunc {
	return ofuscateBody(c.Ofuscate.Response)
}

func NewClient(cl *http.Client, logger ILogger, host, integration string) *resty.Client {
	c := resty.NewWithClient(cl).
		SetHostURL(host).
		SetTimeout(defaultTimeout).
		SetRetryCount(defaultRetryCount).
		SetRetryWaitTime(defaultRetryWaitTime).
		SetRetryMaxWaitTime(defaultRetryMaxWaitTime)

	c.OnBeforeRequest(OnBeforeRequest(logger, Config{Integration: integration})).
		OnAfterResponse(OnAfterResponse(logger, Config{Integration: integration}))

	return c
}

func NewClientFromConfig(cl *http.Client, logger ILogger, config Config) *resty.Client {
	c := resty.NewWithClient(cl).
		SetHostURL(config.Host).
		SetTimeout(config.Timeout).
		SetRetryCount(config.Retry.Count).
		SetRetryWaitTime(config.Retry.WaitTime).
		SetRetryMaxWaitTime(config.Retry.MaxWaitTime)

	if config.LogLevel.Enabled(zapcore.InfoLevel) {
		c = c.OnBeforeRequest(OnBeforeRequest(logger, config)).
			OnAfterResponse(OnAfterResponse(logger, config))
	}

	return c
}

func OnBeforeRequest(logger ILogger, config Config) func(*resty.Client, *resty.Request) error {
	var (
		ofuscateHeaders     = config.OfuscateHeaders()
		ofuscateQueryParams = config.OfuscateQueryParams()
		ofuscateBody        = config.OfuscateRequests()
	)

	return func(c *resty.Client, r *resty.Request) error {
		logger.Info(r.Context(), "",
			zap.String("integration", config.Integration),
			zap.String("method", r.Method),
			zap.String("path", r.URL),
			zap.Any("headers", ofuscateHeaders(r.Header)),
			zap.Any("query", ofuscateQueryParams(r.QueryParam)),
			zap.Any("body", ofuscateBody(r.Body)),
		)

		return nil
	}
}

func OnAfterResponse(logger ILogger, config Config) func(*resty.Client, *resty.Response) error {
	var (
		ofuscateHeaders     = config.OfuscateHeaders()
		ofuscateQueryParams = config.OfuscateQueryParams()
		ofuscateBody        = config.OfuscateResponses()
	)

	return func(c *resty.Client, r *resty.Response) error {
		var m map[string]interface{}
		if err := json.Unmarshal(r.Body(), &m); err != nil {
			logger.Info(r.Request.Context(), "couldn't parse body of the response", zap.Error(err))
		}

		logger.Info(r.Request.Context(), "",
			zap.String("integration", config.Integration),
			zap.Int("statusCode", r.StatusCode()),
			zap.Duration("time", r.Time()),
			zap.String("method", r.Request.Method),
			zap.String("path", r.Request.RawRequest.URL.Path),
			zap.Any("headers", ofuscateHeaders(r.Header())),
			zap.Any("query", ofuscateQueryParams(r.Request.QueryParam)),
			zap.Any("body", ofuscateBody(m)))

		return nil
	}
}
