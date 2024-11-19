package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
)

type CommandRunner interface {
	Run(ctx context.Context) error
}

type Command struct {
	Name      string
	Help      string
	Configure func(name string, args []string) (config CommandRunner, output string, err error)
}

var commands = []Command{
	{Name: "get", Help: "Gets the client configs, either all of them or one", Configure: parseGetFlags},
	{Name: "apply", Help: "Applies a client config", Configure: parseApplyFlags},
	{Name: "help", Help: "Prints help"},
}

func usage() {
	intro := `pcc is a command-line client that manages ping client configs

Usage:
  pcc [flags] <command> [command-flags]`
	fmt.Fprintln(os.Stderr, intro)
	fmt.Fprintln(os.Stderr, "\nCommands:")
	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "  %-18s: %s\n", cmd.Name, cmd.Help)
	}
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Run `pcc <command> -h` to get help for a specific command\n\n")
}

func runCommand(name string, args []string) {

	cmdIdx := slices.IndexFunc(commands, func(cmd Command) bool {
		return cmd.Name == name
	})

	if cmdIdx < 0 {
		fmt.Fprintf(os.Stderr, "command \"%s\" not found\n\n", name)
		flag.Usage()
		os.Exit(1)
	}

	cmd := commands[cmdIdx]
	if cmd.Name == "help" {
		flag.Usage()
		os.Exit(0)
	}

	command, output, err := cmd.Configure(name, args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println("got error:", err)
		fmt.Println("output:\n", output)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := command.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}

	subCmd := flag.Arg(0)
	subArgs := flag.Args()[1:]
	runCommand(subCmd, subArgs)
}
