package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LevantateLabs/basaltrouter/internal/platform/server"
)

func main() {
	ctx := context.Background()
	res, err := server.Bootstrap(ctx, "api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "api: bootstrap: %v\n", err)
		os.Exit(1)
	}
	defer res.Close(ctx)

	if err := server.RunHTTPServer(ctx, res, "api"); err != nil {
		res.Logger.Error("api exited with error", "error", err)
		os.Exit(1)
	}
}
