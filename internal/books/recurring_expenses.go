package books

import (
	"github.com/urfave/cli/v3"
)

func recurringExpensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurringexpenses",
		Usage: "Recurring expense operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List recurring expenses", Flags: listFlags(), Action: runList("recurringexpenses")},
			{Name: "get", Usage: "Get a recurring expense", ArgsUsage: "<recurringexpense-id>", Action: runGet("recurringexpenses")},
			{Name: "create", Usage: "Create a recurring expense", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("recurringexpenses")},
			{Name: "update", Usage: "Update a recurring expense", ArgsUsage: "<recurringexpense-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("recurringexpenses")},
			{Name: "delete", Usage: "Delete a recurring expense", ArgsUsage: "<recurringexpense-id>", Action: runDelete("recurringexpenses")},
			{Name: "stop", Usage: "Stop recurring expense", ArgsUsage: "<recurringexpense-id>", Action: runPost("recurringexpenses", "stop")},
			{Name: "resume", Usage: "Resume recurring expense", ArgsUsage: "<recurringexpense-id>", Action: runPost("recurringexpenses", "resume")},
			{Name: "list-child-expenses", Usage: "List child expenses created", ArgsUsage: "<recurringexpense-id>", Action: runGetSub("recurringexpenses", "childexpenses")},
			{Name: "list-history", Usage: "List recurring expense history", ArgsUsage: "<recurringexpense-id>", Action: runGetSub("recurringexpenses", "history")},
		},
	}
}
