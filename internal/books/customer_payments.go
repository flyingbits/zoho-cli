package books

import (
	"github.com/urfave/cli/v3"
)

func customerPaymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "customerpayments",
		Usage: "Customer payment operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List customer payments", Flags: listFlags(), Action: runList("customerpayments")},
			{Name: "get", Usage: "Get a customer payment", ArgsUsage: "<payment-id>", Action: runGet("customerpayments")},
			{Name: "create", Usage: "Create a customer payment", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("customerpayments")},
			{Name: "update", Usage: "Update a customer payment", ArgsUsage: "<payment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("customerpayments")},
			{Name: "delete", Usage: "Delete a customer payment", ArgsUsage: "<payment-id>", Action: runDelete("customerpayments")},
			{Name: "list-refunds", Usage: "List refunds of a customer payment", ArgsUsage: "<payment-id>", Action: runGetSub("customerpayments", "refunds")},
		},
	}
}
