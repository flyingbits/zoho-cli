package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func organizationCmd() *cli.Command {
	return &cli.Command{
		Name:  "organization",
		Usage: "Organization API",
		Commands: []*cli.Command{
			{
				Name:  "add-child",
				Usage: "Add child organization",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "allowed-ips-add",
				Usage:     "Add allowed IPs",
				ArgsUsage: "<zoid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, _ = mustOrg(cmd)
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/allowedIps", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "get",
				Usage:     "Get organization details",
				ArgsUsage: "<zoid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid, nil)
				},
			},
			{
				Name:      "storage",
				Usage:     "Get org subscription/storage details",
				ArgsUsage: "<zoid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/storage", nil)
				},
			},
			{
				Name:      "storage-user",
				Usage:     "Get user storage details",
				ArgsUsage: "<zoid> <zuid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().Get(0)
					zuid := cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zuid == "" {
						return internal.NewValidationError("zuid required")
					}
					return request(c, "GET", "/organization/"+zoid+"/storage/"+zuid, nil)
				},
			},
			{
				Name:      "spam-list",
				Usage:     "Get org spam listing",
				ArgsUsage: "<zoid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/antispam/data", nil)
				},
			},
			{
				Name:      "allowed-ips-list",
				Usage:     "Get allowed IPs list",
				ArgsUsage: "<zoid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/allowedIps", nil)
				},
			},
			{
				Name:      "storage-update",
				Usage:     "Update user storage",
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
					return request(c, "PUT", "/organization/"+zoid+"/storage/"+zuid, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "spam-type-update",
				Usage:     "Update org spam process type",
				ArgsUsage: "<zoid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "spam-update",
				Usage:     "Update org spam listing",
				ArgsUsage: "<zoid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/antispam/data", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "spam-remove",
				Usage:     "Remove org spam listing",
				ArgsUsage: "<zoid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body (optional)"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					opts := (*zohttp.RequestOpts)(nil)
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts = &zohttp.RequestOpts{JSON: body}
					}
					return request(c, "DELETE", "/organization/"+zoid+"/antispam/data", opts)
				},
			},
			{
				Name:      "allowed-ips-delete",
				Usage:     "Delete allowed IPs",
				ArgsUsage: "<zoid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "DELETE", "/organization/"+zoid+"/allowedIps", &zohttp.RequestOpts{JSON: body})
				},
			},
		},
	}
}
