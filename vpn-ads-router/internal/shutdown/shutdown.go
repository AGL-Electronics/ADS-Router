package shutdown

import (
	"context"
	"net"
	"sync"
)

//this package contains the gracefull shutdown of the go routines for the proxy, without this memory leaks could manifest

type Server struct {
	listener net.Listener
	Wg       sync.WaitGroup
	Ctx      context.Context
	Cancel   context.CancelFunc
}

func NewServer(Addr string) (*Server, error) {

}
