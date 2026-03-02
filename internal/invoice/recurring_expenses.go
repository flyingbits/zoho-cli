package invoice

import (
	"github.com/urfave/cli/v3"
)

func recurringExpensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurring-expenses",
		Usage: "Recurring expense operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List recurring expenses", Flags: listFlags(), Action: runList("recurring-expenses")},
			{Name: "get", Usage: "Get a recurring expense", ArgsUsage: "<recurring-expense-id>", Action: runGet("recurring-expenses")},
			{Name: "create", Usage: "Create a recurring expense", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("recurring-expenses")},
			{Name: "update", Usage: "Update a recurring expense", ArgsUsage: "<recurring-expense-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("recurring-expenses")},
			{Name: "delete", Usage: "Delete a recurring expense", ArgsUsage: "<recurring-expense-id>", Action: runDelete("recurring-expenses")},
			{Name: "stop", Usage: "Stop a recurring expense", ArgsUsage: "<recurring-expense-id>", Action: runPost("recurring-expenses", "status/stop")},
			{Name: "resume", Usage: "Resume a recurring expense", ArgsUsage: "<recurring-expense-id>", Action: runPost("recurring-expenses", "status/resume")},
			{Name: "list-child-expenses", Usage: "List child expenses created", ArgsUsage: "<recurring-expense-id>", Action: runGetSub("recurring-expenses", "child_expenses")},
			{Name: "list-history", Usage: "List recurring expense history", ArgsUsage: "<recurring-expense-id>", Action: runGetSub("recurring-expenses", "history")},
		},
	}
}
