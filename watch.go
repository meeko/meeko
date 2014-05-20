// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"

	"github.com/meeko/meekod/broker/services/logging"
	"github.com/meeko/meekod/supervisor/data"

	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
)

func init() {
	subcmd := &gocli.Command{
		UsageLine: "watch [-level=LEVEL] ALIAS",
		Short:     "stream agent logs",
		Long: `
  Stream agent logs of ALIAS to stdout. Just for completion, it is not possible
  to view older logs, this subcommand really only prints log entries as they
  are emitted by the agent.

  Log levels available for the level flag are
    0) unset (default)
	1) trace
	2) debug
	3) info
	4) warning
	5) error
	6) critical

  Choosing particular log level means that all the levels higher in the list
  are included as well, e.g warning level enables error and critical as well.
        `,
		Action: runWatch,
	}
	subcmd.Flags.StringVar(&watchLevel, "level", watchLevel, "filter logs by log level")

	app.MustRegisterSubcommand(subcmd)
}

var watchLevel = "unset"

func runWatch(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runWatch(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}
}

func _runWatch(alias string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	level, err := logging.ParseLogLevel(watchLevel)
	if err != nil {
		return err
	}

	// Start streaming the logs. This command will remain active until
	// interrupted by the user.
	var reply data.WatchReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodWatch, &data.WatchArgs{
		Token: []byte(cfg.ManagementToken),
		Alias: alias,
		Level: uint32(level),
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
