package inventory

import (
	"fmt"
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func creditNotesCmd() *cli.Command {
	base := crudSubcommands("creditnotes")
	path := "/creditnotes/%s"
	extra := []*cli.Command{
		{Name: "email", Usage: "Email credit note", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invReq("POST", true, path+"/email")},
		{Name: "get-email-content", Usage: "Get email content", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/email")},
		{Name: "void", Usage: "Void credit note", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/void")},
		{Name: "convert-to-draft", Usage: "Convert to draft", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/draft")},
		{Name: "convert-to-open", Usage: "Convert to open", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/open")},
		{Name: "submit-for-approval", Usage: "Submit for approval", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/submit")},
		{Name: "approve", Usage: "Approve credit note", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/approve")},
		{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, path+"/billingaddress")},
		{Name: "update-shipping-address", Usage: "Update shipping address", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, path+"/shippingaddress")},
		{Name: "list-templates", Usage: "List templates", Action: invNoArg("GET", "/creditnotes/templates")},
		{Name: "update-template", Usage: "Update template", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, path+"/templates")},
		{Name: "apply-credits-to-invoices", Usage: "Apply credits to invoices", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/invoices/credits")},
		{Name: "list-invoices-credited", Usage: "List invoices credited", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/invoices/credited")},
		{Name: "delete-credits-applied", Usage: "Delete credits applied to invoice", ArgsUsage: "<credit-note-id> <invoice-id>", Action: invReq2("DELETE", "/creditnotes/%s/invoices/%s/credits")},
		{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/comments")},
		{Name: "list-comments", Usage: "List comments", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/comments")},
		{Name: "delete-comment", Usage: "Delete comment", ArgsUsage: "<credit-note-id> <comment-id>", Action: invReq2("DELETE", "/creditnotes/%s/comments/%s")},
		{Name: "list-refunds", Usage: "List credit note refunds", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/refunds")},
		{Name: "refund", Usage: "Refund credit note", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/refunds")},
		{Name: "list-refunds-of", Usage: "List refunds of credit note", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/refunds")},
		{Name: "update-refund", Usage: "Update credit note refund", ArgsUsage: "<credit-note-id> <refund-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq2JSON("PUT", "/creditnotes/%s/refunds/%s")},
		{Name: "get-refund", Usage: "Get credit note refund", ArgsUsage: "<credit-note-id> <refund-id>", Action: invReq2("GET", "/creditnotes/%s/refunds/%s")},
		{Name: "delete-refund", Usage: "Delete credit note refund", ArgsUsage: "<credit-note-id> <refund-id>", Action: invReq2("DELETE", "/creditnotes/%s/refunds/%s")},
	}
	return &cli.Command{
		Name:     "credit-notes",
		Usage:    "Credit note operations",
		Commands: append(base, extra...),
	}
}

func creditNoteCommentUpdate(_ context.Context, cmd *cli.Command) error {
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
	path := fmt.Sprintf("/creditnotes/%s/comments/%s", cmd.Args().Get(0), cmd.Args().Get(1))
	raw, err := req(c, orgID, "PUT", path, &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
