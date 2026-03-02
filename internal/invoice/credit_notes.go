package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func creditNotesCmd() *cli.Command {
	return &cli.Command{
		Name:  "credit-notes",
		Usage: "Credit note operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List credit notes", Flags: listFlags(), 			Action: runList("credit-notes")},
			{Name: "get", Usage: "Get a credit note", ArgsUsage: "<credit-note-id>", Action: runGet("credit-notes")},
			{Name: "create", Usage: "Create a credit note", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("credit-notes")},
			{Name: "update", Usage: "Update a credit note", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("credit-notes")},
			{Name: "delete", Usage: "Delete a credit note", ArgsUsage: "<credit-note-id>", Action: runDelete("credit-notes")},
			{Name: "email", Usage: "Email a credit note", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("credit-notes", "email")},
			{Name: "void", Usage: "Void a credit note", ArgsUsage: "<credit-note-id>", Action: runPost("credit-notes", "status/void")},
			{Name: "open-voided", Usage: "Open a voided credit note", ArgsUsage: "<credit-note-id>", Action: runPost("credit-notes", "status/open")},
			{Name: "email-history", Usage: "Email history", ArgsUsage: "<credit-note-id>", Action: runGetSub("credit-notes", "emailhistory")},
			{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("credit-notes", "address/billing")},
			{Name: "update-shipping-address", Usage: "Update shipping address", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("credit-notes", "address/shipping")},
			{Name: "list-templates", Usage: "List credit note templates", ArgsUsage: "<credit-note-id>", Action: runGetSub("credit-notes", "templates")},
			{Name: "update-template", Usage: "Update credit note template", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("credit-notes", "templates")},
			{Name: "credit-to-invoice", Usage: "Credit to an invoice", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("credit-notes", "invoices")},
			{Name: "list-invoices-credited", Usage: "List invoices credited", ArgsUsage: "<credit-note-id>", Action: runGetSub("credit-notes", "invoices")},
			{Name: "delete-invoice-credited", Usage: "Delete invoices credited", ArgsUsage: "<credit-note-id> <invoice-id>", Action: creditNoteDeleteInvoiced},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("credit-notes", "comments")},
			{Name: "list-comments", Usage: "List credit note comments and history", ArgsUsage: "<credit-note-id>", Action: runGetSub("credit-notes", "comments")},
			{Name: "delete-comment", Usage: "Delete a comment", ArgsUsage: "<credit-note-id> <comment-id>", Action: creditNoteDeleteComment},
			{Name: "list-refunds", Usage: "List credit note refunds", ArgsUsage: "<credit-note-id>", Action: runGetSub("credit-notes", "refunds")},
			{Name: "refund", Usage: "Refund credit note", ArgsUsage: "<credit-note-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("credit-notes", "refunds")},
			{Name: "update-refund", Usage: "Update credit note refund", ArgsUsage: "<credit-note-id> <refund-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: creditNoteUpdateRefund},
			{Name: "get-refund", Usage: "Get credit note refund", ArgsUsage: "<credit-note-id> <refund-id>", Action: creditNoteGetRefund},
			{Name: "delete-refund", Usage: "Delete credit note refund", ArgsUsage: "<credit-note-id> <refund-id>", Action: creditNoteDeleteRefund},
		},
	}
}

func creditNoteDeleteInvoiced(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/credit-notes/"+cmd.Args().First()+"/invoices/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func creditNoteDeleteComment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/credit-notes/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func creditNoteUpdateRefund(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/credit-notes/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func creditNoteGetRefund(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/credit-notes/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func creditNoteDeleteRefund(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/credit-notes/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
