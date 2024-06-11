package serverpackage

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"log"
	"net"
	"sync"
)

var aesKey []byte
var bSize int

func kharejSide(iranConn net.Conn, lclAddr string, tcpNoDelay bool, encrypt bool, bSize int) {
	kharejConn, err := net.Dial("tcp", lclAddr)
	if err != nil {
		log.Printf("Couldn't Establish the tunnel: %v", err)
		iranConn.Close()
		return
	}
	defer kharejConn.Close()

	if tcpNoDelay {
		if tcpConn, ok := kharejConn.(*net.TCPConn); ok {
			err = tcpConn.SetNoDelay(true)
			if err != nil {
				log.Printf("Enabling TCP nodelay failed: %v", err)
			}
		}
	}

	if bSize > 0 {
		if tcpConn, ok := kharejConn.(*net.TCPConn); ok {
			err = tcpConn.SetReadBuffer(bSize)
			if err != nil {
				log.Printf("Setting up TCP read buffer failed: %v", err)
			}
			err = tcpConn.SetWriteBuffer(bSize)
			if err != nil {
				log.Printf("Setting up TCP write buffer failed: %v", err)
			}
		}
	}

	if encrypt {
		block, err := aes.NewCipher(aesKey)
		if err != nil {
			log.Printf("Creation of AES ciphersec failed: %v", err)
			return
		}
		stream := cipher.NewOFB(block, make([]byte, aes.BlockSize))
		encConn := &cipher.StreamWriter{S: stream, W: kharejConn}

		stream = cipher.NewOFB(block, make([]byte, aes.BlockSize))
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
