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
)

func init() {
	mk.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "install IMPORT_PATH as APP",
		Short:     "install a new app",
		Long: `
  Install a new Cider application from IMPORT_PATH and name it APP.

  IMPORT_PATH is a valid URL defining the repository holding the application
  sources and Cider configuration. The following schemes are available:

    * git+ssh   - clones a Git repository over SSH,
                  URL fragment is treated as Git ref.
	              example: git+ssh://git@github.com:cider/cider.git#develop

    * git+https - clones a Git repository over HTTPS,
                  URL fragment is treated as Git ref.
                  example: git+https://github.com/cider/cider-demo-webapp#master

    * git+file  - clones a Git repository from a local repository,
                  URL fragment is treated as Git ref.
                  example: git+file:///home/tchap/src/mk#develop
        `,
		Action: runInstall,
	})
}

func runInstall(cmd *gocli.Command, args []string) {
	if len(args) != 3 || args[1] != "as" {
		cmd.Usage()
		os.Exit(2)
	}

	if err := _runInstall(args[0], args[2]); err != nil {
		color.Fprintf(os.Stderr, "\n@{r}Error: %s@{|}\n", err)
		os.Exit(1)
	}

	color.Println("\n@{g}Success@{|}")
}

func _runInstall(url string, alias string) error {
	// Get the Cider management token.
	token, err := GetManagementToken()
	if err != nil {
		return err
	}

	// Send the install request to the server.
	var reply data.InstallReply
	err = SendRequest("Cider.Apps.Install", &data.InstallArgs{
		Token:      token,
		Alias:      alias,
		Repository: url,
	}, &reply)
	if err != nil {
		return err
	}
	if reply.Error != "" {
		return errors.New(reply.Error)
	}

	return nil
}
