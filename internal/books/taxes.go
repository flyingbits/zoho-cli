package books

import (
	"github.com/urfave/cli/v3"
)

func taxesCmd() *cli.Command {
	return &cli.Command{
		Name:  "taxes",
		Usage: "Tax operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List taxes", Flags: listFlags(), Action: runList("settings/taxes")},
			{Name: "get", Usage: "Get a tax", ArgsUsage: "<tax-id>", Action: runGet("settings/taxes")},
			{Name: "create", Usage: "Create a tax", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("settings/taxes")},
			{Name: "update", Usage: "Update a tax", ArgsUsage: "<tax-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("settings/taxes")},
			{Name: "delete", Usage: "Delete a tax", ArgsUsage: "<tax-id>", Action: runDelete("settings/taxes")},
			{Name: "get-tax-group", Usage: "Get a tax group", ArgsUsage: "<tax-group-id>", Action: runGet("settings/taxgroups")},
			{Name: "create-tax-group", Usage: "Create a tax group", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("settings/taxgroups")},
			{Name: "update-tax-group", Usage: "Update a tax group", ArgsUsage: "<tax-group-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("settings/taxgroups")},
			{Name: "delete-tax-group", Usage: "Delete a tax group", ArgsUsage: "<tax-group-id>", Action: runDelete("settings/taxgroups")},
		},
	}
}
