package inventory

import (
	"github.com/urfave/cli/v3"
)

func itemAdjustmentsCmd() *cli.Command {
	return &cli.Command{
		Name:     "item-adjustments",
		Usage:    "Item adjustment operations",
		Commands: crudSubcommands("itemadjustments"),
	}
}
