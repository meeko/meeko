// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"

	"github.com/meeko/meekod/supervisor/data"

	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
)

func init() {
	subcmd := &gocli.Command{
		UsageLine: "start [-watch] ALIAS",
		Short:     "start an agent",
		Long: `
  Start a Meeko agent.

  To start an agent, all the required agent variables must be set. When all the
  requirements are met, a new process is started with all the agent variables
  exported as environment variables.
        `,
		Action: runStart,
	}
	subcmd.Flags.BoolVar(&fstartWatch, "watch", fstartWatch, "start watching the agent")

	app.MustRegisterSubcommand(subcmd)
}

var fstartWatch bool

func runStart(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runStart(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runStart(alias string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Send the clone request to the server.
	var reply data.StartReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodStart, &data.StartArgs{
		Token: []byte(cfg.ManagementToken),
		Alias: alias,
		Watch: fstartWatch,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
