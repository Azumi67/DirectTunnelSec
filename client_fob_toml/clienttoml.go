package clienttoml

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

func iranSide(conn net.Conn, localPort int, rmtAddr string, tcpNoDelay bool, encrypt bool, bSize int) {
	kharejConn, err := net.Dial("tcp", rmtAddr)
	if err != nil {
		logrus.Errorf("Couldn't connect to the remote server: %v", err)
		conn.Close()
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
			_, err := io.Copy(encConn, conn)
			if err != nil {
				logrus.Errorf("Encryption failed: %v", err)
			}
			conn.Close()
		}()

		_, err = io.Copy(conn, decryptedConn)
		if err != nil {
			logrus.Errorf("Decryption failed: %v", err)
		}
		conn.Close()

	} else {
		go func() {
			_, err := io.Copy(kharejConn, conn)
			if err != nil {
				logrus.Errorf("Data transfer failed: %v", err)
			}
			conn.Close()
		}()

		_, err = io.Copy(conn, kharejConn)
		if err != nil {
			logrus.Errorf("Data transfer failed: %v", err)
		}
		conn.Close()
	}
}

func iranClient(configFile string) {
	config, err := toml.LoadFile(configFile)
	if err != nil {
		logrus.Fatalf("Error loading %s: %v", configFile, err)
	}

	localPort := config.Get("local_port").(int64)
	rmtAddr := config.Get("target_address").(string)
	tcpNoDelay := config.Get("tcp_no_delay").(bool)
	encrypt := config.Get("encrypt").(bool)
	aesKeyString := config.Get("key").(string)
	bSize := int(config.Get("buffer_size").(int64))

	aesKey = []byte(aesKeyString)
	if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
		logrus.Fatal("AES key must be 16, 24, or 32 bytes in length")
	}

	listenAddr := fmt.Sprintf(":%d", localPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Fatalf("Couldn't listen to %s: %v", listenAddr, err)
	}

	defer listener.Close()
	logrus.Infof("IRAN is listening on %s, forwarding to %s", listenAddr, rmtAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("Couldn't accept connections: %v", err)
			continue
		}

		go iranSide(conn, int(localPort), rmtAddr, tcpNoDelay, encrypt, bSize)
	}
}
