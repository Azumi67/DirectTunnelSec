package clientpackage

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"log"
	"net"
)

var aesKey []byte
var bSize int

func iranSide(conn net.Conn, localPort int, rmtAddr string, tcpNoDelay bool, encrypt bool, bSize int) {
	kharejConn, err := net.Dial("tcp", rmtAddr)
	if err != nil {
		log.Printf("Couldn't Connect to the remote server: %v", err)
		conn.Close()
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

	var copyErr error

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
			_, copyErr = io.Copy(encConn, conn)
			if copyErr != nil {
			}
			conn.Close() 
		}()

		_, copyErr = io.Copy(conn, decryptedConn)
		if copyErr != nil {
		}
		conn.Close() 
	} else {
		go func() {
			_, copyErr = io.Copy(kharejConn, conn)
			if copyErr != nil {
			}
			conn.Close()
		}()

		_, copyErr = io.Copy(conn, kharejConn)
		if copyErr != nil {
		}
		conn.Close() 
	}
}
