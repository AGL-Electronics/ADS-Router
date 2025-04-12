package proxy

import (
	"io"
	"net"

	"vpn-ads-router/internal/network"
	"vpn-ads-router/pkg/logger"
)

var connectionlogger = logger.GetLogger()

func Handleconnection(ClientConn net.Conn) {

	defer ClientConn.Close()
	connectionlogger.Info(logger.ComponentService, "Incomming connection from %s", ClientConn.RemoteAddr())

	PlcAddr := network.BindPlcAddr
	if PlcAddr == "" {
		connectionlogger.Error(logger.ComponentService, "PLC not found.. Droping connection..")
		return
	}

	PlcConn, err := net.Dial("tcp", PlcAddr)
	if err != nil {
		connectionlogger.Error(logger.ComponentService, "Could not connect to PLC at %s with error: %d", PlcAddr, err)
		return
	}
	defer PlcConn.Close()
	connectionlogger.Info(logger.ComponentService, "Proxying data to PLC at %s", PlcAddr)

	go io.Copy(PlcConn, ClientConn)
	io.Copy(ClientConn, PlcConn)
}
