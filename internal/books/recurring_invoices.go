package books

import (
	"github.com/urfave/cli/v3"
)

func recurringInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurringinvoices",
		Usage: "Recurring invoice operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List recurring invoices", Flags: listFlags(), Action: runList("recurringinvoices")},
			{Name: "get", Usage: "Get a recurring invoice", ArgsUsage: "<recurringinvoice-id>", Action: runGet("recurringinvoices")},
			{Name: "create", Usage: "Create a recurring invoice", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("recurringinvoices")},
			{Name: "update", Usage: "Update a recurring invoice", ArgsUsage: "<recurringinvoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("recurringinvoices")},
			{Name: "delete", Usage: "Delete a recurring invoice", ArgsUsage: "<recurringinvoice-id>", Action: runDelete("recurringinvoices")},
			{Name: "stop", Usage: "Stop recurring invoice", ArgsUsage: "<recurringinvoice-id>", Action: runPost("recurringinvoices", "stop")},
			{Name: "resume", Usage: "Resume recurring invoice", ArgsUsage: "<recurringinvoice-id>", Action: runPost("recurringinvoices", "resume")},
			{Name: "list-history", Usage: "List recurring invoice history", ArgsUsage: "<recurringinvoice-id>", Action: runGetSub("recurringinvoices", "history")},
		},
	}
}
