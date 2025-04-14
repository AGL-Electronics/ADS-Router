package proxy

import (
	"fmt"
	"net"

	"vpn-ads-router/pkg/logger"
)

type ClientMsg struct {
	SourceNetId string
	payload     []byte
}

var connectionlogger = logger.GetLogger()
var IncommingConnChan = make(chan ClientMsg, 100)

// StartListener opens a TCP socket on the given address and starts accepting connections.
// Each client gets its own goroutine and feeds messages to the router via IncomingChan.
func StartListener(Address string) {
	ln, err := net.Listen("tcp", Address)
	if err != nil {
		connectionlogger.Error(logger.ComponentProxy, "Error starting listener on %s: %v", Address, err)
	}
	connectionlogger.Info(logger.ComponentProxy, "Proxy listening on %s", Address)

	for {
		Conn, err := ln.Accept()
		if err != nil {
			connectionlogger.Error(logger.ComponentProxy, "Error accepting connection: %v", err)
			continue
		}
		connectionlogger.Info(logger.ComponentProxy, "Incomming connection from %s", Conn.RemoteAddr())
		go HandleClient(Conn)
	}
}

func ParseSourceNetId(payload []byte) string {
	if len(payload) < 16 {
		return "0.0.0.0.0.0" //default value for invalid payloads
	}

	b := payload[10:16] //get the source net id from the payload, this is a 6 byte value
	return fmt.Sprintf("%d.%d.%d.%d.%d.%d", b[0], b[1], b[2], b[3], b[4], b[5])
}


// handleClient handles one incomming tcp connection from a client.
// it contnuesly reads packets and sends them to the sheduler queue.
func HandleClient(Conn net.Conn) {
	defer Conn.Close()

	remoteAddr := Conn.RemoteAddr().String()
	connectionlogger.Info(logger.ComponentProxy, "Handling connection from %s", remoteAddr)

	for {
		buf := make([]byte, 1024)
		n, err := Conn.Read(buf)
		if err != nil {
			connectionlogger.Error(logger.ComponentProxy, "Error reading from connection %s: %v", remoteAddr, err)
			return
		}

		NetID := ParseSourceNetId(buf[:n])
		connectionlogger.Debug(logger.ComponentProxy, "Received %d bytes from %s, NetID: %s", n, remoteAddr, NetID)

		ClientMsg := ClientMsg{
			SourceNetId: NetID,
			payload:     buf[:n],
		}

		IncommingConnChan <- ClientMsg //send the message to the router								THESE 2 REQUIRE TESTING
		connectionlogger.Debug(logger.ComponentProxy, "Sent message to router: %v", ClientMsg)
	}
}
