package main

import (
	"fmt"
	"github.com/jwma/jump-jump/internal/app/cmd/server"
	"os"
)

func main() {
	addr := os.Getenv("J2_LANDING_ADDR")
	if addr == "" {
		_, _ = fmt.Fprint(os.Stderr, "missing J2_LANDING_ADDR environment variable")
		os.Exit(1)
	}

	err := server.RunLanding(addr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
