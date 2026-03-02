package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func bankAccountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "bankaccounts",
		Usage: "Bank account operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List bank accounts", Flags: listFlags(), Action: runList("bankaccounts")},
			{Name: "get", Usage: "Get a bank account", ArgsUsage: "<bankaccount-id>", Action: runGet("bankaccounts")},
			{Name: "create", Usage: "Create a bank account", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("bankaccounts")},
			{Name: "update", Usage: "Update a bank account", ArgsUsage: "<bankaccount-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("bankaccounts")},
			{Name: "delete", Usage: "Delete a bank account", ArgsUsage: "<bankaccount-id>", Action: runDelete("bankaccounts")},
			{Name: "deactivate", Usage: "Deactivate account", ArgsUsage: "<bankaccount-id>", Action: runPost("bankaccounts", "inactive")},
			{Name: "activate", Usage: "Activate account", ArgsUsage: "<bankaccount-id>", Action: runPost("bankaccounts", "active")},
			{
				Name:      "import-statement",
				Usage:     "Import bank/credit card statement",
				ArgsUsage: "<bankaccount-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
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
					raw, err := req(c, orgID, "POST", "/bankaccounts/"+cmd.Args().First()+"/statement", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{Name: "get-last-statement", Usage: "Get last imported statement", ArgsUsage: "<bankaccount-id>", Action: runGetSub("bankaccounts", "statement")},
			{Name: "delete-last-statement", Usage: "Delete last imported statement", ArgsUsage: "<bankaccount-id>", Action: runDeleteSub("bankaccounts", "statement")},
		},
	}
}
