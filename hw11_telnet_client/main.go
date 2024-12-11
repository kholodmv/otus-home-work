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
		log.Fatalf("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer client.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		log.Fatalf("Received SIGINT, closing connection...")
		client.Close()
		os.Exit(0)
	}()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	go func() {
		if err := client.Receive(); err != nil && !client.closed {
			log.Fatalf("...Connection was closed by peer")
			client.Close()
			os.Exit(0)
		}
	}()

	if err := client.Send(); err != nil && !client.closed {
		log.Fatalf("...Error sending to server: %v\n", err)
		client.Close()
		os.Exit(0)
	}

	if client.closed {
		log.Fatalf("...EOF")
		os.Exit(0)
	}
}
