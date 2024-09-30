package pcc

import (
	"context"
	"fmt"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"log"
	"net/http"
	"os"
)

type GetCommand struct {
	AdminURL string
	Username string
	Password string

	All      bool
	ClientID string
}

func (c *GetCommand) Run() error {
	ctx := context.Background()
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	//clientConfig.DefaultHeader["X-BypassExternalValidation"] = strconv.FormatBool(xBypassExternalValidation)
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: c.AdminURL,
		},
	}
	httpClient := http.DefaultClient
	clientConfig.HTTPClient = httpClient
	userAgentSuffix := fmt.Sprintf("pfclientconf/%s %s", "v1210", "go")
	clientConfig.UserAgentSuffix = &userAgentSuffix
	apiClient := client.NewAPIClient(clientConfig)

	apiCtx := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: os.Getenv("PINGFEDERATE_USERNAME"),
		Password: os.Getenv("PINGFEDERATE_PASSWORD"),
	})

	log.Printf("ping admin host %s\n", c.AdminURL)

	if c.All {
		clients, r, err := apiClient.OauthClientsAPI.GetOauthClients(apiCtx).Execute()
		if err != nil {
			return fmt.Errorf("error getting clients: %v", err)
		}

		var rps = createClients(clients.Items)
		clientYaml, err := marshalClients(rps)
		if err != nil {
			return fmt.Errorf("error marshalling rp: %v", err)
		}
		log.Printf("\n%s\n", clientYaml)
		log.Println(r.Status)
	} else {
		// get a single client
		pingClient, r, err := apiClient.OauthClientsAPI.GetOauthClientById(apiCtx, c.ClientID).Execute()
		if err != nil {
			return fmt.Errorf("error getting client %s: %v", c.ClientID, err)
		}
		rp := createClient(*pingClient)
		clientYaml, err := rp.Marshal()
		if err != nil {
			return fmt.Errorf("error marshalling rp: %v", err)
		}
		log.Printf("\n%s\n", clientYaml)
		log.Println(r.Status)
	}

	return nil
}
