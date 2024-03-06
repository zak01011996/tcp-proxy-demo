# TCP Proxy server (demo)

This is a simple TCP proxy example written in Golang.
Configuration file could be found here: [config/sample.yaml](/config/sample.yaml)

However this demo server could also run without any config file, default values will be used (default listen address is **localhost:7373**)

## Client example
```golang
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// Configuration...
const proxyAddress = "localhost:7373"

var cipherKey = []byte(`super_secret_key!`)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", proxyAddress)
	if err != nil {
		panic("Cannot resolve TCP address")
	}

	// Establish connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(fmt.Sprintf("Cannot establish connection: %-v", err))
	}
	defer conn.Close()

	// Prepare message to send
	writeMessage(conn)

	// Close write connection, we've sent everything needed
	conn.CloseWrite()

	// Read result
	reader := bufio.NewReader(conn)
	for {
		buff := make([]byte, 0xffff)
		n, err := reader.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("There were an issue reading results")
			}

			break
		}

		result := buff[:n]

		// Decrypt response
		encryptDecrypt(result)

		// Print results
		fmt.Println(string(result))
	}
}

func writeMessage(conn net.Conn) error {
	// As message we can set some HTTP data for example
	msg := []byte("GET / HTTP/1.0\r\n\r\n")

	// Encrypt message
	encryptDecrypt(msg)

	_, err := conn.Write(msg)
	return err
}

// for encryption and decryption of messages that will be sent to proxy
func encryptDecrypt(data []byte) {
	keyLen := len(cipherKey)
	for i := 0; i < len(data); i++ {
		data[i] = data[i] ^ cipherKey[i%keyLen]
	}
}
```