package books

import (
	"github.com/urfave/cli/v3"
)

func journalsCmd() *cli.Command {
	return &cli.Command{
		Name:  "journals",
		Usage: "Journal operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "Get journal list", Flags: listFlags(), Action: runList("journals")},
			{Name: "get", Usage: "Get a journal", ArgsUsage: "<journal-id>", Action: runGet("journals")},
			{Name: "create", Usage: "Create a journal", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("journals")},
			{Name: "update", Usage: "Update a journal", ArgsUsage: "<journal-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("journals")},
			{Name: "delete", Usage: "Delete a journal", ArgsUsage: "<journal-id>", Action: runDelete("journals")},
			{Name: "mark-published", Usage: "Mark journal as published", ArgsUsage: "<journal-id>", Action: runPost("journals", "publish")},
			{Name: "add-attachment", Usage: "Add attachment", ArgsUsage: "<journal-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("journals", "attachment")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<journal-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("journals", "comments")},
		},
	}
}
