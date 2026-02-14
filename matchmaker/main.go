package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type peerData struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// Global state to store hashes and IPs.
var (
	peers = make(map[string]peerData)
	mutex = &sync.Mutex{}
)

func main() {
	// register the API endpoints
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/lookup", lookupHandler)

	fmt.Println("Matchmaker server running on port 8080")

	// start the http server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failure: ", err)
	}
}

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	hash := request.URL.Query().Get("hash")
	ip := request.URL.Query().Get("ip")
	port := request.URL.Query().Get("port")

	if hash == "" || ip == "" || port == "" {
		http.Error(writer, "Missing parameters", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	peers[hash] = peerData{IP: ip, Port: port}
	mutex.Unlock()

	fmt.Printf("Registered hash %s -> %s:%s\n", hash, ip, port)
	writer.WriteHeader(http.StatusOK)
}

func lookupHandler(writer http.ResponseWriter, request *http.Request) {
	hash := request.URL.Query().Get("hash")

	mutex.Lock()
	peer, exists := peers[hash]

	if !exists {
		http.Error(writer, "Data not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(peer)
}
