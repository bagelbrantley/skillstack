package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bagelbrantley/skillstack/cmd"
	"github.com/bagelbrantley/skillstack/pkg/style"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cmd.ExecuteContext(ctx); err != nil {
		style.PrintError(os.Stderr, err.Error())
		os.Exit(1)
	}
}
