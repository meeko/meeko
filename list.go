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
	mk.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "list",
		Short:     "list installed apps",
		Long: `
  List all applications installed on the local Cider instance.
        `,
		Action: runList,
	})
}

func runList(cmd *gocli.Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runList(); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %v@{|}\n", err)
		os.Exit(1)
	}
}

func _runList() error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the request to the server.
	var reply data.ListReply
	err = SendRequest("Cider.Apps.List", &data.ListArgs{
		Token: token,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	// Print the application list.
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, "ALIAS\tAPPLICATION NAME\tVERSION\tENABLED")
	fmt.Fprintln(tw, "=====\t================\t=======\t=======")

	for _, app := range reply.Apps {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", app.Alias, app.Name, app.Version, app.Enabled)
	}

	fmt.Fprintln(tw, "")
	return nil
}
