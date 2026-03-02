package books

import (
	"github.com/urfave/cli/v3"
)

func billsCmd() *cli.Command {
	return &cli.Command{
		Name:  "bills",
		Usage: "Bill operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List bills", Flags: listFlags(), Action: runList("bills")},
			{Name: "get", Usage: "Get a bill", ArgsUsage: "<bill-id>", Action: runGet("bills")},
			{Name: "create", Usage: "Create a bill", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("bills")},
			{Name: "update", Usage: "Update a bill", ArgsUsage: "<bill-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("bills")},
			{Name: "delete", Usage: "Delete a bill", ArgsUsage: "<bill-id>", Action: runDelete("bills")},
			{Name: "void", Usage: "Void bill", ArgsUsage: "<bill-id>", Action: runPost("bills", "status/void")},
			{Name: "mark-open", Usage: "Mark as open", ArgsUsage: "<bill-id>", Action: runPost("bills", "status/open")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<bill-id>", Action: runPost("bills", "submit")},
			{Name: "approve", Usage: "Approve bill", ArgsUsage: "<bill-id>", Action: runPost("bills", "approve")},
			{Name: "list-payments", Usage: "List bill payments", ArgsUsage: "<bill-id>", Action: runGetSub("bills", "payments")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<bill-id>", Action: runGetSub("bills", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<bill-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("bills", "comments")},
			{Name: "convert-po-to-bill", Usage: "Convert PO to bill", ArgsUsage: "<bill-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("bills", "purchaseorder")},
		},
	}
}
