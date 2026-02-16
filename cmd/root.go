package cmd

import (
	"fmt"          //For printing to the console (think console.log or System.out.println).
	"goP2P/crypto" //The internal package we created earlier

	// internal package
	"io"
	"net"
	"net/http"
	"os" //For interacting with the operating system, such as exiting the program.

	"github.com/spf13/cobra" //For creating command-line interfaces (CLIs) in Go.
)

// The address of our Matchmaker server
const matchmakerURL = "gop2p-production.up.railway.app"
const matchmakerTCP = "shortline.proxy.rlwy.net:56002"

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

		// 1. Check if file exists
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("Error: Could not open file '%s'\n", fileName)
			return
		}
		defer file.Close()

		// 2. Generate Hash
		code := crypto.GenerateCode()
		fmt.Printf("Your secret hash for connection : %s\n", code)

		// 1. Register via HTTP
		registerURL := matchmakerURL + "/register?hash=" + code
		http.Get(registerURL)

		// 3. Connect to the relay via TCP
		conn, err := net.Dial("tcp", matchmakerTCP)
		if err != nil {
			fmt.Println("Could not connect to relay: ", err)
			return
		}
		defer conn.Close()

		// 3. Tell the matchmaker we are the sender
		nametag := fmt.Sprintf("SND:%-6s ", code)
		conn.Write([]byte(nametag[:11]))

		// 4. stream the file to the relay
		fmt.Println("Sending file")
		io.Copy(conn, file)

		fmt.Println("File sent successfully")

	},
}

var recieveCommand = &cobra.Command{
	Use:     "recieve [hash]",
	Aliases: []string{"get"},
	Short:   "Recieve file using secret hash",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := args[0]

		// 1. connecting to the relay via TCP
		conn, err := net.Dial("tcp", matchmakerTCP)
		if err != nil {
			fmt.Println("Error connecting to the matchmaker")
			return
		}
		defer conn.Close()

		// 2. Tell the matchmaker we are the reciever for the given hash
		nametag := fmt.Sprintf("RCV:%-6s ", hash)
		conn.Write([]byte(nametag[:11]))

		// 3. Stream the file from matchmaker to the disk
		outputFile := "downloaded_" + hash + ".txt"
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating file")
			return
		}
		defer file.Close()

		// 4. Download the stream
		fmt.Println("Downloading file...")
		io.Copy(file, conn)
		fmt.Println("Download completed")
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
