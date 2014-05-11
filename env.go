// Copyright (c) 2013 The mk AUTHORS
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
	subcmd := &gocli.Command{
		UsageLine: "env [-include_unset] ALIAS",
		Short:     "show app variable values",
		Long: `
  Show the environmental variables as defined for application ALIAS.

  Unset variables are not shown unless -include_unset is present.
        `,
		Action: runEnv,
	}
	subcmd.Flags.BoolVar(&fenvIncludeUnset, "include_unset", fenvIncludeUnset,
		"include unset variables in the output")

	mk.MustRegisterSubcommand(subcmd)
}

var fenvIncludeUnset bool

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
		if v.Value != "" || fenvIncludeUnset {
			if v.Value == "" {
				if v.Optional {
					v.Value = color.Sprint("@{y}<unset>@{|}")
				} else {
					v.Value = color.Sprint("@{r}<unset>@{|}")
				}
			}
			fmt.Fprintf(tw, "%s\t%v\n", k, v.Value)
		}
	}

	return nil
}
