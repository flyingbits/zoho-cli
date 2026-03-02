package books

import (
	"github.com/urfave/cli/v3"
)

func itemsCmd() *cli.Command {
	return &cli.Command{
		Name:  "items",
		Usage: "Item operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List items", Flags: listFlags(), Action: runList("items")},
			{Name: "get", Usage: "Get an item", ArgsUsage: "<item-id>", Action: runGet("items")},
			{Name: "create", Usage: "Create an item", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("items")},
			{Name: "update", Usage: "Update an item", ArgsUsage: "<item-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("items")},
			{Name: "delete", Usage: "Delete an item", ArgsUsage: "<item-id>", Action: runDelete("items")},
			{Name: "mark-active", Usage: "Mark item as active", ArgsUsage: "<item-id>", Action: runPost("items", "active")},
			{Name: "mark-inactive", Usage: "Mark item as inactive", ArgsUsage: "<item-id>", Action: runPost("items", "inactive")},
		},
	}
}
