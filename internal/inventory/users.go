package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func usersCmd() *cli.Command {
	base := crudSubcommands("users")
	extra := []*cli.Command{
		{
			Name:  "get-current",
			Usage: "Get current user",
			Action: func(_ context.Context, cmd *cli.Command) error {
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
			},
		},
		{
			Name:      "invite",
			Usage:     "Invite a user",
			ArgsUsage: "<user-id>",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				opts := (*zohttp.RequestOpts)(nil)
				if j := cmd.String("json"); j != "" {
					var body any
					json.Unmarshal([]byte(j), &body)
					opts = &zohttp.RequestOpts{JSON: body}
				}
				raw, err := req(c, orgID, "POST", "/users/"+cmd.Args().First()+"/invite", opts)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-active",
			Usage:     "Mark user as active",
			ArgsUsage: "<id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "POST", "/users/"+cmd.Args().First()+"/active", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-inactive",
			Usage:     "Mark user as inactive",
			ArgsUsage: "<id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "POST", "/users/"+cmd.Args().First()+"/inactive", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "users",
		Usage:    "User operations",
		Commands: append(base, extra...),
	}
}
