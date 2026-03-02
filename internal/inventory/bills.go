package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func billsCmd() *cli.Command {
	base := crudSubcommands("bills")
	extra := []*cli.Command{
		{
			Name:      "update-custom-field",
			Usage:     "Update custom field",
			ArgsUsage: "<bill-id>",
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
				raw, err := req(c, orgID, "PUT", "/bills/"+cmd.Args().First()+"/customfield", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-open",
			Usage:     "Mark as open",
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
				raw, err := req(c, orgID, "POST", "/bills/"+cmd.Args().First()+"/open", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-void",
			Usage:     "Mark as void",
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
				raw, err := req(c, orgID, "POST", "/bills/"+cmd.Args().First()+"/void", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "bills",
		Usage:    "Bill operations",
		Commands: append(base, extra...),
	}
}
