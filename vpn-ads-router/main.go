package main

import (
	"net"
	"os"

	"vpn-ads-router/internal/network"
	"vpn-ads-router/internal/proxy"
	"vpn-ads-router/pkg/config"
	"vpn-ads-router/pkg/logger"
)

var BindPlcAddr string

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

	logger.GlobalLogger.Info(logger.ComponentService, "INIT: PLC Port Fingerprint loaded with %d ports", len(config.AppConfig.Fingerprint.PlcFingerprint))
	logger.GlobalLogger.Info(logger.ComponentService, "INIT: PLC Subnet is set to %s", config.AppConfig.Fingerprint.Subnets)
	logger.GlobalLogger.Info(logger.ComponentService, "INIT: Starting VPN-ADS Router...")

}

func main() {
	ListenAddr := ":48898"

	// Initial PLC discovery
	BindPlcAddr = network.PlcDiscover()
	if BindPlcAddr == "" {
		logger.GlobalLogger.Fatal(logger.ComponentService, "Could not discover PLC at startup. Exiting.")
		os.Exit(1) //might need to live here instead of scanner.go, test this with real pc on network
	}

	ln, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		logger.GlobalLogger.Error(logger.ComponentService, "Failed to listen on %s: %v", ListenAddr, err)
	}

	logger.GlobalLogger.Error(logger.ComponentService, "Listening on %s for ADS connections...", ListenAddr)

	for {
		Conn, err := ln.Accept()
		if err != nil {
			logger.GlobalLogger.Error(logger.ComponentService, "MAIN: Connection accept error: %v", err)
			continue
		}

		go func(c net.Conn) {

			logger.GlobalLogger.Info(logger.ComponentService, "MAIN: New connection from %s", c.RemoteAddr())
			logger.GlobalLogger.Info(logger.ComponentService, "MAIN: Current PLC Address: %s", BindPlcAddr)

			if !network.ValidateBind(BindPlcAddr) {
				logger.GlobalLogger.Info(logger.ComponentService, "MAIN: Cached PLC invalid, rescanning...")
				BindPlcAddr = network.PlcDiscover()
			}

			if BindPlcAddr == "" {
				logger.GlobalLogger.Fatal(logger.ComponentService, "MAIN: No valid PLC available, closing connection.")
				c.Close()
				return
			}

			proxy.HandleClient(c)
		}(Conn)
	}
}
