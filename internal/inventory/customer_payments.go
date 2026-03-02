package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func customerPaymentsCmd() *cli.Command {
	base := crudSubcommands("customerpayments")
	extra := []*cli.Command{
		{
			Name:      "update-custom-field",
			Usage:     "Update custom field",
			ArgsUsage: "<payment-id>",
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
				raw, err := req(c, orgID, "PUT", "/customerpayments/"+cmd.Args().First()+"/customfield", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "customer-payments",
		Usage:    "Customer payment operations",
		Commands: append(base, extra...),
	}
}
