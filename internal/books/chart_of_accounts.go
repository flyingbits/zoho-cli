package books

import (
	"github.com/urfave/cli/v3"
)

func chartOfAccountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "chartofaccounts",
		Usage: "Chart of accounts operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List chart of accounts", Flags: listFlags(), Action: runList("chartofaccounts")},
			{Name: "get", Usage: "Get an account", ArgsUsage: "<account-id>", Action: runGet("chartofaccounts")},
			{Name: "create", Usage: "Create an account", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("chartofaccounts")},
			{Name: "update", Usage: "Update an account", ArgsUsage: "<account-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("chartofaccounts")},
			{Name: "delete", Usage: "Delete an account", ArgsUsage: "<account-id>", Action: runDelete("chartofaccounts")},
			{Name: "mark-active", Usage: "Mark account as active", ArgsUsage: "<account-id>", Action: runPost("chartofaccounts", "active")},
			{Name: "mark-inactive", Usage: "Mark account as inactive", ArgsUsage: "<account-id>", Action: runPost("chartofaccounts", "inactive")},
			{Name: "list-transactions", Usage: "List transactions for account", ArgsUsage: "<account-id>", Action: runGetSub("chartofaccounts", "transactions")},
		},
	}
}
