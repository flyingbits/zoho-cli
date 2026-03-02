package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func transferOrdersCmd() *cli.Command {
	base := crudSubcommands("transferorders")
	extra := []*cli.Command{
		{
			Name:      "mark-received",
			Usage:     "Mark transfer order as received",
			ArgsUsage: "<id>",
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
				opts := &zohttp.RequestOpts{}
				if j := cmd.String("json"); j != "" {
					var body any
					json.Unmarshal([]byte(j), &body)
					opts.JSON = body
				}
				raw, err := req(c, orgID, "POST", "/transferorders/"+cmd.Args().First()+"/received", opts)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "transfer-orders",
		Usage:    "Transfer order operations",
		Commands: append(base, extra...),
	}
}
