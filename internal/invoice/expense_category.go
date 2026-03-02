package invoice

import (
	"github.com/urfave/cli/v3"
)

func expenseCategoryCmd() *cli.Command {
	return &cli.Command{
		Name:  "expense-category",
		Usage: "Expense category operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List expense categories", Flags: listFlags(), 			Action: runList("expense-category")},
			{Name: "get", Usage: "Get an expense category", ArgsUsage: "<category-id>", Action: runGet("expense-category")},
			{Name: "create", Usage: "Create an expense category", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("expense-category")},
			{Name: "update", Usage: "Update an expense category", ArgsUsage: "<category-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("expense-category")},
			{Name: "delete", Usage: "Delete an expense category", ArgsUsage: "<category-id>", Action: runDelete("expense-category")},
			{Name: "mark-active", Usage: "Mark expense category as active", ArgsUsage: "<category-id>", Action: runPost("expense-category", "active")},
			{Name: "mark-inactive", Usage: "Mark expense category as inactive", ArgsUsage: "<category-id>", Action: runPost("expense-category", "inactive")},
		},
	}
}
