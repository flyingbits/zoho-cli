package books

import (
	"github.com/urfave/cli/v3"
)

func vendorPaymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "vendorpayments",
		Usage: "Vendor payment operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List vendor payments", Flags: listFlags(), Action: runList("vendorpayments")},
			{Name: "get", Usage: "Get a vendor payment", ArgsUsage: "<vendorpayment-id>", Action: runGet("vendorpayments")},
			{Name: "create", Usage: "Create a vendor payment", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("vendorpayments")},
			{Name: "update", Usage: "Update a vendor payment", ArgsUsage: "<vendorpayment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("vendorpayments")},
			{Name: "delete", Usage: "Delete a vendor payment", ArgsUsage: "<vendorpayment-id>", Action: runDelete("vendorpayments")},
			{Name: "list-refunds", Usage: "List refunds of a vendor payment", ArgsUsage: "<vendorpayment-id>", Action: runGetSub("vendorpayments", "refunds")},
			{Name: "email", Usage: "Email vendor payment", ArgsUsage: "<vendorpayment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("vendorpayments", "email")},
		},
	}
}
