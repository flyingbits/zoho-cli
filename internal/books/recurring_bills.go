package books

import (
	"github.com/urfave/cli/v3"
)

func recurringBillsCmd() *cli.Command {
	return &cli.Command{
		Name:  "recurringbills",
		Usage: "Recurring bill operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List recurring bills", Flags: listFlags(), Action: runList("recurringbills")},
			{Name: "get", Usage: "Get a recurring bill", ArgsUsage: "<recurringbill-id>", Action: runGet("recurringbills")},
			{Name: "create", Usage: "Create a recurring bill", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("recurringbills")},
			{Name: "update", Usage: "Update a recurring bill", ArgsUsage: "<recurringbill-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("recurringbills")},
			{Name: "delete", Usage: "Delete a recurring bill", ArgsUsage: "<recurringbill-id>", Action: runDelete("recurringbills")},
			{Name: "stop", Usage: "Stop recurring bill", ArgsUsage: "<recurringbill-id>", Action: runPost("recurringbills", "stop")},
			{Name: "resume", Usage: "Resume recurring bill", ArgsUsage: "<recurringbill-id>", Action: runPost("recurringbills", "resume")},
			{Name: "list-history", Usage: "List recurring bill history", ArgsUsage: "<recurringbill-id>", Action: runGetSub("recurringbills", "history")},
		},
	}
}
