// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/meeko/meekod/supervisor/data"

	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
)

func init() {
	app.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "list",
		Short:     "list installed agents",
		Long: `
  List all agents installed on the target Meeko instance.
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
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Send the request to the server.
	var reply data.ListReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodList, &data.ListArgs{
		Token: cfg.ManagementToken,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	// Print the agent list.
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, "ALIAS\tAGENT NAME\tVERSION\tENABLED")
	fmt.Fprintln(tw, "=====\t==========\t=======\t=======")

	for _, agent := range reply.Agents {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", agent.Alias, agent.Name, agent.Version, agent.Enabled)
	}

	fmt.Fprintln(tw, "")
	return nil
}
