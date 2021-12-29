package main

type ClientCmd struct {
	Register RegisterCmd `cmd:"" help:"Register email with public key"`
	Send     SendCmd     `cmd:"" help:"Send new post"`
}
