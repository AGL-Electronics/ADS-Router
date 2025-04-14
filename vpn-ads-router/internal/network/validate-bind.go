package network

import (
	"fmt"
	"net"
	"time"

	"vpn-ads-router/configs"
	"vpn-ads-router/pkg/logger"
)

var bindlogger = logger.GetLogger()

func ValidateBind(BindPlcAddr string) bool {
	host, _, err := net.SplitHostPort(BindPlcAddr)
	if err != nil {
		bindlogger.Error(logger.ComponentNetwork, "Invalid cached address: %v", err)
		return false
	}

	timeout := 150 * time.Millisecond
	for _, P := range configs.PlcFingerprint {
		target := fmt.Sprintf("%s: %d", host, P.Port)
		Conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			if P.Required {
				bindlogger.Error(logger.ComponentNetwork, "required port %d (%s) not open on %s", P.Port, P.Label, host)
				return false
			}
			bindlogger.Warn(logger.ComponentNetwork, "optional port %d (%s) not open on %s -continuing", P.Port, P.Label, host)
		} else {
			Conn.Close()
			bindlogger.Info(logger.ComponentNetwork, "port %d (%s) is open on %s", P.Port, P.Label, host)
		}
	}
	return true
}
