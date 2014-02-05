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
		UsageLine: "info ALIAS",
		Short:     "show app info",
		Long: `
  Show application info, which is basically the relevant ciderapp.json,
  just formated nicely. This command, however, does not only print the static
  data from ciderapp.json, but also the current application configuration,
  such as the environmental variables.
        `,
		Action: runInfo,
	})
}

func runInfo(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runInfo(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}
}

func _runInfo(alias string) error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return nil
	}

	// Send the status request to the server.
	var reply data.InfoReply
	err = SendRequest("Cider.Apps.Info", &data.InfoArgs{
		Token: token,
		Alias: alias,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	// Print the app info.
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tw.Flush()

	app := reply.App

	fmt.Fprintf(tw, "Alias:\t%s\n", app.Alias)
	fmt.Fprintf(tw, "App name:\t%s\n", app.Name)
	fmt.Fprintf(tw, "Version:\t%s\n", app.Version)
	fmt.Fprintf(tw, "Description:\t%s\n", app.Description)
	fmt.Fprintf(tw, "Repository:\t%s\n", app.RepositoryURL)
	if len(app.Vars) != 0 {
		fmt.Fprintf(tw, "Variables:\n")
		for k, v := range app.Vars {
			fmt.Fprintf(tw, "\t\tName:\t%s\n", k)
			fmt.Fprintf(tw, "\t\tUsage:\t%s\n", v.Usage)
			fmt.Fprintf(tw, "\t\tType:\t%s\n", v.Type)
			if v.Value == "" {
				if v.Optional {
					v.Value = color.Sprint("@{y}unset@{|}")
				}
				v.Value = color.Sprint("@{r}unset@{|}")
			}
			fmt.Fprintf(tw, "\t\tValue:\t%s\n\n", v.Value)
		}
	}

	return nil
}
