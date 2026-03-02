package books

import (
	"github.com/urfave/cli/v3"
)

func bankTransactionsCmd() *cli.Command {
	return &cli.Command{
		Name:  "banktransactions",
		Usage: "Bank transaction operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "Get transactions list", ArgsUsage: "<bankaccount-id>", Action: runGetSub("bankaccounts", "transactions")},
			{Name: "get", Usage: "Get a transaction", ArgsUsage: "<transaction-id>", Action: runGet("banktransactions")},
			{Name: "create", Usage: "Create a transaction", Flags: []cli.Flag{&cli.StringFlag{Name: "account-id", Required: true, Usage: "Bank account ID"}, &cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("banktransactions")},
			{Name: "update", Usage: "Update a transaction", ArgsUsage: "<transaction-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("banktransactions")},
			{Name: "delete", Usage: "Delete a transaction", ArgsUsage: "<transaction-id>", Action: runDelete("banktransactions")},
			{Name: "match", Usage: "Match a transaction", ArgsUsage: "<transaction-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("banktransactions", "match")},
			{Name: "unmatch", Usage: "Unmatch a matched transaction", ArgsUsage: "<transaction-id>", Action: runPost("banktransactions", "unmatch")},
			{Name: "exclude", Usage: "Exclude a transaction", ArgsUsage: "<transaction-id>", Action: runPost("banktransactions", "exclude")},
			{Name: "restore", Usage: "Restore a transaction", ArgsUsage: "<transaction-id>", Action: runPost("banktransactions", "restore")},
		},
	}
}
