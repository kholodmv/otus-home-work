package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()
	if len(flag.Args()) < 2 {
		os.Exit(1)
	}
	address := fmt.Sprintf("%s:%s", flag.Args()[0], flag.Args()[1])
	client, err := NewTelnetClient(address, *timeout, os.Stdin, os.Stderr)
	if err != nil {
		log.Printf("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer client.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		log.Printf("Received SIGINT, closing connection...")
		client.Close()
		os.Exit(0)
	}()

	log.Printf("...Connected to %s\n", address)

	go func() {
		if err := client.Receive(); err != nil && !client.closed {
			log.Printf("...Connection was closed by peer")
			client.Close()
			return
		}
	}()

	if err := client.Send(); err != nil && !client.closed {
		log.Printf("...Error sending to server: %v\n", err)
		client.Close()
		return
	}

	if client.closed {
		log.Printf("...EOF")
		return
	}
}
