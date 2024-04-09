package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"os"

	fate "github.com/fly-apps/fly-autoscaler-temporal-example"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	concurrency := flag.Int("concurrency", 1, "concurrent activity execution size")
	flag.Parse()

	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	namespace := os.Getenv("TEMPORAL_NAMESPACE")
	certData := os.Getenv("TEMPORAL_TLS_CERT_DATA")
	keyData := os.Getenv("TEMPORAL_TLS_KEY_DATA")

	cert, err := tls.X509KeyPair([]byte(certData), []byte(keyData))
	if err != nil {
		return fmt.Errorf("load key pair: %w", err)
	}

	slog.Info("connecting to temporal cloud")

	// Add the cert to the tls certificates in the ConnectionOptions of the Client
	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
		},
	})
	if err != nil {
		return fmt.Errorf("cannot connect to Temporal Cloud: %w", err)
	}
	defer c.Close()

	slog.Info("registering worker")

	w := worker.New(c, fate.TaskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize: *concurrency,
	})
	w.RegisterWorkflow(fate.ExampleWorkflow)
	w.RegisterActivity(&fate.ExampleActivity{})

	slog.Info("starting worker")

	if err := w.Run(worker.InterruptCh()); err != nil {
		return fmt.Errorf("cannot start worker: %w", err)
	}
	return nil
}
