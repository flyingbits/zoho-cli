package books

import (
	"github.com/urfave/cli/v3"
)

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Books user operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List users", Flags: listFlags(), Action: runList("users")},
			{Name: "get", Usage: "Get a user", ArgsUsage: "<user-id>", Action: runGet("users")},
			{Name: "get-current", Usage: "Get current user", Action: runGetNoID("users/me")},
			{Name: "create", Usage: "Create a user", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("users")},
			{Name: "update", Usage: "Update a user", ArgsUsage: "<user-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("users")},
			{Name: "delete", Usage: "Delete a user", ArgsUsage: "<user-id>", Action: runDelete("users")},
			{Name: "invite", Usage: "Invite a user", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("users/invite")},
			{Name: "mark-active", Usage: "Mark user as active", ArgsUsage: "<user-id>", Action: runPost("users", "active")},
			{Name: "mark-inactive", Usage: "Mark user as inactive", ArgsUsage: "<user-id>", Action: runPost("users", "inactive")},
		},
	}
}
