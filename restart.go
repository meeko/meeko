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
		UsageLine: "restart ALIAS",
		Short:     "restart a running app",
		Long: `
 Restart a running application.
 
 This action can only be performed on a running application, so an error is
 returned if the application is not running.
        `,
		Action: runRestart,
	})
}

func runRestart(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runRestart(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runRestart(alias string) error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the clone request to the server.
	var reply data.RestartReply
	err = SendRequest("Cider.Apps.Restart", &data.RestartArgs{
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
