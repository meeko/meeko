// Copyright (c) 2013 The mk AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"os/signal"

	"github.com/meeko/go-meeko/meeko/services/rpc"
	transport "github.com/meeko/go-meeko/meeko/transports/websocket/rpc"

	"code.google.com/p/go.net/websocket"
	"github.com/wsxiaoys/terminal/color"
)

const (
	MethodEnv     = "Meeko.Agent.Env"
	MethodInfo    = "Meeko.Agent.Info"
	MethodInstall = "Meeko.Agent.Install"
	MethodList    = "Meeko.Agent.List"
	MethodRemove  = "Meeko.Agent.Remove"
	MethodRestart = "Meeko.Agent.Restart"
	MethodSet     = "Meeko.Agent.Set"
	MethodStart   = "Meeko.Agent.Start"
	MethodStatus  = "Meeko.Agent.Status"
	MethodStop    = "Meeko.Agent.Stop"
	MethodUnset   = "Meeko.Agent.Unset"
	MethodUpgrade = "Meeko.Agent.Upgrade"
	MethodWatch   = "Meeko.Agent.Watch"
)

const AccessTokenHeader = "X-Meeko-Token"

func SendRequest(address, token, method string, args, reply interface{}) error {
	// Connect to Meeko.
	debug("Connecting to Meeko ... ")
	service, err := rpc.NewService(func() (rpc.Transport, error) {
		factory := transport.NewTransportFactory()
		factory.Server = address
		factory.Origin = "http://localhost"
		factory.WSConfigFunc = func(wsConfig *websocket.Config) {
			wsConfig.Header.Set(AccessTokenHeader, token)
		}
		factory.MustBeFullyConfigured()
		return factory.NewTransport("meeko#" + mustRandomString())
	})
	if err != nil {
		debug(FAIL)
		return err
	}
	debug(OK)
	defer service.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Dispatch the request and stream the output to the console.
	req := service.NewRemoteCall(method, args)
	req.Stdout = os.Stdout
	req.Stderr = os.Stderr

	debug("Calling ", method, " ... ")
	req.GoExecute()
	debug(OK)

	go func() {
		select {
		case <-signalCh:
			color.Println("@{c}<<< @{r}Interrupting remote call ...")
			req.Interrupt()
		case <-service.Closed():
		}
	}()

	debug(">>> Request output\n")
	err = req.Wait()
	debug("<<< Request output\n")
	if err != nil {
		return err
	}

	debug("Return code: ", req.ReturnCode(), "\n")
	err = req.UnmarshalReturnValue(&reply)
	debug("Return value: ", reply, "\n")
	return err
}

func mustRandomString() string {
	buf := make([]byte, 10)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(buf)
}
