package inventory

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func purchaseOrdersCmd() *cli.Command {
	base := crudSubcommands("purchaseorders")
	extra := []*cli.Command{
		{
			Name:      "mark-issued",
			Usage:     "Mark as issued",
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
				raw, err := req(c, orgID, "POST", "/purchaseorders/"+cmd.Args().First()+"/status/issued", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-cancelled",
			Usage:     "Mark as cancelled",
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
				raw, err := req(c, orgID, "POST", "/purchaseorders/"+cmd.Args().First()+"/status/cancelled", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "purchase-orders",
		Usage:    "Purchase order operations",
		Commands: append(base, extra...),
	}
}
