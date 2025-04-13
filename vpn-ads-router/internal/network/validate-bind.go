package network

import (
	"fmt"
	"net"
	"time"

	"vpn-ads-router/pkg/config"
	"vpn-ads-router/pkg/logger"
)

var bindlogger = logger.GetLogger()

var Subnet string
var PlcFingerprint []configs.PlcFingerprint


func validateBindreadConfig() {
	var fingerprintConfig configs.FingerprintFile

	if err := configs.LoadJSONConfig("configs/PLC-Fingerprint.json", &fingerprintConfig); err != nil {
		scannerlogger.Error(logger.ComponentNetwork, "Error loading fingerprint config: %v", err)
		return
	}

	if len(fingerprintConfig.Subnet) == 0 {
		scannerlogger.Error(logger.ComponentNetwork, "No subnets found in fingerprint config")
		return
	}

	Subnet = fingerprintConfig.Subnet[0].Subnet //get the first subnet from the config, this should be changed to support multiple subnets in the future

	PlcFingerprint = fingerprintConfig.PlcFingerprint //get the plc fingerprint from the config
}

func ValidateBind(BindPlcAddr string) bool {
	host, _, err := net.SplitHostPort(BindPlcAddr)
	if err != nil {
		bindlogger.Error(logger.ComponentNetwork, "Invalid cached address: %v", err)
		return false
	}

	timeout := 150 * time.Millisecond
	for _, P := range PlcFingerprint {
		target := fmt.Sprintf("%s: %d", host, P.Port) //does not work with IPv6
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
