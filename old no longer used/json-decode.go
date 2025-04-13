package configs

import (
	"encoding/json"
	"os"

	"vpn-ads-router/pkg/logger"
)

//-- config file names --
// The config files are stored in the configs folder, the names are hardcoded here
// and should not be changed. The files are loaded at startup and the data is stored in the config structs.
// The config files are in JSON format and should be valid JSON. The files are loaded using the LoadJSONConfig function.

var PlcConfigFile = "configs/Proxy-Config.json"
var PlcFingerprintFile = "configs/PLC-Fingerprint.json"

// --fingerprint.json--
type FingerprintFile struct {
	PlcFingerprint []PlcFingerprint `json:"fingerprint"` // list of fingerprints
	Subnet         []Subnet         `json:"subnet"`      // list of subnets to scan
}

type PlcFingerprint struct {
	Port     int    `json:"port"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
}
type Subnet struct {
	Subnet string `json:"subnet"`
}

// --Proxy-Config.json--
type ProxyConfigFile struct {
	PlcCredential []PlcCredential `json:"plcCredential"` // list of plc credentials
	ProxyConfig   []ProxyConfig   `json:"proxyConfig"`	// list of proxy configs
}

type PlcCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ProxyConfig struct {
	StaticNetid       string `json:"staticNetidSuffix"`
	EthernetInterface string `json:"ethernetInterface"`
}

// --unpacking logic--
var configlogger = logger.GetLogger()

func LoadJSONConfig(path string, out interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		configlogger.Error(logger.ComponentGeneral, "Error reading config file %s: %v", path, err)
		return err
	}

	if err := json.Unmarshal(content, out); err != nil {
		configlogger.Error(logger.ComponentGeneral, "Error parsing config file %s: %v", path, err)
		return err
	}

	return nil

}
