package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func salesReturnsCmd() *cli.Command {
	base := crudSubcommands("salesreturns")
	extra := []*cli.Command{
		{
			Name:      "create-receive",
			Usage:     "Create sales return receive",
			ArgsUsage: "<sales-return-id>",
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
				raw, err := req(c, orgID, "POST", "/salesreturns/"+cmd.Args().First()+"/receive", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "delete-receive",
			Usage:     "Delete sales return receive",
			ArgsUsage: "<sales-return-id> <receive-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "DELETE", "/salesreturns/"+cmd.Args().Get(0)+"/receive/"+cmd.Args().Get(1), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "sales-returns",
		Usage:    "Sales return operations",
		Commands: append(base, extra...),
	}
}
