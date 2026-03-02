package invoice

import (
	"github.com/urfave/cli/v3"
)

func expensesCmd() *cli.Command {
	return &cli.Command{
		Name:  "expenses",
		Usage: "Expense operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List expenses", Flags: listFlags(), Action: runList("expenses")},
			{Name: "get", Usage: "Get an expense", ArgsUsage: "<expense-id>", Action: runGet("expenses")},
			{Name: "create", Usage: "Create an expense", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("expenses")},
			{Name: "update", Usage: "Update an expense", ArgsUsage: "<expense-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("expenses")},
			{Name: "delete", Usage: "Delete an expense", ArgsUsage: "<expense-id>", Action: runDelete("expenses")},
			{Name: "list-history-comments", Usage: "List expense history and comments", ArgsUsage: "<expense-id>", Action: runGetSub("expenses", "comments")},
			{Name: "create-employee", Usage: "Create an employee", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("expenses/employees")},
			{Name: "list-employees", Usage: "List employees", Action: runList("expenses/employees")},
			{Name: "delete-employee", Usage: "Delete an employee", ArgsUsage: "<employee-id>", Action: runDelete("expenses/employees")},
		},
	}
}
