package main

import (
	"fmt"
	"os"

	bcapp "github.com/abdev/fabcar/blockchain_app"
	"github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
)

func startServer() error {
	app := bcapp.NewBlockChainApplication()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// Start the listener
	srv, err := server.NewServer("tcp://0.0.0.0:46658", "socket", app)
	logger.Info("start the server")
	if err != nil {
		return err
	}

	//srv.SetLogger(logger)

	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

	return nil
}

func main() {

	err := startServer()

	if err != nil {
		panic(fmt.Sprintf("Could not start server. Error: %v", err))
	}

}
