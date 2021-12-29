package config

import (
	"fmt"
	"net"
)

type IPAddress string
type Port int

func (a IPAddress) Validate() error {
	if net.ParseIP(string(a)) == nil {
		return fmt.Errorf("Invalid IP address %s", a)
	}
	return nil
}

func (p Port) Validate() error {
	maxPort := 2 << 15
	if p <= 0 || int(p) >= maxPort {
		return fmt.Errorf("Invalid port %d", p)
	}
	return nil
}
