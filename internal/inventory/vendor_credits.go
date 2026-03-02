package inventory

import (
	"github.com/urfave/cli/v3"
)

func vendorCreditsCmd() *cli.Command {
	base := crudSubcommands("vendorcredits")
	path := "/vendorcredits/%s"
	extra := []*cli.Command{
		{Name: "convert-to-open", Usage: "Convert to open", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/open")},
		{Name: "void", Usage: "Void vendor credit", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/void")},
		{Name: "submit-for-approval", Usage: "Submit for approval", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/submit")},
		{Name: "approve", Usage: "Approve vendor credit", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/approve")},
		{Name: "apply-credits-to-bill", Usage: "Apply credits to bill", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/bills/credits")},
		{Name: "list-bills-credited", Usage: "List bills credited", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/bills/credited")},
		{Name: "delete-bills-credited", Usage: "Delete bills credited", ArgsUsage: "<vendor-credit-id> <bill-id>", Action: invReq2("DELETE", "/vendorcredits/%s/bills/%s/credits")},
		{Name: "refund", Usage: "Refund vendor credit", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/refunds")},
		{Name: "list-refunds", Usage: "List refunds of vendor credit", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/refunds")},
		{Name: "list-refund-details", Usage: "List vendor credit refunds", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/refunds")},
		{Name: "update-refund", Usage: "Update vendor credit refund", ArgsUsage: "<vendor-credit-id> <refund-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq2JSON("PUT", "/vendorcredits/%s/refunds/%s")},
		{Name: "get-refund", Usage: "Get vendor credit refund", ArgsUsage: "<vendor-credit-id> <refund-id>", Action: invReq2("GET", "/vendorcredits/%s/refunds/%s")},
		{Name: "delete-refund", Usage: "Delete vendor credit refund", ArgsUsage: "<vendor-credit-id> <refund-id>", Action: invReq2("DELETE", "/vendorcredits/%s/refunds/%s")},
		{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/comments")},
		{Name: "list-comments", Usage: "List comments and history", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/comments")},
		{Name: "delete-comment", Usage: "Delete comment", ArgsUsage: "<vendor-credit-id> <comment-id>", Action: invReq2("DELETE", "/vendorcredits/%s/comments/%s")},
	}
	return &cli.Command{
		Name:     "vendor-credits",
		Usage:    "Vendor credit operations",
		Commands: append(base, extra...),
	}
}

