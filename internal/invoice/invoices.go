package invoice

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func invoicesCmd() *cli.Command {
	return &cli.Command{
		Name:  "invoices",
		Usage: "Invoice operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List invoices", Flags: listFlags(), Action: runList("invoices")},
			{Name: "get", Usage: "Get an invoice", ArgsUsage: "<invoice-id>", Action: runGet("invoices")},
			{Name: "create", Usage: "Create an invoice", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("invoices")},
			{Name: "update", Usage: "Update an invoice", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("invoices")},
			{Name: "delete", Usage: "Delete an invoice", ArgsUsage: "<invoice-id>", Action: runDelete("invoices")},
			{Name: "update-custom-field", Usage: "Update custom field in existing invoices", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("invoices")},
			{Name: "mark-sent", Usage: "Mark invoice as sent", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "status/sent")},
			{Name: "void", Usage: "Void an invoice", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "status/void")},
			{Name: "mark-draft", Usage: "Mark as draft", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "status/draft")},
			{Name: "email", Usage: "Email an invoice", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("invoices", "email")},
			{Name: "get-email-content", Usage: "Get invoice email content", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "email")},
			{Name: "email-multiple", Usage: "Email invoices", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: invoiceEmailMultiple},
			{Name: "remind", Usage: "Remind customer", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("invoices", "paymentreminder")},
			{Name: "get-reminder-mail-content", Usage: "Get payment reminder mail content", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "paymentreminder")},
			{Name: "bulk-reminder", Usage: "Bulk invoice reminder", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: invoiceBulkReminder},
			{Name: "bulk-export", Usage: "Bulk export invoices", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: invoiceBulkExport},
			{Name: "bulk-print", Usage: "Bulk print invoices", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: invoiceBulkPrint},
			{Name: "disable-payment-reminder", Usage: "Disable payment reminder", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "paymentreminder/disable")},
			{Name: "enable-payment-reminder", Usage: "Enable payment reminder", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "paymentreminder/enable")},
			{Name: "write-off", Usage: "Write off invoice", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("invoices", "writeoff")},
			{Name: "cancel-write-off", Usage: "Cancel write off", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "writeoff/cancel")},
			{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("invoices", "address/billing")},
			{Name: "update-shipping-address", Usage: "Update shipping address", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("invoices", "address/shipping")},
			{Name: "list-templates", Usage: "List invoice templates", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "templates")},
			{Name: "update-template", Usage: "Update invoice template", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("invoices", "templates")},
			{Name: "list-payments", Usage: "List invoice payments", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "payments")},
			{Name: "list-credits-applied", Usage: "List credits applied", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "creditsapplied")},
			{Name: "apply-credits", Usage: "Apply credits", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("invoices", "credits")},
			{Name: "delete-payment", Usage: "Delete a payment", ArgsUsage: "<invoice-id> <payment-id>", Action: invoiceDeletePayment},
			{Name: "delete-applied-credit", Usage: "Delete applied credit", ArgsUsage: "<invoice-id> <credit-id>", Action: invoiceDeleteAppliedCredit},
			{Name: "add-attachment", Usage: "Add attachment to an invoice", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}, &cli.StringFlag{Name: "file", Usage: "Path to file"}}, Action: invoiceAddAttachment},
			{Name: "update-attachment-preference", Usage: "Update attachment preference", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("invoices", "attachment")},
			{Name: "get-attachment", Usage: "Get an invoice attachment", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "attachment")},
			{Name: "delete-attachment", Usage: "Delete an attachment", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "attachment/delete")},
			{Name: "delete-expense-receipt", Usage: "Delete the expense receipt", ArgsUsage: "<invoice-id>", Action: runPost("invoices", "expensereceipt/delete")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<invoice-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("invoices", "comments")},
			{Name: "list-comments", Usage: "List invoice comments and history", ArgsUsage: "<invoice-id>", Action: runGetSub("invoices", "comments")},
			{Name: "update-comment", Usage: "Update comment", ArgsUsage: "<invoice-id> <comment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: invoiceUpdateComment},
			{Name: "delete-comment", Usage: "Delete a comment", ArgsUsage: "<invoice-id> <comment-id>", Action: invoiceDeleteComment},
		},
	}
}

func invoiceEmailMultiple(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/invoices/email", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceBulkReminder(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/invoices/bulk/paymentreminder", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceBulkExport(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	opts := &zohttp.RequestOpts{}
	if j := cmd.String("json"); j != "" {
		var body any
		json.Unmarshal([]byte(j), &body)
		opts.JSON = body
	}
	raw, err := req(c, orgID, "POST", "/invoices/bulk/export", opts)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceBulkPrint(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	opts := &zohttp.RequestOpts{}
	if j := cmd.String("json"); j != "" {
		var body any
		json.Unmarshal([]byte(j), &body)
		opts.JSON = body
	}
	raw, err := req(c, orgID, "POST", "/invoices/bulk/print", opts)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceDeletePayment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/invoices/"+cmd.Args().First()+"/payments/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceDeleteAppliedCredit(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/invoices/"+cmd.Args().First()+"/creditsapplied/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceAddAttachment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	opts := &zohttp.RequestOpts{Params: orgParam(orgID)}
	if j := cmd.String("json"); j != "" {
		var body any
		json.Unmarshal([]byte(j), &body)
		opts.JSON = body
	}
	if f := cmd.String("file"); f != "" {
		data, _ := os.ReadFile(f)
		opts.Files = map[string]zohttp.FileUpload{"attachment": {Filename: filepath.Base(f), Data: data}}
	}
	raw, err := req(c, orgID, "POST", "/invoices/"+cmd.Args().First()+"/attachment", opts)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceUpdateComment(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/invoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func invoiceDeleteComment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/invoices/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
