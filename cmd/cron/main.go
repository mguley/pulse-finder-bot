package main

import (
	"application"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle system signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	app := application.NewContainer()
	scheduler := app.CronScheduler.Get()

	go scheduler.Start(ctx)

	<-signalChan
	log.Println("Received termination signal. Shutting down...")
	scheduler.Stop()
}
