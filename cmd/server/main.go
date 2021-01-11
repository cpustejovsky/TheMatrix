package main

import (
	"os"
	"os/signal"
	"syscall"

	matrix "github.com/mailgun/TheMatrix"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infof("Listening on localhost:2319....")
	server, err := matrix.SpawnServer(matrix.ServerConfig{ListenAddress: "localhost:2319"})
	if err != nil {
		logrus.Errorf("while starting server: %s", err)
		os.Exit(1)
	}

	// Wait here for signals to clean up our mess
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for range c {
		logrus.Info("caught signal; shutting down")
		server.Close()
		os.Exit(0)
	}
}
