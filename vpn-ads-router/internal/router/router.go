package router

import (
	"vpn-ads-router/pkg/logger"
)


func PatchSourceNetId(data []byte, netId [6]byte) {
	if len(data) < 16 {
		return //invalid data
	}

	copy(data[10:16], netId[:]) //copy the new net id to the payload
	logger.GlobalLogger.Debug(logger.ComponentRouter, "Patched source NetID in payload: %v", data[10:16])
}

func ExtractInvokeId(data []byte) uint32 {
	if len(data) < 40 {
		return 0 //invalid data
	}

	return uint32(data[36]) |
		(uint32(data[37]) << 8) |
		(uint32(data[38]) << 16) |
		(uint32(data[39]) << 24) //extract the invoke id from the payload, this is a 4 byte value
}
