package inventory

import (
	"github.com/urfave/cli/v3"
)

func purchaseReceivesCmd() *cli.Command {
	return &cli.Command{
		Name:     "purchase-receives",
		Usage:    "Purchase receive operations",
		Commands: crudSubcommands("purchasereceives"),
	}
}
