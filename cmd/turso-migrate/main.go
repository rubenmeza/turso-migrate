package main

import (
	"fmt"
	"os"

	"github.com/rubenmeza/turso-migrate/internal/cli"
)

func main() {
	app := cli.NewApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
