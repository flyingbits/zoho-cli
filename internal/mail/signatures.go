package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func signaturesCmd() *cli.Command {
	return &cli.Command{
		Name:  "signatures",
		Usage: "Signatures API",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add user signature (user context)",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/accounts/signature", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "add-admin",
				Usage:     "Add admin-added user signature",
				ArgsUsage: "<zoid> <zuid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zuid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zuid == "" {
						return internal.NewValidationError("zuid required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/accounts/"+zuid+"/signature", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "get",
				Usage: "Get user signature (user context)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/accounts/signature", nil)
				},
			},
			{
				Name:      "get-admin",
				Usage:     "Get admin-added user signature",
				ArgsUsage: "<zoid> <zuid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zuid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zuid == "" {
						return internal.NewValidationError("zuid required")
					}
					return request(c, "GET", "/organization/"+zoid+"/accounts/"+zuid+"/signature", nil)
				},
			},
			{
				Name:  "update",
				Usage: "Update user signature (user context)",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/accounts/signature", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update-admin",
				Usage:     "Update admin-added user signature",
				ArgsUsage: "<zoid> <zuid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zuid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zuid == "" {
						return internal.NewValidationError("zuid required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/accounts/"+zuid+"/signature", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "delete",
				Usage: "Delete user signature (user context)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "DELETE", "/accounts/signature", nil)
				},
			},
			{
				Name:      "delete-admin",
				Usage:     "Delete admin-added user signature",
				ArgsUsage: "<zoid> <zuid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zuid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zuid == "" {
						return internal.NewValidationError("zuid required")
					}
					return request(c, "DELETE", "/organization/"+zoid+"/accounts/"+zuid+"/signature", nil)
				},
			},
		},
	}
}
