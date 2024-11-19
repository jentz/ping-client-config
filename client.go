package pcc

import (
	"github.com/jentz/ping-client-config/internal/adminapi"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	ClientID                    string   `yaml:"client_id"`
	Name                        string   `yaml:"name"`
	Description                 *string  `yaml:"description,omitempty"`
	RedirectURIs                []string `yaml:"redirect_uris,omitempty"`
	GrantTypes                  []string `yaml:"grant_types,omitempty"`
	OIDCPolicyID                string   `yaml:"oidc_policy_id,omitempty"`
	DefaultAccessTokenManagerID string   `yaml:"default_access_token_manager_id,omitempty"`
	RequirePKCE                 bool     `yaml:"require_pkce,omitempty"`
	RestrictedScopes            []string `yaml:"restricted_scopes,omitempty"`
	ExclusiveScopes             []string `yaml:"exclusive_scopes,omitempty"`
	BypassApprovalPage          bool     `yaml:"bypass_approval_page,omitempty"`
	SSOEnabled                  bool     `yaml:"sso_enabled,omitempty"`
}

func marshalClientConfig(c *ClientConfig) ([]byte, error) {
	return yaml.Marshal(c)
}

func marshalClientConfigs(c []ClientConfig) ([]byte, error) {
	return yaml.Marshal(c)
}

func createClientConfig(pingClient client.Client) ClientConfig {
	return ClientConfig{
		ClientID:                    pingClient.ClientId,
		Name:                        pingClient.Name,
		Description:                 pingClient.Description,
		RedirectURIs:                pingClient.RedirectUris,
		GrantTypes:                  pingClient.GrantTypes,
		OIDCPolicyID:                pingClient.OidcPolicy.PolicyGroup.GetId(),
		DefaultAccessTokenManagerID: pingClient.DefaultAccessTokenManagerRef.GetId(),
		RequirePKCE:                 *pingClient.RequireProofKeyForCodeExchange,
		RestrictedScopes:            pingClient.RestrictedScopes,
		ExclusiveScopes:             pingClient.ExclusiveScopes,
		BypassApprovalPage:          *pingClient.BypassApprovalPage,
		SSOEnabled:                  adminapi.GetSingleExtendedParameterValue(pingClient.ExtendedParameters, "adapter_type", "") == "sso",
	}
}

func createClientConfigs(pingClients []client.Client) []ClientConfig {
	var clients []ClientConfig
	for _, c := range pingClients {
		clients = append(clients, createClientConfig(c))
	}
	return clients
}

func (c *ClientConfig) Marshal() ([]byte, error) {
	return marshalClientConfig(c)
}

func unmarshalClientConfig(data []byte) (*ClientConfig, error) {
	var c ClientConfig
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
