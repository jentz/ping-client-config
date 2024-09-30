package main

import (
	"bytes"
	"flag"
	"fmt"
	pcc "github.com/jentz/ping-client-config"
	"os"
)

func parseGetFlags(name string, args []string) (config CommandRunner, output string, err error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	runner := &pcc.GetCommand{}
	// if the environment variable is set, use it as the default
	defaultUsername := ""
	if username, ok := os.LookupEnv("PINGFEDERATE_USERNAME"); ok {
		defaultUsername = username
	}

	defaultPassword := ""
	if password, ok := os.LookupEnv("PINGFEDERATE_PASSWORD"); ok {
		defaultPassword = password
	}

	flags.StringVar(&runner.AdminURL, "admin-url", "", "set admin url (required)")
	flags.StringVar(&runner.Username, "username", defaultUsername, "set username (required)")
	flags.StringVar(&runner.Password, "password", defaultPassword, "set password (required)")
	flags.StringVar(&runner.ClientID, "client-id", "", "set client ID")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	if runner.ClientID == "" {
		runner.All = true
	}

	var invalidArgsChecks = []struct {
		condition bool
		message   string
	}{
		{
			runner.AdminURL == "",
			"admin-url is required",
		},
		{
			runner.Username == "",
			"username is required",
		},
		{
			runner.Password == "",
			"password is required",
		},
	}

	for _, check := range invalidArgsChecks {
		if check.condition {
			return nil, "", fmt.Errorf(check.message)
		}
	}
	return runner, buf.String(), nil
}
