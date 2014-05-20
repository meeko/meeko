// Copyright (c) 2013 The meeko AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"

	"github.com/meeko/meekod/supervisor/data"

	"github.com/tchap/gocli"
	"github.com/wsxiaoys/terminal/color"
)

func init() {
	cmd := &gocli.Command{
		UsageLine: `
  set VAR for ALIAS to VALUE
  set -ask VAR for ALIAS
  (TBD) set -from FROM_ALIAS [-ours|-theirs] [VAR] for ALIAS
		`,
		Short: "set agent variable",
		Long: `
  Set environment variables for ALIAS.

  The first form simply sets VAR to VALUE.

  If the second form is used, the user is asked for the value with terminal
  echo being turned off. This is useful for passwords and other secret stuff.

  The last form tells Meeko to copy variables from some other agent FROM_AGENT.
  If VAR is skipped, the whole environment is copied. Full copy will fail if it
  would overwrite an existing variable. This behaviour can be altered to either
  prioritize ALIAS or FROM_ALIAS using -ours or -theirs respectively.
        `,
		Action: runSet,
	}
	cmd.Flags.BoolVar(&fsetAsk, "ask", fsetAsk, "ask for the value interactively, without echo")
	cmd.Flags.StringVar(&fsetFrom, "from", fsetFrom, "copy value from another agent")
	cmd.Flags.BoolVar(&fsetOurs, "ours", fsetOurs, "pick the value from ALIAS on conflict")
	cmd.Flags.BoolVar(&fsetTheirs, "theirs", fsetTheirs, "pick the value from FROM_ALIAS on conflict")

	app.MustRegisterSubcommand(cmd)
}

var (
	fsetAsk    bool
	fsetFrom   string
	fsetOurs   bool
	fsetTheirs bool
)

func runSet(cmd *gocli.Command, args []string) {
	switch {
	case fsetAsk && fsetFrom != "":
		cmd.Usage()
		os.Exit(2)
	case fsetAsk:
		if len(args) != 3 || args[1] != "for" {
			cmd.Usage()
			os.Exit(2)
		}
	case fsetFrom != "":
		if fsetOurs && fsetTheirs {
			cmd.Usage()
			os.Exit(2)
		}
		if !(len(args) == 2 && args[0] == "for") && !(len(args) == 3 && args[1] == "for") {
			cmd.Usage()
			os.Exit(2)
		}
	default:
		if len(args) != 5 || args[1] != "for" || args[3] != "to" {
			cmd.Usage()
			os.Exit(2)
		}
	}

	if err := _runSet(args); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runSet(args []string) error {
	// Read the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Ask for the value if requested.
	sargs := &data.SetArgs{
		Token: []byte(cfg.ManagementToken),
	}

	switch {
	case fsetAsk:
		// -ask VAR for ALIAS
		var reply string
		reply, err = promptUser("Insert the value of "+args[0]+": ", true)
		if err != nil {
			return err
		}
		sargs.Variable = args[0]
		sargs.Alias = args[2]
		sargs.Value = reply
	case fsetFrom != "":
		sargs.CopyFrom = fsetFrom
		if fsetOurs {
			sargs.MergeMode = "ours"
		}
		if fsetTheirs {
			sargs.MergeMode = "theirs"
		}
		if len(args) == 2 {
			// -from FROM_ALIAS for ALIAS
			sargs.Alias = args[1]
		} else {
			// -from FROM_ALIAS VAR for ALIAS
			sargs.Variable = args[0]
			sargs.Alias = args[2]
		}
	default:
		// VAR for ALIAS to VALUE
		sargs.Variable = args[0]
		sargs.Alias = args[2]
		sargs.Value = args[4]
	}

	// Send the request to the server.
	var reply data.SetReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodSet, sargs, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
