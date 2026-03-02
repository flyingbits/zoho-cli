package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "User operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List users", Flags: listFlags(), Action: runList("users")},
			{Name: "get", Usage: "Get a user", ArgsUsage: "<user-id>", Action: runGet("users")},
			{Name: "get-current", Usage: "Get current user", Action: usersGetCurrent},
			{Name: "create", Usage: "Create a user", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("users")},
			{Name: "update", Usage: "Update a user", ArgsUsage: "<user-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("users")},
			{Name: "delete", Usage: "Delete a user", ArgsUsage: "<user-id>", Action: runDelete("users")},
			{Name: "invite", Usage: "Invite a user", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: usersInvite},
			{Name: "mark-active", Usage: "Mark user as active", ArgsUsage: "<user-id>", Action: runPost("users", "active")},
			{Name: "mark-inactive", Usage: "Mark user as inactive", ArgsUsage: "<user-id>", Action: runPost("users", "inactive")},
		},
	}
}

func usersGetCurrent(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/users/me", nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func usersInvite(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	var body any
	json.Unmarshal([]byte(cmd.String("json")), &body)
	raw, err := req(c, orgID, "POST", "/users/invite", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
