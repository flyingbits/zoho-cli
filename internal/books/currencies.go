package books

import (
	"github.com/urfave/cli/v3"
)

func currenciesCmd() *cli.Command {
	return &cli.Command{
		Name:  "currencies",
		Usage: "Currency operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List currencies", Flags: listFlags(), Action: runList("settings/currencies")},
			{Name: "get", Usage: "Get a currency", ArgsUsage: "<currency-id>", Action: runGet("settings/currencies")},
			{Name: "create", Usage: "Create a currency", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("settings/currencies")},
			{Name: "update", Usage: "Update a currency", ArgsUsage: "<currency-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("settings/currencies")},
			{Name: "delete", Usage: "Delete a currency", ArgsUsage: "<currency-id>", Action: runDelete("settings/currencies")},
			{Name: "create-exchange-rate", Usage: "Create exchange rate", ArgsUsage: "<currency-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("settings/currencies", "exchangerates")},
			{Name: "list-exchange-rates", Usage: "List exchange rates", ArgsUsage: "<currency-id>", Action: runGetSub("settings/currencies", "exchangerates")},
			{Name: "update-exchange-rate", Usage: "Update exchange rate", ArgsUsage: "<currency-id> <exchange-rate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdateSub("settings/currencies", "exchangerates")},
		},
	}
}
