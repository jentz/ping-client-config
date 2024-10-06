package adminapi

import (
	"context"
	"fmt"
	pfclient "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"net/http"
)

type Config struct {
	HttpClient  *http.Client
	EndpointURL string
	Username    string
	Password    string
}

func AuthContext(ctx context.Context, config Config) context.Context {
	// add other auth methods as required
	return BasicAuthContext(ctx, config.Username, config.Password)
}

// BasicAuthContext Get BasicAuth context with a username and password
func BasicAuthContext(ctx context.Context, username, password string) context.Context {
	return context.WithValue(ctx, pfclient.ContextBasicAuth, pfclient.BasicAuth{
		UserName: username,
		Password: password,
	})
}

// NewConfig creates a new Config
func NewConfig() *Config {
	return &Config{
		HttpClient: http.DefaultClient,
	}
}

// WithUsername sets the username
func (c *Config) WithUsername(username string) *Config {
	c.Username = username
	return c
}

// WithPassword sets the password
func (c *Config) WithPassword(password string) *Config {
	c.Password = password
	return c
}

// WithEndpointURL sets the endpoint URL
func (c *Config) WithEndpointURL(endpointURL string) *Config {
	c.EndpointURL = endpointURL
	return c
}

// WithHTTPClient sets the HTTP client
func (c *Config) WithHTTPClient(httpClient *http.Client) *Config {
	c.HttpClient = httpClient
	return c
}

func newConfiguration(config *Config) *pfclient.Configuration {
	clientConfig := pfclient.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.Servers = pfclient.ServerConfigurations{
		{
			URL: config.EndpointURL,
		},
	}
	clientConfig.HTTPClient = config.HttpClient
	userAgentSuffix := fmt.Sprintf("ping-client-config/%s %s", "v1210", "go")
	clientConfig.UserAgentSuffix = &userAgentSuffix
	return clientConfig
}

func NewAdminClient(config *Config) *pfclient.APIClient {
	return pfclient.NewAPIClient(newConfiguration(config))
}
