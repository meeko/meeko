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
		UsageLine: "status [ALIAS]",
		Short:     "show app status",
		Long: `
  When used without any argument, this command lists all installed applications
  together with their statuses. The status is one of the following:

    * stopped - the app is configured and can be started any time
    * running - the app is running
    * crashed - Cider tried to start the app process, but it crashed
    * killed  - Cider failed to stop the app and had to kill it

  When used with non-empty ALIAS, only the status of the chosen app is printed.
        `,
		Action: runStatus,
	})
}

func runStatus(cmd *gocli.Command, args []string) {
	if len(args) > 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runStatus(args); err != nil {
		os.Exit(1)
	}
}

func _runStatus(args []string) error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the status request to the server.
	var statArgs data.StatusArgs
	statArgs.Token = token
	if len(args) == 1 {
		statArgs.Alias = args[0]
	}

	var reply data.StatusReply
	err = SendRequest("Cider.Apps.Status", &statArgs, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	// Print the reply.
	if len(args) == 1 {
		fmt.Println(colorStatus(reply.Status))
		return nil
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, "ALIAS\tSTATUS")
	fmt.Fprintln(tw, "=====\t======")

	for k, v := range reply.Statuses {
		fmt.Fprintf(tw, "%s\t%s\n", k, colorStatus(v))
	}

	fmt.Fprintln(tw, "")
	return nil
}

func colorStatus(status string) string {
	switch status {
	case "stopped":
		return "stopped"
	case "running":
		return color.Sprint("@{g}running@{|}")
	case "crashed":
		return color.Sprint("@{r}crashed@{|}")
	case "killed":
		return color.Sprint("@{m}killed@{|}")
	}

	panic("Unknown status returned")
}
