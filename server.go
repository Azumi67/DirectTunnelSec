package main

import (
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

)

var (
	aesKey     []byte
	bSize int
)

func main() {
	var listenPort int
	var lclAddr string
	var tcpNoDelay bool
	var encrypt bool
	var fobString string

	flag.IntVar(&listenPort, "listen", 0, "listening on this port")
	flag.StringVar(&lclAddr, "local", "", "local address")
	flag.BoolVar(&tcpNoDelay, "noDelay", false, "Enabling TCP nodelay")
	flag.BoolVar(&encrypt, "encrypt", false, "Enabling AES OFB encryption")
	flag.StringVar(&fobString, "key", "", "AES key")
	flag.IntVar(&bSize, "buffer", 0, "TCP buffer size")
	flag.Parse()

	if listenPort == 0 {
		log.Fatal("Listen port is required")
	}

	if lclAddr == "" {
		log.Fatal("Remote address is required")
	}

	if encrypt && fobString == "" {
		log.Fatal("AES authentication key is required when encryption is enabled. Use [openssl rand -hex 16]")
	}

	if encrypt {
		aesKey = []byte(fobString)
		if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
			log.Fatal("AES key must be one of these lengths 16 or 24 or 32 bytes")
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		log.Fatalf("Couldn't listen to: %v", err)
	}

	log.Printf("Kharej Server listening on port %d", listenPort)

	for {
		iranConn, err := listener.Accept()
		if err != nil {
			log.Printf("Couldn't accept connections from iran client server: %v", err)
			continue
		}
		if tcpNoDelay {
			if tcpConn, ok := iranConn.(*net.TCPConn); ok {
				err = tcpConn.SetNoDelay(true)
				if err != nil {
					log.Printf("Enabling TCP nodelay failed: %v", err)
				}
			}
		}

		go serverpackage.kharejSide(iranConn, lclAddr, tcpNoDelay, encrypt, bSize)
	}
}
