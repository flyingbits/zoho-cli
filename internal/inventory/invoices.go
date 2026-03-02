package inventory

import (
	"context"
	"encoding/json"
	"fmt"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func invoicesCmd() *cli.Command {
	base := crudSubcommands("invoices")
	extra := []*cli.Command{
		{Name: "update-custom-field", Usage: "Update custom field", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", false, "/invoices/%s/customfield")},
		{Name: "mark-sent", Usage: "Mark as sent", ArgsUsage: "<id>", Action: invReq("POST", false, "/invoices/%s/sent")},
		{Name: "void", Usage: "Void invoice", ArgsUsage: "<id>", Action: invReq("POST", false, "/invoices/%s/void")},
		{Name: "mark-draft", Usage: "Mark as draft", ArgsUsage: "<id>", Action: invReq("POST", false, "/invoices/%s/draft")},
		{Name: "email", Usage: "Email invoice", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invReq("POST", true, "/invoices/%s/email")},
		{Name: "get-email-content", Usage: "Get email content", ArgsUsage: "<id>", Action: invReq("GET", false, "/invoices/%s/email")},
		{Name: "email-bulk", Usage: "Email invoices", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invBulk("POST", "/invoices/email")},
		{Name: "get-payment-reminder-content", Usage: "Get payment reminder content", ArgsUsage: "<id>", Action: invReq("GET", false, "/invoices/%s/paymentreminder")},
		{Name: "disable-payment-reminder", Usage: "Disable payment reminder", ArgsUsage: "<id>", Action: invReq("POST", false, "/invoices/%s/paymentreminder/disable")},
		{Name: "enable-payment-reminder", Usage: "Enable payment reminder", ArgsUsage: "<id>", Action: invReq("POST", false, "/invoices/%s/paymentreminder/enable")},
		{Name: "write-off", Usage: "Write off invoice", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invReq("POST", true, "/invoices/%s/writeoff")},
		{Name: "cancel-write-off", Usage: "Cancel write off", ArgsUsage: "<id>", Action: invReq("POST", false, "/invoices/%s/writeoff/cancel")},
		{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, "/invoices/%s/billingaddress")},
		{Name: "update-shipping-address", Usage: "Update shipping address", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, "/invoices/%s/shippingaddress")},
		{Name: "list-templates", Usage: "List invoice templates", Action: invNoArg("GET", "/invoices/templates")},
		{Name: "update-template", Usage: "Update template", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, "/invoices/%s/templates")},
		{Name: "list-payments", Usage: "List payments", ArgsUsage: "<id>", Action: invReq("GET", false, "/invoices/%s/payments")},
		{Name: "list-credits-applied", Usage: "List credits applied", ArgsUsage: "<id>", Action: invReq("GET", false, "/invoices/%s/creditsapplied")},
		{Name: "apply-credits", Usage: "Apply credits", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, "/invoices/%s/credits")},
		{Name: "delete-payment", Usage: "Delete payment", ArgsUsage: "<invoice-id> <payment-id>", Action: invReq2("DELETE", "/invoices/%s/payments/%s")},
		{Name: "delete-applied-credit", Usage: "Delete applied credit", ArgsUsage: "<invoice-id> <credit-id>", Action: invReq2("DELETE", "/invoices/%s/creditsapplied/%s")},
		{Name: "add-attachment", Usage: "Add attachment", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invReq("POST", true, "/invoices/%s/attachment")},
		{Name: "update-attachment-preference", Usage: "Update attachment preference", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, "/invoices/%s/attachment")},
		{Name: "get-attachment", Usage: "Get attachment", ArgsUsage: "<id>", Action: invReq("GET", false, "/invoices/%s/attachment")},
		{Name: "delete-attachment", Usage: "Delete attachment", ArgsUsage: "<id>", Action: invReq("DELETE", false, "/invoices/%s/attachment")},
		{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("POST", true, "/invoices/%s/comments")},
		{Name: "list-comments", Usage: "List comments and history", ArgsUsage: "<id>", Action: invReq("GET", false, "/invoices/%s/comments")},
		{Name: "update-comment", Usage: "Update comment", ArgsUsage: "<invoice-id> <comment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq2JSON("PUT", "/invoices/%s/comments/%s")},
		{Name: "delete-comment", Usage: "Delete comment", ArgsUsage: "<invoice-id> <comment-id>", Action: invReq2("DELETE", "/invoices/%s/comments/%s")},
		{Name: "bulk-export", Usage: "Bulk export", Flags: []cli.Flag{&cli.StringFlag{Name: "json"}}, Action: invBulk("POST", "/invoices/bulk-export")},
		{Name: "bulk-print", Usage: "Bulk print", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invBulk("POST", "/invoices/bulk-print")},
	}
	return &cli.Command{
		Name:     "invoices",
		Usage:    "Invoice operations",
		Commands: append(base, extra...),
	}
}

func invReq(method string, withJSON bool, pathFmt string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		path := fmt.Sprintf(pathFmt, cmd.Args().First())
		opts := (*zohttp.RequestOpts)(nil)
		if withJSON && cmd.String("json") != "" {
			var body any
			json.Unmarshal([]byte(cmd.String("json")), &body)
			opts = &zohttp.RequestOpts{JSON: body}
		}
		raw, err := req(c, orgID, method, path, opts)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func invReq2(method, pathFmt string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		path := fmt.Sprintf(pathFmt, cmd.Args().Get(0), cmd.Args().Get(1))
		raw, err := req(c, orgID, method, path, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func invReq2JSON(method, pathFmt string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
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
		path := fmt.Sprintf(pathFmt, cmd.Args().Get(0), cmd.Args().Get(1))
		raw, err := req(c, orgID, method, path, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func invBulk(method, path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
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
		raw, err := req(c, orgID, method, path, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func invNoArg(method, path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, method, path, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}
