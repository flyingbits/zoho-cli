package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func estimatesCmd() *cli.Command {
	return &cli.Command{
		Name:  "estimates",
		Usage: "Estimate operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List estimates", Flags: listFlags(), Action: runList("estimates")},
			{Name: "get", Usage: "Get an estimate", ArgsUsage: "<estimate-id>", Action: runGet("estimates")},
			{Name: "create", Usage: "Create an estimate", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("estimates")},
			{Name: "update", Usage: "Update an estimate", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("estimates")},
			{Name: "delete", Usage: "Delete an estimate", ArgsUsage: "<estimate-id>", Action: runDelete("estimates")},
			{Name: "update-custom-field", Usage: "Update custom field in existing estimates", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("estimates")},
			{Name: "mark-sent", Usage: "Mark estimate as sent", ArgsUsage: "<estimate-id>", Action: runPost("estimates", "status/sent")},
			{Name: "mark-accepted", Usage: "Mark estimate as accepted", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("estimates", "status/accepted")},
			{Name: "mark-declined", Usage: "Mark estimate as declined", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("estimates", "status/declined")},
			{Name: "email", Usage: "Email an estimate", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("estimates", "email")},
			{Name: "get-email-content", Usage: "Get estimate email content", ArgsUsage: "<estimate-id>", Action: runGetSub("estimates", "email")},
			{Name: "email-multiple", Usage: "Email multiple estimates", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: estimateEmailMultiple},
			{Name: "bulk-export", Usage: "Bulk export estimates", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: estimateBulkExport},
			{Name: "bulk-print", Usage: "Bulk print estimates", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: estimateBulkPrint},
			{Name: "update-billing-address", Usage: "Update billing address", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("estimates", "address/billing")},
			{Name: "update-shipping-address", Usage: "Update shipping address", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("estimates", "address/shipping")},
			{Name: "list-templates", Usage: "List estimate templates", ArgsUsage: "<estimate-id>", Action: runGetSub("estimates", "templates")},
			{Name: "update-template", Usage: "Update estimate template", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runPostJSON("estimates", "templates")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("estimates", "comments")},
			{Name: "list-comments", Usage: "List estimate comments and history", ArgsUsage: "<estimate-id>", Action: runGetSub("estimates", "comments")},
			{Name: "update-comment", Usage: "Update comment", ArgsUsage: "<estimate-id> <comment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: estimateUpdateComment},
			{Name: "delete-comment", Usage: "Delete a comment", ArgsUsage: "<estimate-id> <comment-id>", Action: estimateDeleteComment},
		},
	}
}

func estimateEmailMultiple(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/estimates/email", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func estimateBulkExport(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/estimates/bulk/export", opts)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func estimateBulkPrint(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/estimates/bulk/print", opts)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func estimateUpdateComment(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/estimates/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func estimateDeleteComment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/estimates/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
