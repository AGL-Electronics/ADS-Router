package proxy

import (
	"net"
	"time"

	"vpn-ads-router/internal/network"
	"vpn-ads-router/internal/router"
	"vpn-ads-router/pkg/logger"
)

var schedulerlogger = logger.GetLogger()

var PlcAddr = network.BindPlcAddr        //cached plc address, used to avoid scanning the network every time
var PlcConn net.Conn                     //active PLC connection
var RequestMap = make(map[uint32]string) //map of request id to source net id, used to identify the source of the request

// StartScheduler starts the scheduler, which listens for incoming connections and lissens for incoming messages and sends them to the router.
func StartScheduler() {
	var err error
	timeout := 150 * time.Millisecond
	PlcConn, err = net.DialTimeout("tcp", PlcAddr, timeout)
	if err != nil {
		schedulerlogger.Error(logger.ComponentProxy, "Error connecting to PLC at %s: %v", PlcAddr, err)
	}
	schedulerlogger.Info(logger.ComponentProxy, "Connected to PLC at %s", PlcAddr)

	for msg := range IncommingConnChan {
		err := router.ProcessMsg(msg, PlcConn)
		if err != nil {
			schedulerlogger.Error(logger.ComponentProxy, "Error processing message: %v", err)
			schedulerlogger.Error(logger.ComponentProxy, "Error processing message: %v", msg.payload)
			continue
		}
	}
}
