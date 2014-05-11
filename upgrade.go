// Copyright (c) 2013 The mk AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"github.com/cider/cider/apps/data"
	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
	"os"
)

func init() {
	mk.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "upgrade ALIAS",
		Short:     "upgrade an existing app",
		Long: `
  Upgrading an application means that its sources are updated and
  the application is rebuilt, replacing the old executable if successful.
  As the last step, the application process is restarted.
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
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the clone request to the server.
	var reply data.UpgradeReply
	err = SendRequest("Cider.Apps.Upgrade", &data.RemoveArgs{
		Token: token,
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
