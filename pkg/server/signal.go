package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Callback func(signal os.Signal)

func newSystemContext(ctx context.Context, delay time.Duration, callbacks ...Callback) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		sig := <-sigint
		for _, cb := range callbacks {
			go cb(sig)
		}

		time.Sleep(delay)

		cancel()
	}()

	return ctx
}
