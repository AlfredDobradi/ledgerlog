package main

import (
	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}

var CLI struct {
	Debug bool `help:"Enable debug mode"`

	Register RegisterCmd `cmd:"" help:"Register email with public key"`
	Send     SendCmd     `cmd:"" help:"Send new post"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
