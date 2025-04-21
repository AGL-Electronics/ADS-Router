package router

import (
	"net"
)

type NetID [6]byte

type Message struct {
	SourceNetId NetID
	TargetnetId NetID
	invokeId    uint32
	payload     []byte
	Conn        net.Conn
}
