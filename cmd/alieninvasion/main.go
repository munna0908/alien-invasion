package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/munna0908/alien-invasion/cmd/alieninvasion/cli"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	cli.Execute(closeCh)
}
