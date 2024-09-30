package pcc

import (
	"context"
	"fmt"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"log"
	"net/http"
	"os"
	"strconv"
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
			log.Fatal(err)
		}

		log.Printf("Found %d clients\n", len(clients.Items))
		// loop through the clients
		for _, c := range clients.Items {
			log.Printf("%s - %s - %s\n", c.ClientId, c.Name, strconv.FormatBool(*c.Enabled))
		}

		log.Println(r.Status)
		return nil
	}

	// get a single client
	client, r, err := apiClient.OauthClientsAPI.GetOauthClientById(apiCtx, c.ClientID).Execute()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s - %s - %s\n", client.ClientId, client.Name, strconv.FormatBool(*client.Enabled))
	log.Println(r.Status)
	return nil
}
