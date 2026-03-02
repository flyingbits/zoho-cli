package books

import (
	"github.com/urfave/cli/v3"
)

func baseCurrencyAdjustmentCmd() *cli.Command {
	return &cli.Command{
		Name:  "basecurrencyadjustment",
		Usage: "Base currency adjustment operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List base currency adjustment", Flags: listFlags(), Action: runList("basecurrencyadjustment")},
			{Name: "get", Usage: "Get a base currency adjustment", ArgsUsage: "<adjustment-id>", Action: runGet("basecurrencyadjustment")},
			{Name: "create", Usage: "Create a base currency adjustment", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("basecurrencyadjustment")},
			{Name: "delete", Usage: "Delete a base currency adjustment", ArgsUsage: "<adjustment-id>", Action: runDelete("basecurrencyadjustment")},
			{Name: "list-account-details", Usage: "List account details for base currency adjustment", Action: runList("basecurrencyadjustment/accountdetails")},
		},
	}
}
