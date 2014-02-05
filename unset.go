// Copyright (c) 2013 The ciderapp AUTHORS
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
	ciderapp.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "unset VAR for APP",
		Short:     "unset app variable",
		Long: `
  Unset an environmental variables VAR defined for application APP.
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
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the request to the server.
	var reply data.UnsetReply
	err = SendRequest("Cider.Apps.Unset", &data.UnsetArgs{
		Token:    token,
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
