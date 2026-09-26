package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		slog.Error("Missing required params, usage: go-telnet [--timeout] <host> <port>")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

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

	// from stdin
	go func() {
		for {
			if err := client.Send(); err != nil {
				if errors.Is(err, io.EOF) {
					slog.Error("EOF:", "error", err)
					break
				}
				slog.Error("send error:", "error", err)
			}
		}
		cancel()
	}()

	// to stdout
	go func() {
		for {
			if err := client.Receive(); err != nil {
				if errors.Is(err, io.EOF) {
					slog.Error("connection was closed by server:", "error", err)
					break
				}
				slog.Error("receive error:", "error", err)
			}
		}
		cancel()
	}()

	<-ctx.Done()
}
