package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func packagesCmd() *cli.Command {
	base := crudSubcommands("packages")
	extra := []*cli.Command{
		{
			Name:  "bulk-print",
			Usage: "Bulk print packages",
			Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body (package_ids)"}},
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
				raw, err := req(c, orgID, "POST", "/packages/bulk-print", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "packages",
		Usage:    "Package operations",
		Commands: append(base, extra...),
	}
}
