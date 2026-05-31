package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LevantateLabs/basaltrouter/internal/platform/server"
)

func main() {
	ctx := context.Background()
	res, err := server.Bootstrap(ctx, "gateway")
	if err != nil {
		fmt.Fprintf(os.Stderr, "gateway: bootstrap: %v\n", err)
		os.Exit(1)
	}
	defer res.Close(ctx)

	if err := server.RunHTTPServer(ctx, res, "gateway"); err != nil {
		res.Logger.Error("gateway exited with error", "error", err)
		os.Exit(1)
	}
}
