package network

import (
	"fmt"
	"net"
	"time"

	"vpn-ads-router/pkg/config"
	"vpn-ads-router/pkg/logger"
)


func ValidateBind(BindPlcAddr string) bool {
	host, _, err := net.SplitHostPort(BindPlcAddr)
	if err != nil {
		logger.GlobalLogger.Error(logger.ComponentNetwork, "Invalid cached address: %v", err)
		return false
	}

	timeout := 150 * time.Millisecond
	for _, P := range config.AppConfig.Fingerprint.PlcFingerprint {
		target := fmt.Sprintf("%s: %d", host, P.Port) //does not work with IPv6
		Conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			if P.Required {
				logger.GlobalLogger.Error(logger.ComponentNetwork, "required port %d (%s) not open on %s", P.Port, P.Label, host)
				return false
			}
			logger.GlobalLogger.Warn(logger.ComponentNetwork, "optional port %d (%s) not open on %s -continuing", P.Port, P.Label, host)
		} else {
			Conn.Close()
			logger.GlobalLogger.Info(logger.ComponentNetwork, "port %d (%s) is open on %s", P.Port, P.Label, host)
		}
	}
	return true
}
