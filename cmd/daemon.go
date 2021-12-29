package main

type DaemonCmd struct {
	Start   StartCmd   `cmd:"" help:"Start the daemon"`
	FindKey FindKeyCmd `cmd:"" help:"Scan for a key pattern"`
}
