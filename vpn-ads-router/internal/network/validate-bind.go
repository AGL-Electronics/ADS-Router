package network

import (
	"fmt"
	"net"
	"time"

	"vpn-ads-router/configs"
	"vpn-ads-router/pkg/logger"
)

func validateBind(BindPlcAddr string) bool {
	logger := logger.GetLogger()
	host, _, err := net.SplitHostPort(BindPlcAddr)
	if err != nil {
		logger.Error(logger.ComponentNetwork, "Invalid cached address: %v", err)
		return false
	}

	timeout := 150 * time.Millisecond
	for _, P := range configs.PlcFingerprint {
		target := fmt.Sprintf("%s:%d", host, P.Port)
		Conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			if P.Required {
				logger.Error(logger.Componentnetwork, "required port %d (%s) not open on %s", P.Port, P.Label, host)
				return false
			}
			logger.Warn(logger.Componentnetwork, "optional port %d (%s) not open on %s -continuing", P.Port, P.Label, host)
		} else {
			Conn.Close()
			logger.Info(logger.Componentnetwork, "port %d (%s) is open on %s", P.Port, P.Label, host)
		}
	}
	return true
}
