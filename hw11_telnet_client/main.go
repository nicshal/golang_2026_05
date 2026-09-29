package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultTimeout = time.Second * 10

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "connection timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		slog.Error("Missing required params, usage: go-telnet [--timeout] <host> <port>")
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		slog.Error("error connecting:", "error", err)
		os.Exit(1)
	}

	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			slog.Error("error client close:", "error", err)
		}
	}(client)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	// from stdin
	go func() {
		defer cancel()
		for {
			if err := client.Send(); err != nil {
				slog.Error("send error:", "error", err)
			} else {
				slog.Info("send success")
				break
			}
		}
	}()

	// to stdout
	go func() {
		defer cancel()
		for {
			if err := client.Receive(); err != nil {
				slog.Error("receive error:", "error", err)
			} else {
				slog.Info("receive success")
				break
			}
		}
	}()

	<-ctx.Done()
}
