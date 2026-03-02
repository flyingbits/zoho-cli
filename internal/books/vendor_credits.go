package books

import (
	"github.com/urfave/cli/v3"
)

func vendorCreditsCmd() *cli.Command {
	return &cli.Command{
		Name:  "vendorcredits",
		Usage: "Vendor credit operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List vendor credits", Flags: listFlags(), Action: runList("vendorcredits")},
			{Name: "get", Usage: "Get a vendor credit", ArgsUsage: "<vendorcredit-id>", Action: runGet("vendorcredits")},
			{Name: "create", Usage: "Create a vendor credit", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("vendorcredits")},
			{Name: "update", Usage: "Update a vendor credit", ArgsUsage: "<vendorcredit-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("vendorcredits")},
			{Name: "delete", Usage: "Delete a vendor credit", ArgsUsage: "<vendorcredit-id>", Action: runDelete("vendorcredits")},
			{Name: "convert-to-open", Usage: "Convert to open", ArgsUsage: "<vendorcredit-id>", Action: runPost("vendorcredits", "status/open")},
			{Name: "void", Usage: "Void vendor credit", ArgsUsage: "<vendorcredit-id>", Action: runPost("vendorcredits", "status/void")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<vendorcredit-id>", Action: runPost("vendorcredits", "submit")},
			{Name: "approve", Usage: "Approve vendor credit", ArgsUsage: "<vendorcredit-id>", Action: runPost("vendorcredits", "approve")},
			{Name: "list-bills-credited", Usage: "List bills credited", ArgsUsage: "<vendorcredit-id>", Action: runGetSub("vendorcredits", "billscredited")},
			{Name: "list-refunds", Usage: "List refunds of vendor credit", ArgsUsage: "<vendorcredit-id>", Action: runGetSub("vendorcredits", "refunds")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<vendorcredit-id>", Action: runGetSub("vendorcredits", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<vendorcredit-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("vendorcredits", "comments")},
		},
	}
}
