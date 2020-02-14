package main

import (
	"fmt"
	"os"
	"os/signal"
)

func handleSignals() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)
	go func() {
		fmt.Println("Waiting for signals")
		s := <-sigc
		// Log event
		s.String()
		fmt.Println(fmt.Sprintf("RECEIVED SIGNAL: %s", s.String()))
	}()
}
