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
	fmt.Println("Error:", err)
}

func Parse(args []string) (config Config) {
	var services []string
	var timeoutSeconds uint

	rootCommand := &cobra.Command{
		Use:  "wait-for-it [flags] [-s|--service [HOST]:PORT]... [--] [COMMAND [ARG ..]]",
		Long: "Wait for service(s) to be available before executing a command.",
		Run: func(cmd *cobra.Command, args []string) {
			config.Argv = args
			config.Timeout = time.Duration(timeoutSeconds) * time.Second

			for _, service := range services {
				address, err := syntax.ParseAddress(service)
				if err != nil {
					report(err)
					os.Exit(1)
				}
				config.Addresses = append(config.Addresses, address)
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
	rootCommand.Flags().BoolVarP(&config.Quiet, "quiet", "q",
		false, "do not output any status messages")

	rootCommand.SetArgs(args)

	if err := rootCommand.Execute(); err != nil {
		report(err)
		os.Exit(1)
	}

	return config
}
