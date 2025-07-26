package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func ipToUint32(ip net.IP) uint32 {
	ip4 := ip.To4()
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
}

func isGCNAT(ipStr string) bool {
	ipAddr := net.ParseIP(ipStr)
	if ipAddr == nil {
		return false
	}
	ipNum := ipToUint32(ipAddr)
	const (
		cgnatStart = 0x64400000
		cgnatEnd   = 0x647FFFFF
	)
	return ipNum >= cgnatStart && ipNum <= cgnatEnd
}

func getExternalIP() (string, error) {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", fmt.Errorf("failed to get external IP: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

func getLocalIP() (string, error) {
	conn, err := net.DialTimeout("udp", "8.8.8.8:80", time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to dial UDP for local IP: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func main() {
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Printf("Error getting local IP: %v\n", err)
		return
	}

	extIP, err := getExternalIP()
	if err != nil {
		fmt.Printf("Error getting external IP: %v\n", err)
		return
	}

	cgnatStatus := "No."
	if isGCNAT(extIP) {
		cgnatStatus = "Yes."
	}

	fmt.Printf("Local IP: %s\nExternal IP: %s\nCGNAT: %s\n", localIP, extIP, cgnatStatus)

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}