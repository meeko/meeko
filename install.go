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
	app.MustRegisterSubcommand(&gocli.Command{
		UsageLine: "install IMPORT_PATH as ALIAS",
		Short:     "install a new agent",
		Long: `
  Install a new Meeko agent from IMPORT_PATH and name it ALIAS.

  IMPORT_PATH is a valid URL defining the repository holding the agent sources
  and its Meeko configuration. The following schemes are available:

    * git+ssh   - clones a Git repository over SSH,
                  URL fragment is treated as Git ref.
	              example: git+ssh://git@github.com/meeko/meeko#develop

    * git+https - clones a Git repository over HTTPS,
                  URL fragment is treated as Git ref.
                  example: git+https://github.com/meeko/meeko#master

    * git+file  - clones a Git repository from a local repository,
                  URL fragment is treated as Git ref.
                  example: git+file:///home/foobar/src/meeko#develop
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
	// Get the config file.
	cfg, err := LoadConfig(flagConfig)
	if err != nil {
		return err
	}

	// Send the install request to the server.
	var reply data.InstallReply
	err = SendRequest(cfg.Address, cfg.AccessToken, MethodInstall, &data.InstallArgs{
		Token:      []byte(cfg.ManagementToken),
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
