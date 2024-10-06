package pcc

import (
	"context"
	"fmt"
	"github.com/jentz/ping-client-config/internal/adminapi"
	"log"
)

type GetCommand struct {
	AdminURL string
	Username string
	Password string

	All      bool
	ClientID string
}

func (c *GetCommand) Run(ctx context.Context) error {
	cfg := adminapi.NewConfig().WithEndpointURL(c.AdminURL).
		WithUsername(c.Username).WithPassword(c.Password)

	adminClient := adminapi.NewAdminClient(cfg)
	ctx = adminapi.AuthContext(ctx, *cfg)

	log.Printf("ping admin host %s\n", c.AdminURL)

	if c.All {
		clients, r, err := adminClient.OauthClientsAPI.GetOauthClients(ctx).Execute()
		if err != nil {
			return fmt.Errorf("error getting clients: %v", err)
		}

		var rps = createClientConfigs(clients.Items)
		clientYaml, err := marshalClientConfigs(rps)
		if err != nil {
			return fmt.Errorf("error marshalling rp: %v", err)
		}
		log.Printf("\n%s\n", clientYaml)
		log.Println(r.Status)
	} else {
		// get a single client
		pingClient, r, err := adminClient.OauthClientsAPI.GetOauthClientById(ctx, c.ClientID).Execute()
		if err != nil {
			return fmt.Errorf("error getting client %s: %v", c.ClientID, err)
		}
		rp := createClientConfig(*pingClient)
		clientYaml, err := rp.Marshal()
		if err != nil {
			return fmt.Errorf("error marshalling rp: %v", err)
		}
		log.Printf("\n%s\n", clientYaml)
		log.Println(r.Status)
	}

	return nil
}
