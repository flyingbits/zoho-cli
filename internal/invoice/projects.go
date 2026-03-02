package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func projectsCmd() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Invoice project operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List projects", Flags: listFlags(), Action: runList("projects")},
			{Name: "get", Usage: "Get a project", ArgsUsage: "<project-id>", Action: runGet("projects")},
			{Name: "create", Usage: "Create a project", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("projects")},
			{Name: "update", Usage: "Update a project", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("projects")},
			{Name: "delete", Usage: "Delete a project", ArgsUsage: "<project-id>", Action: runDelete("projects")},
			{Name: "activate", Usage: "Activate a project", ArgsUsage: "<project-id>", Action: runPost("projects", "active")},
			{Name: "deactivate", Usage: "Deactivate a project", ArgsUsage: "<project-id>", Action: runPost("projects", "inactive")},
			{Name: "clone", Usage: "Clone a project", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("projects", "clone")},
			{Name: "assign-users", Usage: "Assign users", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("projects", "users")},
			{Name: "list-users", Usage: "List users", ArgsUsage: "<project-id>", Action: runGetSub("projects", "users")},
			{Name: "invite-user", Usage: "Invite user", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("projects", "users/invite")},
			{Name: "update-user", Usage: "Update user", ArgsUsage: "<project-id> <user-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: projectUpdateUser},
			{Name: "get-user", Usage: "Get a user", ArgsUsage: "<project-id> <user-id>", Action: projectGetUser},
			{Name: "delete-user", Usage: "Delete user", ArgsUsage: "<project-id> <user-id>", Action: projectDeleteUser},
			{Name: "add-comment", Usage: "Post comment", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("projects", "comments")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<project-id>", Action: runGetSub("projects", "comments")},
			{Name: "delete-comment", Usage: "Delete comment", ArgsUsage: "<project-id> <comment-id>", Action: projectDeleteComment},
			{Name: "list-invoices", Usage: "List invoices", ArgsUsage: "<project-id>", Action: runGetSub("projects", "invoices")},
		},
	}
}

func projectUpdateUser(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func projectDeleteUser(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func projectDeleteComment(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/projects/"+cmd.Args().First()+"/comments/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func projectGetUser(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/projects/"+cmd.Args().First()+"/users/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
