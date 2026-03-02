package inventory

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func itemGroupsCmd() *cli.Command {
	base := crudSubcommands("itemgroups")
	extra := []*cli.Command{
		{
			Name:      "mark-active",
			Usage:     "Mark item group as active",
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
				raw, err := req(c, orgID, "POST", "/itemgroups/"+cmd.Args().First()+"/active", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-inactive",
			Usage:     "Mark item group as inactive",
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
				raw, err := req(c, orgID, "POST", "/itemgroups/"+cmd.Args().First()+"/inactive", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "item-groups",
		Usage:    "Item group operations",
		Commands: append(base, extra...),
	}
}
