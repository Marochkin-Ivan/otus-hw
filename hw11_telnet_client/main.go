package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout. Default: 10s")
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Println(errors.New("not enough arguments to start telnet"))
		return
	}

	// Place your code here,
	// P.S. Do not rush to throw context down, think if it is useful with blocking operation?
	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	telnet := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	err := telnet.Connect()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer telnet.Close()

	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err = telnet.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "...Receive error: %v\n", err)
			stop() // Завершаем контекст при ошибке
		}
	}()

	if err = telnet.Send(); err != nil {
		fmt.Fprintf(os.Stderr, "...Send error: %v\n", err)
		stop() // Завершаем контекст при ошибке
	}
}
