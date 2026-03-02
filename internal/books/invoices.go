package books

import (
	"github.com/urfave/cli/v3"
)

func invoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "invoices",
		Usage: "Invoice operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List invoices", Flags: listFlags(), Action: runList("invoices")},
			{Name: "get", Usage: "Get an invoice", ArgsUsage: "<invoice-id>", Action: runGet("invoices")},
			{Name: "create", Usage: "Create an invoice", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("invoices")},
			{Name: "update", Usage: "Update an invoice", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("invoices")},
			{Name: "delete", Usage: "Delete an invoice", ArgsUsage: "<invoice-id>", Action: runDelete("invoices")},
			{Name: "mark-sent", Usage: "Mark as sent", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "status/sent")},
			{Name: "void", Usage: "Void invoice", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "status/void")},
			{Name: "mark-draft", Usage: "Mark as draft", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "status/draft")},
			{Name: "email", Usage: "Email invoice", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("invoices", "email")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "submit")},
			{Name: "approve", Usage: "Approve invoice", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "approve")},
			{Name: "list-payments", Usage: "List invoice payments", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "payments")},
			{Name: "list-credits-applied", Usage: "List credits applied", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "creditsapplied")},
			{Name: "apply-credits", Usage: "Apply credits", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("invoices", "credits")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("invoices", "comments")},
			{Name: "payment-link", Usage: "Generate payment link", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("invoices", "paymentlink")},
		},
	}
}
