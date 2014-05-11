package main

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"os/signal"

	"github.com/cider/go-cider/cider/services/rpc"
	zrpc "github.com/cider/go-cider/cider/transports/zmq3/rpc"
	zmq "github.com/pebbe/zmq3"
	"github.com/wsxiaoys/terminal/color"
)

func SendRequest(method string, args interface{}, reply interface{}) (err error) {
	defer zmq.Term()

	debug("Connecting to Cider ... ")
	client, err := rpc.NewService(func() (rpc.Transport, error) {
		factory := zrpc.NewTransportFactory()
		factory.MustReadConfigFromEnv("CIDER_ZMQ3_RPC_")

		if fendpoint != "" {
			factory.Endpoint = fendpoint
		}

		factory.MustBeFullyConfigured()
		return factory.NewTransport("mk#" + mustRandomString())
	})
	if err != nil {
		debug(FAIL)
		return err
	}
	debug(OK)
	defer client.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	req := client.NewRemoteCall(method, args)
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
		case <-client.Closed():
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
