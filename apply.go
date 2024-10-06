package pcc

import (
	"context"
	"fmt"
	"github.com/jentz/ping-client-config/internal/adminapi"
	pfclient "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"io"
	"log"
	"net/http"
	"os"
)

type ApplyCommand struct {
	AdminURL string
	Username string
	Password string

	ClientConfigFileName string
}

func (c *ApplyCommand) Run(ctx context.Context) error {
	cfg := adminapi.NewConfig().WithEndpointURL(c.AdminURL).
		WithUsername(c.Username).WithPassword(c.Password)

	adminClient := adminapi.NewAdminClient(cfg)
	ctx = adminapi.AuthContext(ctx, *cfg)

	log.Printf("ping admin host %s\n", c.AdminURL)

	// open the client config file and unmarshal the client
	clientConfigFile, err := os.Open(c.ClientConfigFileName)
	if err != nil {
		return fmt.Errorf("error opening client config file: %v", err)
	}
	defer func(clientConfigFile *os.File) {
		err := clientConfigFile.Close()
		if err != nil {
			log.Printf("error closing client config file: %v", err)
		}
	}(clientConfigFile)

	// read the file contents
	data, err := io.ReadAll(clientConfigFile)
	if err != nil {
		return fmt.Errorf("error reading client config file: %v", err)
	}

	clientIn, err := unmarshalClientConfig(data)
	if err != nil {
		return fmt.Errorf("error unmarshalling client: %v", err)
	}

	client, r, err := adminClient.OauthClientsAPI.GetOauthClientById(ctx, clientIn.ClientID).Execute()
	if err != nil {
		if r.StatusCode != http.StatusNotFound {
			return fmt.Errorf("error getting client %s: %v", clientIn.ClientID, err)
		}
		// create the client
		log.Printf("creating client %s\n", clientIn.ClientID)
		client = clientFromClientConfig(clientIn)
		_, r, err = adminClient.OauthClientsAPI.CreateOauthClient(ctx).Body(*client).Execute()
		if err != nil {
			return fmt.Errorf("error creating client %s: %v", clientIn.ClientID, err)
		}

		return nil
	}

	// update the client
	log.Printf("updating client %s\n", client.ClientId)

	pccManagedString := adminapi.GetSingleExtendedParameterValue(client.ExtendedParameters, "pcc_managed", "false")

	if pccManagedString != "true" {
		return fmt.Errorf("client %s is not managed by pcc", client.ClientId)
	}
	clientFromConfig := clientFromClientConfig(clientIn)
	clientFromConfig.ClientAuth = client.ClientAuth
	_, r, err = adminClient.OauthClientsAPI.UpdateOauthClient(ctx, client.ClientId).Body(*clientFromConfig).Execute()
	if err != nil {
		return fmt.Errorf("error updating client %s: %v", client.ClientId, err)
	}

	return nil
}

func clientFromClientConfig(clientConfig *ClientConfig) *pfclient.Client {
	authType := "NONE"
	oauthClient := &pfclient.Client{
		ClientId:     clientConfig.ClientID,
		Name:         clientConfig.Name,
		Description:  clientConfig.Description,
		RedirectUris: clientConfig.RedirectURIs,
		GrantTypes:   clientConfig.GrantTypes,
		ClientAuth: &pfclient.ClientAuth{
			Type: &authType,
		},
		RequireProofKeyForCodeExchange: boolPtr(clientConfig.RequirePKCE),
		BypassApprovalPage:             boolPtr(clientConfig.BypassApprovalPage),
		OidcPolicy: &pfclient.ClientOIDCPolicy{
			PolicyGroup: &pfclient.ResourceLink{
				Id: clientConfig.OIDCPolicyID,
			},
		},
		DefaultAccessTokenManagerRef: &pfclient.ResourceLink{
			Id: clientConfig.DefaultAccessTokenManagerID,
		},
	}

	// set the extended parameters
	parameters := make(map[string]pfclient.ParameterValues)
	oauthClient.ExtendedParameters = &parameters

	parameters["pcc_managed"] = pfclient.ParameterValues{
		Values: []string{"true"},
	}

	// set the adapter type
	var adapterType string
	if clientConfig.SSOEnabled {
		adapterType = "sso"
	} else {
		adapterType = "no_sso"
	}

	parameters["adapter_type"] = pfclient.ParameterValues{
		Values: []string{adapterType},
	}
	return oauthClient
}
