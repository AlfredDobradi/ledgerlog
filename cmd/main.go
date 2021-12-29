package main

import (
	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}

var CLI struct {
	Debug bool `help:"Enable debug mode"`

	Client ClientCmd `cmd:"" help:"Client related commands"`
	Daemon DaemonCmd `cmd:"" help:"Daemon related commands"`
}

func main() {
	ctx := kong.Parse(&CLI, kong.Description("Ledgerlog - A journaling microblog"), kong.Name("ledgerlog"), kong.UsageOnError())
	err := ctx.Run(&Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
