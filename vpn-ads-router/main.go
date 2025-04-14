package main

import (
	"net"

	"vpn-ads-router/internal/network"
	"vpn-ads-router/pkg/config"

	// "vpn-ads-router/internal/proxy" needs to be added back later
	"vpn-ads-router/pkg/logger"
)

var BindPlcAddr string
var systemlogger = logger.GetLogger()

func init() {
	// Initialize the logger
	logger.InitGlobalLogger("logs", logger.LogLevelInfo, []logger.Component{
		logger.ComponentRouter,
		logger.ComponentProxy,
		logger.ComponentNetwork,
		logger.ComponentADS,
		logger.ComponentVPN,
		logger.ComponentGeneral,
		logger.ComponentService,
	})

	systemlogger.Info(logger.ComponentService, "INIT: PLC Port Fingerprint loaded with %d ports", len(config.AppConfig.Fingerprint.PlcFingerprint))
	systemlogger.Info(logger.ComponentService, "INIT: PLC Subnet is set to %s", config.AppConfig.Fingerprint.Subnets[0])
	systemlogger.Info(logger.ComponentService, "INIT: Starting VPN-ADS Router...")

}

func main() {
	ListenAddr := ":48898"

	// Initial PLC discovery
	BindPlcAddr = network.PlcDiscover()
	if BindPlcAddr == "" {
		systemlogger.Error(logger.ComponentService, "Could not discover PLC at startup. Exiting.")
	}

	ln, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		systemlogger.Error(logger.ComponentService, "Failed to listen on %s: %v", ListenAddr, err)
	}

	systemlogger.Error(logger.ComponentService, "Listening on %s for ADS connections...", ListenAddr)

	for {
		Conn, err := ln.Accept()
		if err != nil {
			systemlogger.Error(logger.ComponentService, "MAIN: Connection accept error: %v", err)
			continue
		}

		go func(c net.Conn) {

			systemlogger.Info(logger.ComponentService, "MAIN: New connection from %s", c.RemoteAddr())
			systemlogger.Info(logger.ComponentService, "MAIN: Current PLC Address: %s", BindPlcAddr)

			if !network.ValidateBind(BindPlcAddr) {
				systemlogger.Info(logger.ComponentService, "MAIN: Cached PLC invalid, rescanning...")
				BindPlcAddr = network.PlcDiscover()
			}

			if BindPlcAddr == "" {
				systemlogger.Fatal(logger.ComponentService, "MAIN: No valid PLC available, closing connection.")
				c.Close()
				return
			}
			// add handeling for connection back here

		}(Conn)
	}
}
