package inventory

import (
	"github.com/urfave/cli/v3"
)

func currencyCmd() *cli.Command {
	return &cli.Command{
		Name:     "currency",
		Usage:    "Currency operations",
		Commands: crudSubcommands("currency"),
	}
}
