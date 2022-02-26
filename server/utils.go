package server

import (
	"log"
	"net"
	"strings"
)

func writeToClient(conn net.Conn, data []byte) {
	// add a newline to the end of the data
	_, err := conn.Write(data)
	if err != nil {
		log.Printf("failed to write to client: %v", err)
	}
}

func parseRequest(scannedData string) []string {
	processedData := strings.ToLower(strings.TrimSpace(scannedData))
	return strings.Split(processedData, " ")
}
