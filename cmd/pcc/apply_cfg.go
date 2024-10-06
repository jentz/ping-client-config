package main

import (
	"bytes"
	"errors"
	"flag"
	pcc "github.com/jentz/ping-client-config"
	"os"
)

func parseApplyFlags(name string, args []string) (config CommandRunner, output string, err error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	runner := &pcc.ApplyCommand{}

	flags.StringVar(&runner.AdminURL, "admin-url", "", "set admin url (required)")
	flags.StringVar(&runner.Username, "username", "", "set username (required)")
	flags.StringVar(&runner.Password, "password", "", "set password (required)")
	flags.StringVar(&runner.ClientConfigFileName, "client-config-file", "", "set client config file (required)")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	// if the environment variable is set, use it as the default
	if runner.Username == "" {
		if username, ok := os.LookupEnv("PINGFEDERATE_USERNAME"); ok {
			runner.Username = username
		}
	}

	if runner.Password == "" {
		if password, ok := os.LookupEnv("PINGFEDERATE_PASSWORD"); ok {
			runner.Password = password
		}
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
		{
			runner.ClientConfigFileName == "",
			"client-config-file is required",
		},
	}
	for _, check := range invalidArgsChecks {
		if check.condition {
			return nil, buf.String(), errors.New(check.message)
		}
	}
	return runner, buf.String(), nil
}
