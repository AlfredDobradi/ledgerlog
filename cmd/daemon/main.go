package main

import (
	"log"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/alecthomas/kong"
)

var (
	// tag is the version tag created by go-semrel
	tag = "v0.0.0"

	// commitHash is the HEAD commit when the application was compiled
	commitHash = "00000000"

	// buildTime is the full date time when the application was compiled
	buildTime = ""
)

type Context struct {
	Debug bool
}

var CLI struct {
	Debug      bool   `help:"Enable debug mode"`
	ConfigPath string `help:"Path to TOML config file" type:"existingfile" name:"cfg" env:"LEDGER_CFG" short:"c"`

	Start   StartCmd   `cmd:"" help:"Start the daemon"`
	Version VersionCmd `cmd:"" help:"Version information"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Description("LedgerD - A journaling microblog server"),
		kong.Name("ledgerd"), kong.UsageOnError(),
	)

	if CLI.ConfigPath != "" {
		if err := config.Parse(CLI.ConfigPath); err != nil {
			log.Panicln(err)
		}
	}

	err := ctx.Run(&Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
