package server

import (
	"net"
	"strings"
)

func writeToClient(conn net.Conn, data []byte) error {
	return nil
}

func parseRequest(scannedData string) []string {
	processedData := strings.ToLower(strings.TrimSpace(scannedData))
	return strings.Split(processedData, " ")
}
