package p2p

import (
	"fmt" // Formatting strings (printing to console)
	"io"  // Input/Output utilities (streaming data)
	"net" // Networking (TCP/UDP sockets)
	"os"  // Operating System (reading files from disk)
)

func StartServer(port string, filePath string) error {
	// 1. ask the os to listen on a TSP port
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Server is listening on port %s... : ", port)

	for {
		// 2. waiting for the reciever to connect
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Errer accepting connection, ", err)
			continue
		}

		fmt.Println("Reciever connected")

		// 3. Handle the file transfer in a goRoutine (concurrent)
		go handleSend(conn, filePath)
	}
}

func handleSend(conn net.Conn, filePath string) {
	defer conn.Close()

	// 4. open the file user wants to send
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("An error occured while opening file ", err)
		return
	}
	defer file.Close()

	// 5. send the file over the connection
	res, err := io.Copy(conn, file)
	if err != nil {
		fmt.Println("An error occured while sending file")
		return
	}

	fmt.Printf("Successfully sent %d bytes to %s\n", res, conn.RemoteAddr())
}
