package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"bufio"
)

func ResponseHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"result" : "OK"}`)
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/test", ResponseHandler)

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("../certs/ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("../certs/server.pem", "../certs/server-key.pem")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},        // server certificate which is validated by the client
		ClientCAs:    caCertPool,                     // used to verify the client cert is signed by the CA and is therefore valid
		ClientAuth:   tls.RequireAndVerifyClientCert, // this requires a valid client certificate to be supplied during handshake
	}
	ln, err := tls.Listen("tcp", "localhost:9443", tlsConfig)
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}

	log.Println("listen: ", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("failed to accept conn: %s", err)
		}

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			log.Fatalf("failed to cast conn to tls.Conn")
		}

		go handleTLSConnection(tlsConn)
	}
}

func handleTLSConnection(tlsConn *tls.Conn) {

	tag := fmt.Sprintf("[%s -> %s]", tlsConn.LocalAddr(), tlsConn.RemoteAddr())
	log.Printf("%s accept", tag)

	defer tlsConn.Close()

	// this is required to complete the handshake and populate the connection state
	// we are doing this so we can print the peer certificates prior to reading / writing to the connection
	err := tlsConn.Handshake()
	if err != nil {
		log.Fatalf("failed to complete handshake: %s", err)
	}

	if len(tlsConn.ConnectionState().PeerCertificates) > 0 {
		log.Printf("%s client common name: %+v", tag, tlsConn.ConnectionState().PeerCertificates[0].Subject.CommonName)
	}

	b := bufio.NewReader(tlsConn)

	line, err := b.ReadBytes('\n')
	if err != nil {
		log.Fatalf("%s failed to read line from conn: %s", tag, err)
	}

	log.Printf("%s line: %s", tag, line)

	tlsConn.Write(line)
}