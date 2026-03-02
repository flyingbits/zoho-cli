package books

import (
	"github.com/urfave/cli/v3"
)

func fixedAssetsCmd() *cli.Command {
	return &cli.Command{
		Name:  "fixedassets",
		Usage: "Fixed asset operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "Get fixed asset list", Flags: listFlags(), Action: runList("fixedassets")},
			{Name: "get", Usage: "Get a fixed asset", ArgsUsage: "<fixedasset-id>", Action: runGet("fixedassets")},
			{Name: "create", Usage: "Create a fixed asset", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("fixedassets")},
			{Name: "update", Usage: "Update a fixed asset", ArgsUsage: "<fixedasset-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("fixedassets")},
			{Name: "delete", Usage: "Delete a fixed asset", ArgsUsage: "<fixedasset-id>", Action: runDelete("fixedassets")},
			{Name: "get-history", Usage: "Get fixed asset history", ArgsUsage: "<fixedasset-id>", Action: runGetSub("fixedassets", "history")},
			{Name: "get-forecast-depreciation", Usage: "Get forecast depreciation", ArgsUsage: "<fixedasset-id>", Action: runGetSub("fixedassets", "depreciation")},
			{Name: "mark-active", Usage: "Mark fixed asset as active", ArgsUsage: "<fixedasset-id>", Action: runPost("fixedassets", "active")},
			{Name: "cancel", Usage: "Cancel fixed asset", ArgsUsage: "<fixedasset-id>", Action: runPost("fixedassets", "inactive")},
			{Name: "mark-draft", Usage: "Mark fixed asset as draft", ArgsUsage: "<fixedasset-id>", Action: runPost("fixedassets", "draft")},
			{Name: "write-off", Usage: "Write off fixed asset", ArgsUsage: "<fixedasset-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("fixedassets", "writeoff")},
			{Name: "sell", Usage: "Sell fixed asset", ArgsUsage: "<fixedasset-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("fixedassets", "sell")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<fixedasset-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("fixedassets", "comments")},
		},
	}
}
