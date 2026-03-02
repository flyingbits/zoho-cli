package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func currencyCmd() *cli.Command {
	return &cli.Command{
		Name:  "currency",
		Usage: "Currency operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List currencies", Flags: listFlags(), Action: runList("currencies")},
			{Name: "get", Usage: "Get a currency", ArgsUsage: "<currency-id>", Action: runGet("currencies")},
			{Name: "create", Usage: "Create a currency", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("currencies")},
			{Name: "update", Usage: "Update a currency", ArgsUsage: "<currency-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("currencies")},
			{Name: "delete", Usage: "Delete a currency", ArgsUsage: "<currency-id>", Action: runDelete("currencies")},
			{Name: "create-exchange-rate", Usage: "Create exchange rate", ArgsUsage: "<currency-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("currencies", "exchangerates")},
			{Name: "list-exchange-rates", Usage: "List exchange rates", ArgsUsage: "<currency-id>", Action: runGetSub("currencies", "exchangerates")},
			{Name: "update-exchange-rate", Usage: "Update exchange rate", ArgsUsage: "<currency-id> <rate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: currencyUpdateRate},
			{Name: "get-exchange-rate", Usage: "Get exchange rate", ArgsUsage: "<currency-id> <rate-id>", Action: currencyGetRate},
			{Name: "delete-exchange-rate", Usage: "Delete exchange rate", ArgsUsage: "<currency-id> <rate-id>", Action: currencyDeleteRate},
		},
	}
}

func currencyUpdateRate(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/currencies/"+cmd.Args().First()+"/exchangerates/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func currencyGetRate(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/currencies/"+cmd.Args().First()+"/exchangerates/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func currencyDeleteRate(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/currencies/"+cmd.Args().First()+"/exchangerates/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
