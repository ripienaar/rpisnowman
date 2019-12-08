package main

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/choria-io/go-backplane/backplane"
)

// StartBackplane starts the management backplane if brokers are set
func (s *SnowMan) StartBackplane(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	brokers := os.Getenv("BACKPLANE_BROKERS")
	if brokers == "" {
		s.log.Infof("Backplane management is not enabled, now BACKPLANE_BROKERS set")
		return
	}

	opts := []backplane.Option{
		backplane.ManageInfoSource(s),
		backplane.ManagePausable(s),
	}

	config := &backplane.StandardConfiguration{
		AppName:  s.name,
		Loglevel: loglevel,
		Brokers:  strings.Split(brokers, ","),
		Authorization: backplane.Authorization{
			Insecure: true,
		},
	}

	_, err := backplane.Run(ctx, wg, config, opts...)
	if err != nil {
		s.log.Fatalf("Could not start backplane: %s", err)
	}
}

// FactData implements backplane.InfoSource
func (s *SnowMan) FactData() interface{} {
	data := make(map[string]string)
	data["name"] = s.name

	return data
}

// Version implements backplane.InfoSource
func (s *SnowMan) Version() string {
	return Version
}

// Pause implements backplane.Pausable
func (s *SnowMan) Pause() {
	s.log.Info("Pausing snowman display")
	s.paused = true
}

// Resume implements backplane.Pausable
func (s *SnowMan) Resume() {
	s.log.Info("Resuming snowman display")
	s.paused = false
}

// Flip implements backplane.Pausable
func (s *SnowMan) Flip() {
	s.log.Info("Toggling snowman display")
	s.paused = !s.paused
}

// Paused implements backplane.Pausable
func (s *SnowMan) Paused() bool {
	return s.paused
}
