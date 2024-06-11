package main

import (
	"flag"
	"fmt"
	"log"
	"net"

)

var (
	aesKey     []byte
	bSize int
)

func main() {
	var localPort int
	var rmtAddr string
	var tcpNoDelay bool
	var encrypt bool
	var aesKeyString string

	flag.IntVar(&localPort, "local", 0, "Local port to listen on")
	flag.StringVar(&rmtAddr, "target", "", "Target address")
	flag.BoolVar(&tcpNoDelay, "noDelay", false, "Enable TCP no-delay")
	flag.BoolVar(&encrypt, "encrypt", false, "Enable AES encryption")
	flag.StringVar(&aesKeyString, "key", "", "AES key")
	flag.IntVar(&bSize, "buffer", 0, "TCP buffer size")
	flag.Parse()

	if localPort == 0 {
		log.Fatal("Local port is required")
	}

	if rmtAddr == "" {
		log.Fatal("Kharej address is required")
	}

	if encrypt && aesKeyString == "" {
		log.Fatal("AES authentication key is required when encryption is enabled. Use [openssl rand -hex 16]")
	}

	if encrypt {
		aesKey = []byte(aesKeyString)
		if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
			log.Fatal("AES key must be one of these lengths 16 or 24 or 32 bytes")
		}
	}

	listenAddr := fmt.Sprintf(":%d", localPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Couldn't listen to %s: %v", listenAddr, err)
	}

	log.Printf("Client is Listening on %s, forwarding to %s", listenAddr, rmtAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Couldn't accept connections from Kharej server: %v", err)
			continue
		}

		go clientpackage.iranSide(conn, localPort, rmtAddr, tcpNoDelay, encrypt, bSize)
	}
}
