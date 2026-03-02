package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func compositeItemsCmd() *cli.Command {
	base := crudSubcommands("compositeitems")
	extra := []*cli.Command{
		{
			Name:      "mark-active",
			Usage:     "Mark composite item as active",
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
				raw, err := req(c, orgID, "POST", "/compositeitems/"+cmd.Args().First()+"/active", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-inactive",
			Usage:     "Mark composite item as inactive",
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
				raw, err := req(c, orgID, "POST", "/compositeitems/"+cmd.Args().First()+"/inactive", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "create-assemblies",
			Usage:     "Create assemblies",
			ArgsUsage: "<composite-item-id>",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
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
				raw, err := req(c, orgID, "POST", "/compositeitems/"+cmd.Args().First()+"/assemblies", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "list-assemblies",
			Usage:     "Retrieve assemblies",
			ArgsUsage: "<composite-item-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "GET", "/compositeitems/"+cmd.Args().First()+"/assemblies", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "assemblies-history",
			Usage:     "Assemblies history",
			ArgsUsage: "<composite-item-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "GET", "/compositeitems/"+cmd.Args().First()+"/assemblies/history", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "delete-assemblies",
			Usage:     "Delete assemblies",
			ArgsUsage: "<composite-item-id>",
			Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body (assembly IDs)"}},
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
				raw, err := req(c, orgID, "DELETE", "/compositeitems/"+cmd.Args().First()+"/assemblies", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "composite-items",
		Usage:    "Composite item operations",
		Commands: append(base, extra...),
	}
}
