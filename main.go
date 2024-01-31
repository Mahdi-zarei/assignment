package main

import (
	"assignment/internal"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	svc := internal.Service{}
	svc.Start(ctx)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals

	fmt.Println("Signal received, shutting down gracefully...")
	cancel()
	time.Sleep(2 * time.Second)
	svc.Close()
}
