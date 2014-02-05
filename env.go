// Copyright (c) 2013 The ciderapp AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/cider/cider/apps/data"
	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"text/tabwriter"
)

func init() {
	ciderapp.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "env ALIAS",
		Short:     "show app variable values",
		Long: `
  Show the environmental variables as defined for application ALIAS.
        `,
		Action: runEnv,
	})
}

func runEnv(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runEnv(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}
}

func _runEnv(alias string) error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the status request to the server.
	var reply data.EnvReply
	err = SendRequest("Cider.Apps.Env", &data.EnvArgs{
		Token: token,
		Alias: alias,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	// Print the environment.
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tw.Flush()

	for k, v := range reply.Vars {
		if v.Value != "" {
			fmt.Fprintf(tw, "%s\t%v\n", k, v.Value)
		}
	}

	return nil
}
