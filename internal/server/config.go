package server

import "github.com/AlfredDobradi/ledgerlog/internal/config"

var (
	ipAddress  config.IPAddress = "0.0.0.0"
	port       config.Port      = 8080
	publicPath string           = "./public"
)

func IPAddress() config.IPAddress {
	return ipAddress
}

func Port() config.Port {
	return port
}

func PublicPath() string {
	return publicPath
}

func SetIPAddress(newIP config.IPAddress) {
	ipAddress = newIP
}

func SetPort(newPort config.Port) {
	port = newPort
}

func SetPublicPath(newPublicPath string) {
	publicPath = newPublicPath
}
