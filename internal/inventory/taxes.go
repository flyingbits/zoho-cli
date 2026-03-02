package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func taxesCmd() *cli.Command {
	base := crudSubcommands("taxes")
	extra := []*cli.Command{
		{
			Name:  "create-group",
			Usage: "Create tax group",
			Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}},
			Action: func(_ context.Context, cmd *cli.Command) error {
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
				raw, err := req(c, orgID, "POST", "/taxgroups", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "get-group",
			Usage:     "Get tax group",
			ArgsUsage: "<tax-group-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "GET", "/taxgroups/"+cmd.Args().First(), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "update-group",
			Usage:     "Update tax group",
			ArgsUsage: "<tax-group-id>",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true}},
			Action: func(_ context.Context, cmd *cli.Command) error {
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
				raw, err := req(c, orgID, "PUT", "/taxgroups/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "delete-group",
			Usage:     "Delete tax group",
			ArgsUsage: "<tax-group-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "DELETE", "/taxgroups/"+cmd.Args().First(), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "taxes",
		Usage:    "Tax operations",
		Commands: append(base, extra...),
	}
}
