package router

import (
	"fmt"
	"net"
	"sync"

	"vpn-ads-router/pkg/logger"
)

type router struct {
	routingTable map[NetID]net.Conn	//netID -> client conn
	invokeMap   map[uint32]net.Conn	//invokeId -> original sender NetID

	mappingMutex sync.Mutex	
	invokeMutex  sync.Mutex

	plcConn net.Conn		//permanent connection to the PLC
	broacastIp net.IP		//derived from config
	ifaceName net.Interface	//derived from config
}


func NewRouter(broadcastIp net.IP, ifaceName string) *router {


}

