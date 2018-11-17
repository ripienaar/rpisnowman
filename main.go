package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry
var Version = "2.0.0"
var mu = &sync.Mutex{}
var wg = &sync.WaitGroup{}
var loglevel = "info"

var ctx context.Context
var cancel func()

func main() {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	loglevel = os.Getenv("SNOWMAN_LOGLEVEL")
	if loglevel == "" {
		loglevel = "info"
	}

	log = logrus.NewEntry(logrus.New())

	switch loglevel {
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		log.Logger.SetLevel(logrus.WarnLevel)
	default:
		loglevel = "info"
		log.Logger.SetLevel(logrus.InfoLevel)
	}

	name := os.Getenv("SNOWMAN_NAME")
	if name == "" {
		name = "snowman"
	}

	log.Infof("Ryanteck RTK-000-00A GPIO Snowman '%s' version %s starting", name, Version)

	go interruptWatcher()

	man := NewSnowMan(name, log)

	wg.Add(1)
	go man.StartBackplane(ctx, wg)

	err := man.Open()
	if err != nil {
		fmt.Printf("Could not open rpi: %s", err)
		os.Exit(1)
	}
	defer man.Close()

	wg.Add(1)
	go man.Run(ctx, wg)

	wg.Wait()
}

func interruptWatcher() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				log.Warnf("Shutting down on %s", sig)
				cancel()
			case syscall.SIGQUIT:
				dumpGoRoutines()
			}
		case <-ctx.Done():
			return
		}
	}
}

func dumpGoRoutines() {
	mu.Lock()
	defer mu.Unlock()

	outname := filepath.Join(os.TempDir(), fmt.Sprintf("snowman-threaddump-%d-%d.txt", os.Getpid(), time.Now().UnixNano()))

	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)

	err := ioutil.WriteFile(outname, buf[:stacklen], 0644)
	if err != nil {
		log.Errorf("Could not produce thread dump: %s", err)
		return
	}

	log.Warnf("Produced thread dump to %s", outname)
}
