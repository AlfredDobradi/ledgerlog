package server

import "github.com/AlfredDobradi/ledgerlog/internal/config"

var (
	ipAddress config.IPAddress = "0.0.0.0"
	port      config.Port      = 8080
)

func IPAddress() config.IPAddress {
	return ipAddress
}

func Port() config.Port {
	return port
}

func SetIPAddress(newIP config.IPAddress) {
	ipAddress = newIP
}

func SetPort(newPort config.Port) {
	port = newPort
}
