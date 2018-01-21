// +build choria

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/mcorpc"
	"github.com/choria-io/go-choria/server"
	"github.com/choria-io/go-choria/server/agents"
	"github.com/choria-io/go-protocol/protocol"
	"github.com/sirupsen/logrus"
)

var (
	cserver      *server.Instance
	cfg          *choria.Config
	fw           *choria.Framework
	middleware   string
	enableChoria = true
)

// reply structure for the RPC agent
type switchReply struct {
	Message string `json:"message"`
	Paused  bool   `json:"paused"`
}

func registerAgent(ctx context.Context) error {
	// mcollective like metadata
	m := &agents.Metadata{
		Name:        "snowman",
		Description: "Snowman Control",
		Author:      "R.I.Pienaar <rip@devco.net>",
		Version:     "0.0.1",
		License:     "Apache-2.0",
		Timeout:     1,
		URL:         "https://www.devco.net/",
	}

	agent := mcorpc.New("snowman", m, fw, logrus.WithFields(logrus.Fields{"agent": "snowman"}))
	agent.MustRegisterAction("switch", switchAction)

	// adds the agent to the running instance of the server
	// this has to happen after its initial connect
	return cserver.RegisterAgent(ctx, "snowman", agent)
}

func switchAction(ctx context.Context, req *mcorpc.Request, reply *mcorpc.Reply, agent *mcorpc.Agent, conn choria.ConnectorInfo) {
	paused = !paused

	reply.Data = &switchReply{fmt.Sprintf("Flipped the snowman switch"), paused}
}

func setupChoria() (*server.Instance, error) {
	var err error

	cfg, err = choria.NewConfig("/dev/null")
	if err != nil {
		return nil, err
	}

	cfg.Choria.MiddlewareHosts = strings.Split(middleware, ",")
	if len(cfg.Choria.MiddlewareHosts) == 0 {
		return nil, errors.New("No middleware configured, cannot start choria")
	}

	// basic setup thats needed because there is no config at all
	cfg.LogLevel = "debug"
	cfg.MainCollective = "snowmen"
	cfg.Collectives = []string{"snowmen"}

	// disable TLS so this works on the plain text demo.nats.io
	cfg.DisableTLS = true
	protocol.Secure = "false"

	fw, err = choria.NewWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	return server.NewInstance(fw)
}

func runChoria() {
	if !enableChoria || middleware == "" {
		logrus.Errorf("Choria is disabled or no middleware have been compiled in")
		return
	}

	wg := &sync.WaitGroup{}
	ctx := context.Background()

	cserver, err = setupChoria()
	if err != nil {
		logrus.Errorf("Could not start choria: %s", err)
		return
	}

	wg.Add(1)
	cserver.Run(ctx, wg)

	err := registerAgent(ctx)
	if err != nil {
		logrus.Errorf("Could not register snowman agent: %s", err)
	}
}
