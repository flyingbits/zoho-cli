package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func customModulesCmd() *cli.Command {
	return &cli.Command{
		Name:  "custommodules",
		Usage: "Custom module operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "Get record list of a custom module",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module-name", Required: true, Usage: "Custom module API name"},
					&cli.StringFlag{Name: "page", Usage: "Page"},
					&cli.StringFlag{Name: "per-page", Usage: "Per page"},
					&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
					&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := mergeParams(orgID, map[string]string{"module_name": cmd.String("module-name")})
					if v := cmd.String("page"); v != "" {
						params["page"] = v
					}
					if v := cmd.String("per-page"); v != "" {
						params["per_page"] = v
					}
					if v := cmd.String("sort-column"); v != "" {
						params["sort_column"] = v
					}
					if v := cmd.String("sort-order"); v != "" {
						params["sort_order"] = v
					}
					raw, err := req(c, orgID, "GET", "/custommodules", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get individual record details",
				ArgsUsage: "<record-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "module-name", Required: true, Usage: "Custom module API name"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := mergeParams(orgID, map[string]string{"module_name": cmd.String("module-name")})
					raw, err := req(c, orgID, "GET", "/custommodules/"+cmd.Args().First(), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create custom module",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
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
					raw, err := req(c, orgID, "POST", "/custommodules", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update custom module",
				ArgsUsage: "<module-id>",
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
					raw, err := req(c, orgID, "PUT", "/custommodules/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete custom module or delete individual record",
				ArgsUsage: "<id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "module-name", Usage: "Custom module API name (for deleting record)"}},
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
					if m := cmd.String("module-name"); m != "" {
						opts.Params = mergeParams(orgID, map[string]string{"module_name": m})
					} else {
						opts.Params = orgParam(orgID)
					}
					raw, err := req(c, orgID, "DELETE", "/custommodules/"+cmd.Args().First(), opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "bulk-update",
				Usage: "Bulk update custom module",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "module-name", Required: true, Usage: "Custom module API name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
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
					params := mergeParams(orgID, map[string]string{"module_name": cmd.String("module-name")})
					raw, err := req(c, orgID, "PUT", "/custommodules", &zohttp.RequestOpts{JSON: body, Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
