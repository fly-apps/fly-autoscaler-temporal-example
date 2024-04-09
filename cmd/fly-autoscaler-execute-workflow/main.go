package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	fate "github.com/fly-apps/fly-autoscaler-temporal-example"
	"go.temporal.io/sdk/client"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	count := flag.Int("count", 1, "number of workflows to start")
	every := flag.Duration("every", 1*time.Second, "frequency to create new workflows")
	duration := flag.Duration("duration", 2*time.Second, "how long each workflow will run for")
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

	// Create new N workflows on a given interval.
	slog.Info("initiating loop", slog.Int("n", *count), slog.Duration("every", *every))

	ticker := time.NewTicker(*every)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for i := 0; i < *count; i++ {
				we, err := c.ExecuteWorkflow(
					context.Background(),
					client.StartWorkflowOptions{TaskQueue: fate.TaskQueue},
					fate.ExampleWorkflow,
					*duration,
				)
				if err != nil {
					return fmt.Errorf("cannot execute workflow: %w", err)
				}
				slog.Info("workflow started", slog.String("id", we.GetRunID()), slog.Duration("duration", *duration))
			}
		}
	}
}

type ActivityResult struct {
	Msg string
}
