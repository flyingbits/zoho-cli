package books

import (
	"github.com/urfave/cli/v3"
)

func customerDebitNotesCmd() *cli.Command {
	return &cli.Command{
		Name:  "customerdebitnotes",
		Usage: "Customer debit note operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List customer debit notes", Flags: listFlags(), Action: runList("customerdebitnotes")},
			{Name: "get", Usage: "Get a customer debit note", ArgsUsage: "<customerdebitnote-id>", Action: runGet("customerdebitnotes")},
			{Name: "create", Usage: "Create a customer debit note", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("customerdebitnotes")},
			{Name: "update", Usage: "Update a customer debit note", ArgsUsage: "<customerdebitnote-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("customerdebitnotes")},
			{Name: "delete", Usage: "Delete a customer debit note", ArgsUsage: "<customerdebitnote-id>", Action: runDelete("customerdebitnotes")},
		},
	}
}
