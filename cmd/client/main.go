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
	ConfigPath string `help:"Path to TOML config file" type:"existingfile" required:"" name:"cfg" env:"LEDGER_CFG" short:"c"`

	Register RegisterCmd `cmd:"" help:"Register email with public key"`
	Send     SendCmd     `cmd:"" help:"Send new post"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Description("Ledgerlog Client - A journaling microblog client"),
		kong.Name("ledgerlog"), kong.UsageOnError(),
	)

	if err := config.Parse(CLI.ConfigPath); err != nil {
		log.Panicln(err)
	}

	err := ctx.Run(&Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
