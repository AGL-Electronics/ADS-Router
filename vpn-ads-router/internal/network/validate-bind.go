package network

import (
	"fmt"
	"net"
	"time"

	"vpn-ads-router/pkg/logger"
)

func validateBind(BindPlcAddr string) bool {
	logger := logger.GetLogger()
	host, _, err := net.SplitHostPort(BindPlcAddr)
	if err != nil {
		logger.Error(logger.Componentnetwork, "Invalid cached address: %v", err)
		return false
	}

	timeout := 150 * time.Millisecond
	for _, P := range PlcFingerprint {
		target := fmt.Sprintf("%s:%d", host, P.port)
		Conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			if P.required {
				logger.Error(logger.Componentnetwork, "required port %d (%s) not open on %s", P.port, P.label, host)
				return false
			}
			logger.Warn(logger.Componentnetwork, "optional port %d (%s) not open on %s -continuing", P.port, P.label, host)
		} else {
			Conn.Close()
			logger.Info(logger.Componentnetwork, "port %d (%s) is open on %s", P.port, P.label, host)
		}
	}
	return true
}
