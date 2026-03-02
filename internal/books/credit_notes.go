package books

import (
	"github.com/urfave/cli/v3"
)

func creditNotesCmd() *cli.Command {
	return &cli.Command{
		Name:  "creditnotes",
		Usage: "Credit note operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List credit notes", Flags: listFlags(), Action: runList("creditnotes")},
			{Name: "get", Usage: "Get a credit note", ArgsUsage: "<creditnote-id>", Action: runGet("creditnotes")},
			{Name: "create", Usage: "Create a credit note", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("creditnotes")},
			{Name: "update", Usage: "Update a credit note", ArgsUsage: "<creditnote-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("creditnotes")},
			{Name: "delete", Usage: "Delete a credit note", ArgsUsage: "<creditnote-id>", Action: runDelete("creditnotes")},
			{Name: "email", Usage: "Email credit note", ArgsUsage: "<creditnote-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("creditnotes", "email")},
			{Name: "void", Usage: "Void credit note", ArgsUsage: "<creditnote-id>", Action: runPost("creditnotes", "status/void")},
			{Name: "convert-to-open", Usage: "Convert to open", ArgsUsage: "<creditnote-id>", Action: runPost("creditnotes", "status/open")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<creditnote-id>", Action: runPost("creditnotes", "submit")},
			{Name: "approve", Usage: "Approve credit note", ArgsUsage: "<creditnote-id>", Action: runPost("creditnotes", "approve")},
			{Name: "list-invoices-credited", Usage: "List invoices credited", ArgsUsage: "<creditnote-id>", Action: runGetSub("creditnotes", "invoicescredited")},
			{Name: "list-refunds", Usage: "List credit note refunds", ArgsUsage: "<creditnote-id>", Action: runGetSub("creditnotes", "refunds")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<creditnote-id>", Action: runGetSub("creditnotes", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<creditnote-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("creditnotes", "comments")},
		},
	}
}
