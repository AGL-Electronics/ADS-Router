package clearer

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"vpn-ads-router/pkg/logger"
)

//this package contains the gracefull shutdown of the go routines for the router
//without this memory leaks could manifest

type GraceFullShutdown struct {
	Wg     sync.WaitGroup
	Ctx    context.Context
	Cancel context.CancelFunc
}

var shutdownlogger = logger.GetLogger()

func New() *GraceFullShutdown {
	ctx, cancel := context.WithCancel(context.Background())

	Gs := &GraceFullShutdown{
		Ctx:    ctx,
		Cancel: cancel,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdownlogger.Info(logger.ComponentService, "Received shutdown signal, shutting down...")
		Gs.Cancel()
	}()
	return Gs
}

func (Gs *GraceFullShutdown) Add(delta int) {
	Gs.Wg.Add(delta)
}

func (Gs *GraceFullShutdown) Done() {
	Gs.Wg.Done()
}

func (Gs *GraceFullShutdown) Wait() {
	Gs.Wg.Wait()
}
