package network

import (
	"fmt"
	"log"
	"net"
	"time"

	"vpn-ads-router/configs"
)

//this fine handels the the scanning of the network and discovering the plc

var BindPlcAddr string //cached plc address, used to avoid scanning the network every time
var Subnet string      //subnet to scan, set in the config file

func plcDiscover() string { //check common beckhoff ports to identify the plc based on open ports, this to filter out false positives, if false pos still occur add more ports to fingerprint.
	if BindPlcAddr != "" {
		if validateBind(BindPlcAddr) {
			log.Println("PLC DISC: Using cached PLC Address:", BindPlcAddr)
			return BindPlcAddr
		}
	}

	Timeout := 150 * time.Millisecond //timeout for conn

	log.Println("PLC DISC: Scanning for PLC...")
	log.Println("PLC DISC: Atempting to identify plc with port fingerprint, Can be changed in fingerprint file")

	for i := 1; i <= 254; i++ {
		BaseIp := fmt.Sprintf("%s%d", Subnet, i)
		matched := true
		for _, P := range configs.PlcFingerprint {
			Addr := fmt.Sprintf("%s:%d", BaseIp, P.port) //does not work with IPv6
			Conn, err := net.DialTimeout("tcp", Addr, Timeout)
			if err == nil {
				log.Printf("PLC DISC: Port %d (%s) open on %s", P.port, P.label, BaseIp)
				Conn.Close()
			} else {
				if P.Required {
					log.Printf("PLC DISC: Required port %d (%s) not open on %s", P.port, P.label, BaseIp)
					matched = false
					break
				} else {
					log.Printf("PLC DISC: Optional (not required) port %d (%s) not open at %s", P.port, P.label, BaseIp)
				}
			}

		}
		if matched {
			log.Printf("PLC DISC: Found device matching fingerpirint at %s, Likely a PLC, Caching ip", BaseIp)
			BindPlcAddr = fmt.Sprintf("%s:48898", BaseIp)
			return BindPlcAddr
		}
	}
	log.Printf("PLC DISC: No PLC found with fingerprint on subnet %s\n", Subnet)
	return ""
}
