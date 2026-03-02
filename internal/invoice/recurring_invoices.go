package invoice

import (
	"github.com/urfave/cli/v3"
)

func recurringInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurring-invoices",
		Usage: "Recurring invoice operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List recurring invoices", Flags: listFlags(), 			Action: runList("recurring-invoices")},
			{Name: "get", Usage: "Get a recurring invoice", ArgsUsage: "<recurring-invoice-id>", Action: runGet("recurring-invoices")},
			{Name: "create", Usage: "Create a recurring invoice", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("recurring-invoices")},
			{Name: "update", Usage: "Update a recurring invoice", ArgsUsage: "<recurring-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("recurring-invoices")},
			{Name: "delete", Usage: "Delete a recurring invoice", ArgsUsage: "<recurring-invoice-id>", Action: runDelete("recurring-invoices")},
			{Name: "stop", Usage: "Stop a recurring invoice", ArgsUsage: "<recurring-invoice-id>", Action: runPost("recurring-invoices", "status/stop")},
			{Name: "resume", Usage: "Resume a recurring invoice", ArgsUsage: "<recurring-invoice-id>", Action: runPost("recurring-invoices", "status/resume")},
			{Name: "update-template", Usage: "Update recurring invoice template", ArgsUsage: "<recurring-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("recurring-invoices", "templates")},
			{Name: "list-history", Usage: "List recurring invoice history", ArgsUsage: "<recurring-invoice-id>", Action: runGetSub("recurring-invoices", "history")},
		},
	}
}
