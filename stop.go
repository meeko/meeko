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
	"time"
)

func init() {
	subcmd := &gocli.Command{
		UsageLine: "stop [-timeout] ALIAS",
		Short:     "stop a running app",
		Long: `
  Stop a Cider application, if it is actually running, that is.
  If that is not the case, this action is a NOOP.
        `,
		Action: runStop,
	}
	subcmd.Flags.DurationVar(&fstopTimeout, "timeout", fstopTimeout, "kill the app after timeout")

	mk.MustRegisterSubcommand(subcmd)
}

var fstopTimeout time.Duration = -1

func runStop(cmd *gocli.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runStop(args[0]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runStop(alias string) error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the clone request to the server.
	var reply data.StopReply
	err = SendRequest("Cider.Apps.Stop", &data.StopArgs{
		Token:   token,
		Alias:   alias,
		Timeout: fstopTimeout,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
