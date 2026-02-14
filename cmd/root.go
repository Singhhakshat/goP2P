package cmd

import (
	"encoding/json"
	"fmt"          //For printing to the console (think console.log or System.out.println).
	"goP2P/crypto" //The internal package we created earlier
	"goP2P/p2p"    // internal package
	"io"
	"net/http"
	"os" //For interacting with the operating system, such as exiting the program.

	"github.com/spf13/cobra" //For creating command-line interfaces (CLIs) in Go.
)

// The address of our Matchmaker server
const matchmakerURL = "http://localhost:8080"

// PeerData represents the JSON structure we expect from the Matchmaker
type PeerData struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

var rootCommand = &cobra.Command{
	Use:     "p2p-share",
	Aliases: []string{"goP2P"},
	Short:   "A secure P2P file sharing tool",
}

var sendCommand = &cobra.Command{
	Use:     "send [filename]",
	Aliases: []string{"send-file"},
	Short:   "send file to the peer",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		code := crypto.GenerateCode()
		fmt.Printf("Your secret hash for connection : %s\n", code)

		// 1. getting public IP for sender
		fmt.Println(">> Getting public IP <<")
		resp, IPErr := http.Get("https://api.ipify.org")
		if IPErr != nil {
			fmt.Printf("Error getting public IP")
			return
		}
		defer resp.Body.Close()

		ip, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			fmt.Println("Error parsing IP response")
		}
		publicIP := string(ip)
		fmt.Printf("Public IP detected : %s\n", publicIP)

		// 2. register this hash with the matchmaker
		registerURL := matchmakerURL + "/register?hash=" + code + "&ip=" + publicIP + "&port=3000"
		resp, err := http.Get(registerURL)
		if err != nil && resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Could not reach the matchmaker server")
			return
		}
		defer resp.Body.Close()

		fmt.Println("Registered with matchmaker successfully")

		// 3. starting our p2p server
		error := p2p.StartServer("3000", fileName)
		if error != nil {
			fmt.Println("Error starting server", err)
		}
	},
}

var recieveCommand = &cobra.Command{
	Use:     "recieve [hash]",
	Aliases: []string{"get"},
	Short:   "Recieve file using secret hash",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := args[0]
		fmt.Printf("looking up sender using hash %s\n", hash)

		// 1. HTTP GET: Ask the Matchmaker for the IP associated with this hash
		lookupURL := matchmakerURL + "/lookup?hash=" + hash
		resp, lookupErr := http.Get(lookupURL)
		if lookupErr != nil && resp.StatusCode != http.StatusOK {
			fmt.Println("Hash not found of matchmaker is offline")
			return
		}
		defer resp.Body.Close()

		// 2. decode the JSON response body to get the sender IP and port
		var peer PeerData
		decoderError := json.NewDecoder(resp.Body).Decode(&peer)
		if decoderError != nil {
			fmt.Println("Error parsing matchmaker response")
			return
		}

		// 3. starting p2p server
		outputFileName := "downloaded_" + hash + ".txt"
		err := p2p.ReceiveFile(peer.IP, peer.Port, outputFileName)
		if err != nil {
			fmt.Printf("Error occured while recieving the file %v\n", err)
		}
	},
}

func Execute() {
	rootCommand.AddCommand(sendCommand)
	rootCommand.AddCommand(recieveCommand)
	err := rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
