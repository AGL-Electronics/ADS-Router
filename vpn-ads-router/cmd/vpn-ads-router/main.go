// authur TaireCat
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"vpn-ads-router/configs"
)

var BindPlcAddr string

/*
fingerprint is defined in the fingerprint.go file in the folder.
this file defines the ports the main.go file searches for aswell as subdomain.
change this and reguild the systemd package to run as service, or run in shell.

for running in the "go-proxy-netbird" folder run (assuming you in start folder)

	cd go-proxy-netbird
	nix-shell
	go run main.go fingerprint.go

this will run it in the shell (handy for testing)
a dummy command can be ran from win pc (assuming netbird connection to server) to test if the proxy connects

	ncat 100.77.95.203 {replace that with netbird ip} 48898
*/
func init() {
	log.Println("INIT: PLCPortFingerprint loaded with", len(configs.PlcFingerprint), "ports")
	log.Println("INIT: PLCSubnet is set to", configs.Subnet)
	log.Println("...")
}

func main() {
	ListenAddr := ":48898"

	// Initial PLC discovery
	BindPlcAddr = plcDiscover()
	if BindPlcAddr == "" {
		log.Fatalln("MAIN: Could not discover PLC at startup. Exiting.")
	}

	ln, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Fatalf("MAIN: Failed to listen on %s: %v", ListenAddr, err)
	}
	log.Printf("MAIN: Listening on %s for ADS connections...\n", ListenAddr)

	for {
		Conn, err := ln.Accept()
		if err != nil {
			log.Printf("MAIN: Connection accept error: %v", err)
			continue
		}

		go func(c net.Conn) {
			log.Println("MAIN: New connection from", c.RemoteAddr())

			if !validateBind(BindPlcAddr) {
				log.Println("MAIN: Cached PLC invalid, rescanning...")
				BindPlcAddr = plcDiscover()
			}

			if BindPlcAddr == "" {
				log.Println("MAIN: No valid PLC available, closing connection.")
				c.Close()
				return
			}

			Handleconnection(c)
		}(Conn)
	}
}

func Handleconnection(ClientConn net.Conn) {

	defer ClientConn.Close()
	log.Println("CON HANDLE: Incomming connection from", ClientConn.RemoteAddr())

	PlcAddr := BindPlcAddr
	if PlcAddr == "" {
		log.Println("CON HANDLE: PLC not found.. Droping connection..")
		return
	}

	PlcConn, err := net.Dial("tcp", PlcAddr)
	if err != nil {
		log.Printf("CON HANDLE: Could not connect to plc at %s: %v\n", PlcAddr, err)
		return
	}
	defer PlcConn.Close()
	log.Printf("CON HANDLE: Proxying data to PLC at %s\n", PlcAddr)

	go io.Copy(PlcConn, ClientConn)
	io.Copy(ClientConn, PlcConn)
}

func validateBind(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Printf("BIND VALID: Invalid cached address: %v", err)
		return false
	}

	timeout := 150 * time.Millisecond
	for _, P := range PlcFingerprint {
		target := fmt.Sprintf("%s:%d", host, P.port)
		Conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			if P.required {
				log.Printf("BIND VLAID: required port %d (%s) not open on %s", P.port, P.label, host)
				return false
			}
			log.Printf("BIND VALID: optional port %d (%s) not open on %s -continuing", P.port, P.label, host)
		} else {
			Conn.Close()
			log.Printf("BIND VALID: port %d (%s) is open on %s", P.port, P.label, host)
		}
	}
	return true
}

/*

this func is handeld by discorver now with fingerprint file as a "config"

func IsProbbalyMaybePlc(IP string) bool {
	Ports := []int{48898, 8016}
	open := 0
	for _, Port := range Ports {
		Conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", IP, Port), 200*time.Millisecond)
		if err == nil {
			open++
			Conn.Close()
		}
	}
	if open == len(Ports) {
		log.Printf("Found device at %s with ports 48898 and 8016 open, is-probbly-maybe a PLC", IP)
		return true
	}
	log.Printf("%s does not have bolth port 48898 and 8016 open (%d found)", IP, open)
}
*/
