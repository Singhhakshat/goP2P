package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
)

// Global state to store hashes and IPs.
var (
	senders = make(map[string]net.Conn)
	mutex   = &sync.Mutex{}
)

func main() {
	// 1. start a HTTP server for Hash registration
	go func() {
		http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
			hash := request.URL.Query().Get("hash")

			fmt.Printf("registered hash %s\n", hash)
			writer.WriteHeader(http.StatusOK)
		})

		fmt.Println("Matchmaker HTTP running on port 8080")
		http.ListenAndServe(":8080", nil)
	}()

	// 2. start the TCP relay server
	fmt.Println(">> matchmaker relay working on :9000 <<")
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("TCP relay failed:", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleRelayConnection(conn)
	}
}

func handleRelayConnection(conn net.Conn) {
	buf := make([]byte, 11)
	_, error := conn.Read(buf)
	if error != nil {
		conn.Close()
		return
	}

	message := string(buf)
	role := message[:3]
	hash := message[4:]

	mutex.Lock()
	if role == "SND" {
		fmt.Printf("Sender connected for hash %s\n", hash)
		senders[hash] = conn
		mutex.Unlock()
	}

	if role == "RCV" {
		senderConn, exists := senders[hash]
		if exists {
			delete(senders, hash)
			mutex.Unlock()

			fmt.Printf("Stitching connections for hash %s\n", hash)
			//stitching connections together
			go func() {
				io.Copy(conn, senderConn)
				conn.Close()
				senderConn.Close()
			}()
		} else {
			fmt.Println("Reciever connected but no sender found")
			conn.Close()
			mutex.Unlock()
		}
	}
}
