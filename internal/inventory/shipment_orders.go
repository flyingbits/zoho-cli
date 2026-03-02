package inventory

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func shipmentOrdersCmd() *cli.Command {
	base := crudSubcommands("shipmentorders")
	extra := []*cli.Command{
		{
			Name:      "mark-delivered",
			Usage:     "Mark shipment order as delivered",
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
				raw, err := req(c, orgID, "POST", "/shipmentorders/"+cmd.Args().First()+"/delivered", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "shipment-orders",
		Usage:    "Shipment order operations",
		Commands: append(base, extra...),
	}
}
