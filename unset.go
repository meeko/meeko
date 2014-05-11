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
	app.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "unset VAR for ALIAS",
		Short:     "unset agent variable",
		Long: `
  Unset an environment variables VAR defined for agent ALIAS.
        `,
		Action: runUnset,
	})
}

func runUnset(cmd *gocli.Command, args []string) {
	if len(args) != 3 || args[1] != "for" {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runUnset(args[0], args[2]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runUnset(variable string, alias string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Send the request to the server.
	var reply data.UnsetReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodUnset, &data.UnsetArgs{
		Token:    cfg.ManagementToken,
		Alias:    alias,
		Variable: variable,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
