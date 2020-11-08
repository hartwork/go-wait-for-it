/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/spf13/cobra"
)

type Config struct {
	Addresses []syntax.Address
	Argv      []string
	Quiet     bool
	Timeout   time.Duration
}

func report(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
}

func Parse(args []string) (config *Config, err error) {
	var services []string
	var quiet bool
	var timeoutSeconds uint

	rootCommand := &cobra.Command{
		Use:  "wait-for-it [flags] [-s|--service [HOST]:PORT]... [--] [COMMAND [ARG ..]]",
		Long: "Wait for service(s) to be available before executing a command.",
		Run: func(cmd *cobra.Command, args []string) {
			timeout := time.Duration(timeoutSeconds) * time.Second

			var addresses []syntax.Address
			for _, service := range services {
				address, syntaxError := syntax.ParseAddress(service)
				if syntaxError != nil {
					report(syntaxError)
					err = syntaxError // the first error is as good as the last, here
					continue
				}
				addresses = append(addresses, address)
			}

			config = &Config{
				addresses,
				args,
				quiet,
				timeout,
			}
		},
		Version: "1.0.0",
	}

	rootCommand.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	rootCommand.Flags().StringSliceVarP(&services, "service", "s",
		[]string{},
		"services to test (format '[HOST]:PORT')")
	rootCommand.Flags().UintVarP(&timeoutSeconds, "timeout", "t",
		15, "timeout in seconds, 0 for no timeout")
	rootCommand.Flags().BoolVarP(&quiet, "quiet", "q",
		false, "do not output any status messages")

	rootCommand.SetArgs(args)

	if executeError := rootCommand.Execute(); executeError != nil {
		report(executeError)
		err = executeError
	}

	return config, err
}
