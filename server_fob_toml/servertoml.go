package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"flag"
)

var aesKey []byte

func kharejSide(iranConn net.Conn, lclAddr string, tcpNoDelay bool, encrypt bool, bSize int) {
	kharejConn, err := net.Dial("tcp", lclAddr)
	if err != nil {
		logrus.Errorf("Couldn't establish the tunnel: %v", err)
		iranConn.Close()
		return
	}
	defer kharejConn.Close()

	if tcpNoDelay {
		if tcpConn, ok := kharejConn.(*net.TCPConn); ok {
			err = tcpConn.SetNoDelay(true)
			if err != nil {
				logrus.Errorf("Enabling TCP nodelay failed: %v", err)
			}
		}
	}

	if bSize > 0 {
		if tcpConn, ok := kharejConn.(*net.TCPConn); ok {
			err = tcpConn.SetReadBuffer(bSize)
			if err != nil {
				logrus.Errorf("Setting up TCP read buffer failed: %v", err)
			}
			err = tcpConn.SetWriteBuffer(bSize)
			if err != nil {
				logrus.Errorf("Setting up TCP write buffer failed: %v", err)
			}
		}
	}

	if encrypt {
		block, err := aes.NewCipher(aesKey)
		if err != nil {
			logrus.Errorf("Creation of AES ciphersec failed: %v", err)
			return
		}
		stream := cipher.NewOFB(block, make([]byte, aes.BlockSize))
		encConn := &cipher.StreamWriter{S: stream, W: kharejConn}
		decryptedConn := &cipher.StreamReader{S: stream, R: kharejConn}

		go func() {
			_, err := io.Copy(encConn, iranConn)
			if err != nil {
				
			}
			iranConn.Close()
		}()

		_, err = io.Copy(iranConn, decryptedConn)
		if err != nil {
			
		}
		iranConn.Close()
	} else {
		go func() {
			_, err := io.Copy(kharejConn, iranConn)
			if err != nil {
				
			}
			iranConn.Close()
		}()

		_, err = io.Copy(iranConn, kharejConn)
		if err != nil {
			
		}
		iranConn.Close()
	}
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "TOML config filepath needed")
	flag.Parse()

	if configFile == "" {
		logrus.Fatal("Config filepath is required")
	}

	config, err := toml.LoadFile(configFile)
	if err != nil {
		logrus.Fatalf("error loading %s: %v", configFile, err)
	}

	listenPort := config.Get("listen_port").(int64)
	lclAddr := config.Get("local_address").(string)
	tcpNoDelay := config.Get("tcp_no_delay").(bool)
	encrypt := config.Get("encrypt").(bool)
	aesKeyString, _ := config.Get("key").(string) 
	bSize := int(config.Get("buffer_size").(int64))

	if encrypt {
		aesKey = []byte(aesKeyString)
		if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
			logrus.Fatal("AES key must be 16, 24, or 32 bytes in length")
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		logrus.Fatalf("Couldn't listen to: %v", err)
	}

	defer listener.Close()
	logrus.Infof("Kharej Server listening on port %d", listenPort)

	for {
		iranConn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("Couldn't accept connections from iran client server: %v", err)
			continue
		}

		go kharejSide(iranConn, lclAddr, tcpNoDelay, encrypt, bSize)
	}
}
