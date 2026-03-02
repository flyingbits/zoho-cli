package invoice

import (
	"github.com/urfave/cli/v3"
)

func taxesCmd() *cli.Command {
	return &cli.Command{
		Name:  "taxes",
		Usage: "Tax operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List taxes", Flags: listFlags(), Action: runList("taxes")},
			{Name: "get", Usage: "Get a tax", ArgsUsage: "<tax-id>", Action: runGet("taxes")},
			{Name: "create", Usage: "Create a tax", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("taxes")},
			{Name: "update", Usage: "Update a tax", ArgsUsage: "<tax-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("taxes")},
			{Name: "delete", Usage: "Delete a tax", ArgsUsage: "<tax-id>", Action: runDelete("taxes")},
			{Name: "update-group", Usage: "Update a tax group", ArgsUsage: "<tax-group-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("taxes/tax_groups")},
			{Name: "get-group", Usage: "Get a tax group", ArgsUsage: "<tax-group-id>", Action: runGet("taxes/tax_groups")},
			{Name: "delete-group", Usage: "Delete a tax group", ArgsUsage: "<tax-group-id>", Action: runDelete("taxes/tax_groups")},
			{Name: "create-group", Usage: "Create a tax group", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("taxes/tax_groups")},
			{Name: "create-exemption", Usage: "Create tax exemption", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("taxes/exemptions")},
			{Name: "list-exemptions", Usage: "List tax exemptions", Action: runList("taxes/exemptions")},
			{Name: "update-exemption", Usage: "Update tax exemption", ArgsUsage: "<exemption-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("taxes/exemptions")},
			{Name: "delete-exemption", Usage: "Delete tax exemption", ArgsUsage: "<exemption-id>", Action: runDelete("taxes/exemptions")},
			{Name: "create-authority", Usage: "Create tax authority", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("taxes/authorities")},
			{Name: "list-authorities", Usage: "List tax authorities", Action: runList("taxes/authorities")},
			{Name: "update-authority", Usage: "Update tax authority", ArgsUsage: "<authority-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("taxes/authorities")},
			{Name: "delete-authority", Usage: "Delete tax authority", ArgsUsage: "<authority-id>", Action: runDelete("taxes/authorities")},
		},
	}
}
