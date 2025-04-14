package network

import (
	"net"

	"vpn-ads-router/pkg/logger"
)
var ProxynetID [6]byte

// GetLocalIP returns the local IP address of the machine.
func GetLocalIP(interfacename string) (net.IP, error) {
	iface, err := net.InterfaceByName(interfacename)
	if err != nil {
		logger.GlobalLogger.Error(logger.ComponentNetwork, "Error getting interface %s: %v", interfacename, err)
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
	logger.GlobalLogger.Error(logger.ComponentNetwork, "No valid IPv4 address found for interface %s", interfacename)
	return nil, err
}

//build netid from local ip and subnet mask
func BuildNetId(ip net.IP, suffix [2]byte) [6]byte {
	return [6]byte{ip[0], ip[1], ip[2], ip[3], suffix[0], suffix[1]}
	
}
