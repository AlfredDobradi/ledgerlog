package main

import (
	"log"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}

var CLI struct {
	Debug      bool   `help:"Enable debug mode"`
	ConfigPath string `help:"Path to JSON config file" type:"existingfile" name:"cfg" env:"LEDGER_CFG"`

	Client ClientCmd `cmd:"" help:"Client related commands"`
	Daemon DaemonCmd `cmd:"" help:"Daemon related commands"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Description("Ledgerlog - A journaling microblog"),
		kong.Name("ledgerlog"), kong.UsageOnError(),
	)

	if CLI.ConfigPath != "" {
		if err := config.Parse(CLI.ConfigPath); err != nil {
			log.Panicln(err)
		}
	}

	err := ctx.Run(&Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
