package domain

import (
	"context"
	"time"
)

// LogPayload represents a unified, structured log entry inside LogsBee.
type LogPayload struct {
	Timestamp   time.Time `json:"timestamp"`
	Severity    string    `json:"severity"`
	Component   string    `json:"component"`
	SystemdUnit string    `json:"systemd_unit,omitempty"`
	Message     string    `json:"message"`
	RawData     []byte    `json:"raw_data"`
}

// LogBeeWorker defines the structural rules for host platform data extraction.
// Every target distribution logging system must implement this interface contract.
type LogBeeWorker interface {
	// Harvest actively streams structural system logs and pipes entries into the channel.
	// It relies on a context object for graceful termination and cleanup.
	Harvest(ctx context.Context, stream chan<- LogPayload) error

	// GetProviderSignature exposes the tracking driver identity ('journald', 'syslog_file', etc.)
	GetProviderSignature() string
}
