package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/LevantateLabs/basaltrouter/internal/platform/server"
)

func main() {
	ctx := context.Background()
	res, err := server.Bootstrap(ctx, "worker")
	if err != nil {
		fmt.Fprintf(os.Stderr, "worker: bootstrap: %v\n", err)
		os.Exit(1)
	}
	defer res.Close(ctx)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.RunHTTPServer(ctx, res, "worker")
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				res.Logger.Error("worker http server exited with error", "error", err)
				os.Exit(1)
			}
			return
		case <-ticker.C:
			res.Logger.Debug("worker heartbeat")
		}
	}
}
