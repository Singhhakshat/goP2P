package p2p

import (
	"fmt"
	"io"
	"net"
	"os"
)

func ReceiveFile(ip string, port string, outputFileName string) error {
	address := ip + ":" + port
	fmt.Printf("Connecting with sender at %s\n", address)

	// 1.  dial to connect to the server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("Failed to connect to the sender %v", err)
	}
	defer conn.Close()

	// 2. create an empty file on the reciever's harddrive
	file, err := os.Create("goP2P/test/reciever" + outputFileName)
	if err != nil {
		return fmt.Errorf("Error creating file %v", err)
	}
	defer file.Close()

	// 3. Stream the data from the network directly into the new file
	n, err := io.Copy(file, conn)

	if err != nil {
		return fmt.Errorf("Error writing file %v", err)
	}

	fmt.Printf("Successfully recieved %d bytes and saved as %s\n", n, outputFileName)

	return nil
}
