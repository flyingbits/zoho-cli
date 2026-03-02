package inventory

import (
	"context"
	"encoding/json"
	"fmt"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func retainerInvoicesCmd() *cli.Command {
	base := crudSubcommands("retainerinvoices")
	path := "/retainerinvoices/%s"
	extra := []*cli.Command{
		{Name: "mark-sent", Usage: "Mark as sent", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/sent")},
		{Name: "void", Usage: "Void", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/void")},
		{Name: "mark-draft", Usage: "Mark as draft", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/draft")},
		{Name: "submit-for-approval", Usage: "Submit for approval", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/submit")},
		{Name: "approve", Usage: "Approve", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/approve")},
		{Name: "email", Usage: "Email", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invReq("POST", true, path+"/email")},
		{Name: "get-email-content", Usage: "Get email content", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/email")},
		{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, path+"/billingaddress")},
		{Name: "list-templates", Usage: "List templates", Action: invNoArg("GET", "/retainerinvoices/templates")},
		{Name: "add-attachment", Usage: "Add attachment", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invReq("POST", true, path+"/attachment")},
		{Name: "get-attachment", Usage: "Get attachment", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/attachment")},
		{Name: "delete-attachment", Usage: "Delete attachment", ArgsUsage: "<id>", Action: invReq("DELETE", false, path+"/attachment")},
		{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, path+"/comments")},
		{Name: "list-comments", Usage: "List comments", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/comments")},
		{Name: "update-comment", Usage: "Update comment", ArgsUsage: "<invoice-id> <comment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: retainerCommentUpdate},
		{Name: "delete-comment", Usage: "Delete comment", ArgsUsage: "<invoice-id> <comment-id>", Action: invReq2("DELETE", "/retainerinvoices/%s/comments/%s")},
	}
	return &cli.Command{
		Name:     "retainer-invoices",
		Usage:    "Retainer invoice operations",
		Commands: append(base, extra...),
	}
}

func retainerCommentUpdate(_ context.Context, cmd *cli.Command) error {
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
	path := fmt.Sprintf("/retainerinvoices/%s/comments/%s", cmd.Args().Get(0), cmd.Args().Get(1))
	raw, err := req(c, orgID, "PUT", path, &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
