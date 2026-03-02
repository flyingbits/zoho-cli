package books

import (
	"github.com/urfave/cli/v3"
)

func bankRulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "bankrules",
		Usage: "Bank rule operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "Get rules list", Flags: listFlags(), Action: runList("bankrules")},
			{Name: "get", Usage: "Get a rule", ArgsUsage: "<rule-id>", Action: runGet("bankrules")},
			{Name: "create", Usage: "Create a rule", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("bankrules")},
			{Name: "update", Usage: "Update a rule", ArgsUsage: "<rule-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("bankrules")},
			{Name: "delete", Usage: "Delete a rule", ArgsUsage: "<rule-id>", Action: runDelete("bankrules")},
		},
	}
}
