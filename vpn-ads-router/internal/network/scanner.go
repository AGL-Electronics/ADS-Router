package network

import (
	"fmt"
	"net"
	"time"

	"vpn-ads-router/pkg/config"
	"vpn-ads-router/pkg/logger"
)

//this fine handels the the scanning of the network and discovering the plc
//it uses the plc fingerprint to identify the plc based on open ports

var BindPlcAddr string                 //cached plc address, used to avoid scanning the network every time
var scannerSubnet string               //subnet to scan, set in the config file

func PlcDiscover() string { //check common beckhoff ports to identify the plc based on open ports, this to filter out false positives, if false pos still occur add more ports to fingerprint.
	if BindPlcAddr != "" {
		if ValidateBind(BindPlcAddr) {
			logger.GlobalLogger.Info(logger.ComponentNetwork, "PLC DISC: Using cached PLC Address: %s", BindPlcAddr)
			return BindPlcAddr
		}
	}

	Timeout := 150 * time.Millisecond //timeout for conn
	logger.GlobalLogger.Info(logger.ComponentNetwork, "PLC DISC: Scanning for PLC...")
	logger.GlobalLogger.Info(logger.ComponentNetwork, "PLC DISC: Attempting to identify PLC with port fingerprint, Can be changed in fingerprint file")

	for _, subnet := range config.AppConfig.Fingerprint.Subnets {
		for i := 1; i <= 254; i++ {
			BaseIp := fmt.Sprintf("%s%d", subnet, i)
			matched := true

			for _, P := range config.AppConfig.Fingerprint.PlcFingerprint {
				Addr := fmt.Sprintf("%s: %d", BaseIp, P.Port) //does not work with IPv6
				Conn, err := net.DialTimeout("tcp", Addr, Timeout)

				if err == nil {
					logger.GlobalLogger.Info(logger.ComponentNetwork, "PLC DISC: Port %d (%s) open on %s", P.Port, P.Label, BaseIp)
					Conn.Close()
				} else {
					if P.Required {
						logger.GlobalLogger.Error(logger.ComponentNetwork, "PLC DISC: Required port %d (%s) not open on %s", P.Port, P.Label, BaseIp)
						matched = false
						break
					} else {
						logger.GlobalLogger.Warn(logger.ComponentNetwork, "PLC DISC: Optional (not required) port %d (%s) not open at %s", P.Port, P.Label, BaseIp)
					}
				}
			}
		if matched {
			logger.GlobalLogger.Info(logger.ComponentNetwork, "PLC DISC: Found device matching fingerprint at %s, likely a PLC, caching IP", BaseIp)
			BindPlcAddr = fmt.Sprintf("%s:48898", BaseIp)
			return BindPlcAddr
		}
	}
	logger.GlobalLogger.Fatal(logger.ComponentNetwork, "PLC DISC: No PLC found with fingerprint on subnet %s\n", Subnet)
	//os.Exit(1) 	//might need to live at main.go, test this with real pc on network
	}
	return ""
}
