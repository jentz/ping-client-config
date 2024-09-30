package pcc

import (
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"gopkg.in/yaml.v3"
)

type Client struct {
	ClientID     string   `yaml:"client_id"`
	Name         string   `yaml:"name"`
	Description  *string  `yaml:"description,omitempty"`
	RedirectURIs []string `yaml:"redirect_uris,omitempty"`
	GrantTypes   []string `yaml:"grant_types,omitempty"`
}

func marshalClient(c *Client) ([]byte, error) {
	return yaml.Marshal(c)
}

func marshalClients(c []Client) ([]byte, error) {
	return yaml.Marshal(c)
}

func createClient(pingClient client.Client) Client {
	return Client{
		ClientID:     pingClient.ClientId,
		Name:         pingClient.Name,
		Description:  pingClient.Description,
		RedirectURIs: pingClient.RedirectUris,
		GrantTypes:   pingClient.GrantTypes,
	}
}

func createClients(pingClients []client.Client) []Client {
	var clients []Client
	for _, c := range pingClients {
		clients = append(clients, createClient(c))
	}
	return clients
}

func (c *Client) Marshal() ([]byte, error) {
	return marshalClient(c)
}
