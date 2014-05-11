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
		UsageLine: "upgrade ALIAS",
		Short:     "upgrade an installed agent",
		Long: `
  Upgrading an agent means that its sources are updated and the executable
  is rebuilt, replacing the old executable if successful. Then the agent
  process is restarted if running.
        `,
		Action: runUpgrade,
	})
}

func runUpgrade(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runUpgrade(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runUpgrade(alias string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Send the clone request to the server.
	var reply data.UpgradeReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodUpgrade, &data.UpgradeArgs{
		Token: cfg.ManagementToken,
		Alias: alias,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
