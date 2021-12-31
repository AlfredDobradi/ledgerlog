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

	Start   StartCmd   `cmd:"" help:"Start the daemon"`
	FindKey FindKeyCmd `cmd:"" help:"Scan for a key pattern"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Description("LedgerD - A journaling microblog server"),
		kong.Name("ledgerd"), kong.UsageOnError(),
	)

	if err := config.Parse(CLI.ConfigPath); err != nil {
		log.Panicln(err)
	}

	err := ctx.Run(&Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
