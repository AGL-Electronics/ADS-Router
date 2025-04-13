package network

import (
	"net"

	"vpn-ads-router/pkg/logger"
)
var localIPlogger = logger.GetLogger()

var ProxynetID [6]byte

func mainReadConfig() {
	var fingerprintConfig configs.FingerprintFile

	if err := configs.LoadJSONConfig("configs/PLC-Fingerprint.json", &fingerprintConfig); err != nil {
		systemlogger.Error(logger.ComponentNetwork, "Error loading fingerprint config: %v", err)
		return
	}

	if len(fingerprintConfig.Subnet) == 0 {
		systemlogger.Error(logger.ComponentNetwork, "No subnets found in fingerprint config")
		return
	}

	Subnet = fingerprintConfig.Subnet[0].Subnet //get the first subnet from the config, this should be changed to support multiple subnets in the future

}

// GetLocalIP returns the local IP address of the machine.
func getLocalIP(interfacename string) (net.IP, error) {
	iface, err := net.InterfaceByName(interfacename)
	if err != nil {
		localIPlogger.Error(logger.ComponentNetwork, "Error getting interface %s: %v", interfacename, err)
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err 
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil { // Check if the IP is IPv4
				return ipnet.IP, nil
			}
		}
	}
	localIPlogger.Error(logger.ComponentNetwork, "No valid IPv4 address found for interface %s", interfacename)
	return nil, err
}

//build netid from local ip and subnet mask
func buildNetId(ip net.IP, suffix [2]byte) [6]byte {
	return [6]byte{ip[0], ip[1], ip[2], ip[3], suffix[0], suffix[1]}
	
}
