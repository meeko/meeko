// Copyright (c) 2013 The ciderapp AUTHORS
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
)

func init() {
	cmd := &gocli.Command{
		UsageLine: `
  set VAR for ALIAS to VALUE
  set -ask VAR for ALIAS
  set -from FROM_ALIAS [-ours|-theirs] [VAR] for ALIAS
		`,
		Short: "set app variable",
		Long: `
  Set application environmental variables for ALIAS.

  The first form simply sets VAR to VALUE.

  If the second form is used, the user is asked for the value with terminal
  echo being turned off. This is useful for passwords and other secret stuff.

  The last form tells the server to copy variables from some other application
  FROM_ALIAS. If VAR is skipped, the whole environment is copied. Full copy
  will fail if it would overwrite an existing variable. This behaviour can be
  altered to either prioritize ALIAS or FROM_ALIAS using -ours or -theirs
  respectively.
        `,
		Action: runSet,
	}
	cmd.Flags.BoolVar(&fsetAsk, "ask", fsetAsk, "ask for the value interactively, without echo")
	cmd.Flags.StringVar(&fsetFrom, "from", fsetFrom, "copy value from another app")
	cmd.Flags.BoolVar(&fsetOurs, "ours", fsetOurs, "pick the value from ALIAS on conflict")
	cmd.Flags.BoolVar(&fsetTheirs, "theirs", fsetTheirs, "pick the value from FROM_ALIAS on conflict")

	ciderapp.MustRegisterSubcommand(cmd)
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
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Ask for the value if requested.
	sargs := &data.SetArgs{
		Token: token,
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
	if err := SendRequest("Cider.Apps.Set", sargs, &reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
