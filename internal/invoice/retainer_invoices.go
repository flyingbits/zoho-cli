package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func retainerInvoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "retainer-invoices",
		Usage: "Retainer invoice operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List retainer invoices", Flags: listFlags(), 			Action: runList("retainer-invoices")},
			{Name: "get", Usage: "Get a retainer invoice", ArgsUsage: "<retainer-invoice-id>", Action: runGet("retainer-invoices")},
			{Name: "create", Usage: "Create a retainer invoice", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("retainer-invoices")},
			{Name: "update", Usage: "Update a retainer invoice", ArgsUsage: "<retainer-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("retainer-invoices")},
			{Name: "delete", Usage: "Delete a retainer invoice", ArgsUsage: "<retainer-invoice-id>", Action: runDelete("retainer-invoices")},
			{Name: "mark-sent", Usage: "Mark retainer invoice as sent", ArgsUsage: "<retainer-invoice-id>", Action: runPost("retainer-invoices", "status/sent")},
			{Name: "update-template", Usage: "Update retainer invoice template", ArgsUsage: "<retainer-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("retainer-invoices", "templates")},
			{Name: "void", Usage: "Void a retainer invoice", ArgsUsage: "<retainer-invoice-id>", Action: runPost("retainer-invoices", "status/void")},
			{Name: "mark-draft", Usage: "Mark as draft", ArgsUsage: "<retainer-invoice-id>", Action: runPost("retainer-invoices", "status/draft")},
			{Name: "email", Usage: "Email a retainer invoice", ArgsUsage: "<retainer-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("retainer-invoices", "email")},
			{Name: "get-email-content", Usage: "Get retainer invoice email content", ArgsUsage: "<retainer-invoice-id>", Action: runGetSub("retainer-invoices", "email")},
			{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<retainer-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("retainer-invoices", "address/billing")},
			{Name: "list-templates", Usage: "List retainer invoice templates", Action: retainerListTemplates},
			{Name: "add-attachment", Usage: "Add attachment", ArgsUsage: "<retainer-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("retainer-invoices", "attachment")},
			{Name: "get-attachment", Usage: "Get attachment", ArgsUsage: "<retainer-invoice-id>", Action: runGetSub("retainer-invoices", "attachment")},
			{Name: "delete-attachment", Usage: "Delete attachment", ArgsUsage: "<retainer-invoice-id>", Action: runPost("retainer-invoices", "attachment/delete")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<retainer-invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("retainer-invoices", "comments")},
			{Name: "list-comments", Usage: "List comments and history", ArgsUsage: "<retainer-invoice-id>", Action: runGetSub("retainer-invoices", "comments")},
			{Name: "update-comment", Usage: "Update comment", ArgsUsage: "<retainer-invoice-id> <comment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: retainerUpdateComment},
			{Name: "delete-comment", Usage: "Delete comment", ArgsUsage: "<retainer-invoice-id> <comment-id>", Action: retainerDeleteComment},
		},
	}
}

func retainerListTemplates(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/retainer-invoices/templates", nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func retainerUpdateComment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	var body any
	json.Unmarshal([]byte(cmd.String("json")), &body)
	raw, err := req(c, orgID, "PUT", "/retainer-invoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func retainerDeleteComment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/retainer-invoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
