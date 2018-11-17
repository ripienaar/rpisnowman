package main

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/choria-io/go-backplane/backplane"
	"github.com/sirupsen/logrus"
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
		backplane.ManageLogLevel(s),
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

// SetLogLevel implements backplane.LogLevelSetable
func (s *SnowMan) SetLogLevel(level backplane.LogLevel) {
	switch level {
	case backplane.InfoLevel:
		s.log.Logger.SetLevel(logrus.InfoLevel)
		s.log.Infof("Set log level to info")
	case backplane.WarnLevel:
		s.log.Logger.SetLevel(logrus.WarnLevel)
		s.log.Warnf("Set log level to warn")
	default:
		s.log.Logger.SetLevel(logrus.DebugLevel)
		s.log.Debugf("Set log level to debug")
	}
}

// GetLogLevel implements backplane.LogLevelSetable
func (s *SnowMan) GetLogLevel() backplane.LogLevel {
	switch s.log.Logger.Level {
	case logrus.DebugLevel:
		return backplane.DebugLevel
	case logrus.InfoLevel:
		return backplane.InfoLevel
	case logrus.WarnLevel:
		return backplane.WarnLevel
	default:
		return backplane.CriticalLevel
	}
}
