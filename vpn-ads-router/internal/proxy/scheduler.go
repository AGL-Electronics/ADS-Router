package proxy

import (
	"net"
	"time"

	"vpn-ads-router/configs"
	"vpn-ads-router/pkg/logger"
	"vpn-ads-router/internal/network"
)


var schedulerlogger = logger.GetLogger()

var PlcAddr = network.BindPlcAddr 			//cached plc address, used to avoid scanning the network every time
var PlcConn net.Conn 						//active PLC connection
var RequestMap = make(map[uint32]string)	//map of request id to source net id, used to identify the source of the request

func validateBindreadConfig() {
	var fingerprintConfig configs.FingerprintFile

	if err := configs.LoadJSONConfig("configs/PLC-Fingerprint.json", &fingerprintConfig); err != nil {
		schedulerlogger.Error(logger.ComponentProxy, "Error loading fingerprint config: %v", err)
		return
	}

	if len(fingerprintConfig.Subnet) == 0 {
		schedulerlogger.Error(logger.ComponentProxy, "No subnets found in fingerprint config")
		return
	}

	Subnet = fingerprintConfig.Subnet[0].Subnet //get the first subnet from the config, this should be changed to support multiple subnets in the future

	PlcFingerprint = fingerprintConfig.PlcFingerprint //get the plc fingerprint from the config
}


func StartScheduler() {
    // 1. Connect to the PLC
    // 2. Loop over IncomingChan
    // 3. Patch source NetID
    // 4. Track invokeID → SourceNetID
    // 5. Forward to PLC

	var err error
	PlcConn, err = net.DialTimeout("tcp", PlcAddr)
	if err != nil {
		schedulerlogger.Error(logger.ComponentProxy, "Error connecting to PLC: %v", err)
	}
	schedulerlogger.Info(logger.ComponentProxy, "Connected to PLC at %s", PlcAddr)
	
	go StartResponesHandler(PlcConn)

	for msg := range IncommingConnChan {
		msg.SourceNetId = ParseSourceNetId(msg.payload) //parse the source net id from the payload
		msg.payload = msg.payload[16:] //remove the source net id from the payload


	}
}
