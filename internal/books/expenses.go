package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func expensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "expenses",
		Usage: "Expense operations (Books)",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List expenses", Flags: listFlags(), Action: runList("expenses")},
			{Name: "get", Usage: "Get an expense", ArgsUsage: "<expense-id>", Action: runGet("expenses")},
			{Name: "create", Usage: "Create an expense", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("expenses")},
			{Name: "update", Usage: "Update an expense", ArgsUsage: "<expense-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("expenses")},
			{Name: "delete", Usage: "Delete an expense", ArgsUsage: "<expense-id>", Action: runDelete("expenses")},
			{Name: "list-history", Usage: "List expense history and comments", ArgsUsage: "<expense-id>", Action: runGetSub("expenses", "comments")},
			{
				Name:      "add-receipt",
				Usage:     "Add receipt to an expense",
				ArgsUsage: "<expense-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
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
					raw, err := req(c, orgID, "POST", "/expenses/"+cmd.Args().First()+"/receipt", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{Name: "get-receipt", Usage: "Get expense receipt", ArgsUsage: "<expense-id>", Action: runGetSub("expenses", "receipt")},
			{Name: "delete-receipt", Usage: "Delete receipt", ArgsUsage: "<expense-id>", Action: runDeleteSub("expenses", "receipt")},
		},
	}
}
