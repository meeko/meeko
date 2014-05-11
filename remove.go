// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"strings"

	"github.com/meeko/meekod/supervisor/data"

	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
)

func init() {
	app.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "remove ALIAS",
		Short:     "uninstall an agent",
		Long: `
  Removal of an agent means that all the files associated with it are
  deleted. The database records meet the same bitter end.
        `,
		Action: runRemove,
	})
}

func runRemove(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runRemove(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runRemove(alias string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Ask user for confirmation.
	answer, err := promptUser("Are you sure you want to proceed? [y/N]: ", false)
	if err != nil {
		return err
	}
	if strings.ToLower(answer) != "y" {
		return errors.New("Operation canceled")
	}

	// Send the clone request to the server.
	var reply data.RemoveReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodRemove, &data.RemoveArgs{
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
