// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"time"

	"github.com/meeko/meekod/supervisor/data"

	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
)

func init() {
	subcmd := &gocli.Command{
		UsageLine: "stop [-timeout=TIMEOUT] ALIAS",
		Short:     "stop a running agent",
		Long: `
  Stop a Meeko agent, if it is actually running. Otherwise return an error.
        `,
		Action: runStop,
	}
	subcmd.Flags.DurationVar(&fstopTimeout, "timeout", fstopTimeout, "kill the agent after TIMEOUT")

	app.MustRegisterSubcommand(subcmd)
}

var fstopTimeout time.Duration = -1

func runStop(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runStop(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runStop(alias string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Send the clone request to the server.
	var reply data.StopReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodStop, &data.StopArgs{
		Token:   cfg.ManagementToken,
		Alias:   alias,
		Timeout: fstopTimeout,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
