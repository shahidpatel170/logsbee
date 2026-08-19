package service

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/shahidpatel170/logsbee/internal/domain"
)

// DetectHostEnvironment examines host binaries to match the correct logging subsystem.
// It searches the host's PATH execution variables for known system daemons.
func DetectHostEnvironment() string {
	// 1. Probe for native systemd tracking binaries (Fedora, Arch, Ubuntu, Debian)
	if _, err := exec.LookPath("journalctl"); err == nil {
		return "systemd"
	}

	// 2. Probe for Alpine Linux OpenRC circular memory log utilities
	if _, err := exec.LookPath("logread"); err == nil {
		return "openrc"
	}

	// 3. Fallback signature targeting Runit/standard plaintext targets (Void Linux/ Legacy)
	return "syslog"
}

// Orchestrator coordinates asynchronous log streams from workers into database transations.
type Orchestrator struct {
	worker domain.LogBeeWorker
}

// NewOrchestrartor instantiates a decoupled service runner bound to a polymorphic strategy.
func NewOrchestrartor(w domain.LogBeeWorker) *Orchestrator {
	return &Orchestrator{worker: w}
}

// StartIngestion opens the ingestion pipeline channels and boots the current harvester strategy.
func (o *Orchestrator) StartIngestion(ctx context.Context) error {

	stream := make(chan domain.LogPayload)
	defer close(stream)

	// Spawn the concrete worker loop inside an isolated concurrent goroutine
	go func() {
		if err := o.worker.Harvest(ctx, stream); err != nil {
			fmt.Printf("[Orchestrator Error] worker harvest failure: %v\n", err)
		}
	}()

	// infinite select monitoring loop processing active incoming log streams
	for {
		select {
		case <-ctx.Done():
			// Handle clean environment contextual cancel triggers
			return ctx.Err()
		case payload, ok := <-stream:
			if !ok {
				return nil
			}
			// This is where our incoming records hit the database.
			// For now we print a stub until we wire up the pgx repository pool.
			fmt.Printf("[Harvested Event] [%s] %s -> %s \n", payload.Severity, payload.Component, payload.Message)
		}
	}
}
