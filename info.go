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
		UsageLine: "info ALIAS",
		Short:     "show agent info",
		Long: `
  Show agent info, which is basically the relevant agent.json, just formatted
  nicely. This command, however, does not only print the static data from
  agent.json, but also the current agent configuration, such as the environment
  variables.
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
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return nil
	}

	// Send the status request to the server.
	var reply data.InfoReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodInfo, &data.InfoArgs{
		Token: []byte(cfg.ManagementToken),
		Alias: alias,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	// Print the agent info.
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tw.Flush()

	agent := reply.Agent

	fmt.Fprintf(tw, "Alias:\t%s\n", agent.Alias)
	fmt.Fprintf(tw, "Name:\t%s\n", agent.Name)
	fmt.Fprintf(tw, "Version:\t%s\n", agent.Version)
	fmt.Fprintf(tw, "Description:\t%s\n", agent.Description)
	fmt.Fprintf(tw, "Repository:\t%s\n", agent.Repository)
	if len(agent.Vars) != 0 {
		fmt.Fprintf(tw, "Variables:\n")
		for k, v := range agent.Vars {
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
