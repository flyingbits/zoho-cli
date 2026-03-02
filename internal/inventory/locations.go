package inventory

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func locationsCmd() *cli.Command {
	base := crudSubcommands("locations")
	extra := []*cli.Command{
		{
			Name:  "enable",
			Usage: "Enable locations",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "POST", "/locations/enable", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-active",
			Usage:     "Mark location as active",
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
				raw, err := req(c, orgID, "POST", "/locations/"+cmd.Args().First()+"/active", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-inactive",
			Usage:     "Mark location as inactive",
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
				raw, err := req(c, orgID, "POST", "/locations/"+cmd.Args().First()+"/inactive", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-primary",
			Usage:     "Mark as primary location",
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
				raw, err := req(c, orgID, "POST", "/locations/"+cmd.Args().First()+"/primary", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "locations",
		Usage:    "Location operations",
		Commands: append(base, extra...),
	}
}
