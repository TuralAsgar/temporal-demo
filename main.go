package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	logTemporal "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Temporal client (connects to Frontend)
	c, err := client.Dial(client.Options{
		// Discard debug logs not to distract viewers of demo
		Logger: logTemporal.NewStructuredLogger(slog.New(slog.NewJSONHandler(io.Discard, nil))),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Worker
	w := worker.New(c, "lead-calling-task-queue", worker.Options{})

	w.RegisterWorkflow(LeadCallingWorkflow)

	w.RegisterActivity(PickLeadActivity)
	w.RegisterActivity(CallCustomerActivity)
	w.RegisterActivity(UpdateLeadActivity)
	w.RegisterActivity(CreateInteractionActivity)

	// API
	api := NewAPI(c)

	mux := http.NewServeMux()
	mux.HandleFunc("/start-call", api.startCallHandler)
	mux.HandleFunc("/3cx/webhook", api.threeCXWebhookHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run worker
	go func() {
		log.Println("Worker started")
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Run HTTP server
	go func() {
		log.Println("HTTP API listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for signal
	<-sigChan
	log.Println("Shutting down...")

	// Gracefully shutdown HTTP server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
	}

	// Stop worker
	w.Stop()
	log.Println("Shutdown complete")
}
