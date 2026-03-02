package books

import (
	"github.com/urfave/cli/v3"
)

func retainerInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "retainerinvoices",
		Usage: "Retainer invoice operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List retainer invoices", Flags: listFlags(), Action: runList("retainerinvoices")},
			{Name: "get", Usage: "Get a retainer invoice", ArgsUsage: "<retainerinvoice-id>", Action: runGet("retainerinvoices")},
			{Name: "create", Usage: "Create a retainer invoice", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("retainerinvoices")},
			{Name: "update", Usage: "Update a retainer invoice", ArgsUsage: "<retainerinvoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("retainerinvoices")},
			{Name: "delete", Usage: "Delete a retainer invoice", ArgsUsage: "<retainerinvoice-id>", Action: runDelete("retainerinvoices")},
			{Name: "mark-sent", Usage: "Mark as sent", ArgsUsage: "<retainerinvoice-id>", Action: runPost("retainerinvoices", "status/sent")},
			{Name: "void", Usage: "Void retainer invoice", ArgsUsage: "<retainerinvoice-id>", Action: runPost("retainerinvoices", "status/void")},
			{Name: "mark-draft", Usage: "Mark as draft", ArgsUsage: "<retainerinvoice-id>", Action: runPost("retainerinvoices", "status/draft")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<retainerinvoice-id>", Action: runPost("retainerinvoices", "submit")},
			{Name: "approve", Usage: "Approve retainer invoice", ArgsUsage: "<retainerinvoice-id>", Action: runPost("retainerinvoices", "approve")},
			{Name: "email", Usage: "Email retainer invoice", ArgsUsage: "<retainerinvoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("retainerinvoices", "email")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<retainerinvoice-id>", Action: runGetSub("retainerinvoices", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<retainerinvoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("retainerinvoices", "comments")},
		},
	}
}
