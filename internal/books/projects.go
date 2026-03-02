package books

import (
	"github.com/urfave/cli/v3"
)

func projectsCmd() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Books project operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List projects", Flags: listFlags(), Action: runList("projects")},
			{Name: "get", Usage: "Get a project", ArgsUsage: "<project-id>", Action: runGet("projects")},
			{Name: "create", Usage: "Create a project", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("projects")},
			{Name: "update", Usage: "Update a project", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("projects")},
			{Name: "delete", Usage: "Delete a project", ArgsUsage: "<project-id>", Action: runDelete("projects")},
			{Name: "activate", Usage: "Activate project", ArgsUsage: "<project-id>", Action: runPost("projects", "active")},
			{Name: "inactivate", Usage: "Inactivate project", ArgsUsage: "<project-id>", Action: runPost("projects", "inactive")},
			{Name: "clone", Usage: "Clone project", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("projects", "clone")},
			{Name: "list-users", Usage: "List users", ArgsUsage: "<project-id>", Action: runGetSub("projects", "users")},
			{Name: "invite-user", Usage: "Invite user", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("projects", "users")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<project-id>", Action: runGetSub("projects", "comments")},
			{Name: "add-comment", Usage: "Post comment", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("projects", "comments")},
			{Name: "list-invoices", Usage: "List project invoices", ArgsUsage: "<project-id>", Action: runGetSub("projects", "invoices")},
		},
	}
}
