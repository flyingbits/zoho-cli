package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func customerPaymentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "customer-payments",
		Usage: "Customer payment operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List customer payments", Flags: listFlags(), Action: runList("customer-payments")},
			{Name: "get", Usage: "Retrieve a payment", ArgsUsage: "<payment-id>", Action: runGet("customer-payments")},
			{Name: "create", Usage: "Create a payment", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("customer-payments")},
			{Name: "update", Usage: "Update a payment", ArgsUsage: "<payment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("customer-payments")},
			{Name: "delete", Usage: "Delete a payment", ArgsUsage: "<payment-id>", Action: runDelete("customer-payments")},
			{Name: "update-custom-field", Usage: "Update custom field", ArgsUsage: "<payment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("customer-payments")},
			{Name: "refund", Usage: "Refund excess payment", ArgsUsage: "<payment-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("customer-payments", "refunds")},
			{Name: "list-refunds", Usage: "List refunds", ArgsUsage: "<payment-id>", Action: runGetSub("customer-payments", "refunds")},
			{Name: "get-refund", Usage: "Details of a refund", ArgsUsage: "<payment-id> <refund-id>", Action: paymentGetRefund},
			{Name: "update-refund", Usage: "Update a refund", ArgsUsage: "<payment-id> <refund-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: paymentUpdateRefund},
			{Name: "delete-refund", Usage: "Delete a refund", ArgsUsage: "<payment-id> <refund-id>", Action: paymentDeleteRefund},
		},
	}
}

func paymentGetRefund(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/customer-payments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func paymentUpdateRefund(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/customer-payments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func paymentDeleteRefund(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/customer-payments/"+cmd.Args().First()+"/refunds/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
