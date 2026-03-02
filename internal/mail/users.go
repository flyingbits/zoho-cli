package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Users API (organization accounts)",
		Commands: []*cli.Command{
			{Name: "add", Usage: "Add user to organization",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/accounts/", &zohttp.RequestOpts{JSON: body})
				}},
			{Name: "list", Usage: "Fetch all org users details", ArgsUsage: "[zoid]",
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
					return request(c, "GET", "/organization/"+zoid+"/accounts/", nil)
				}},
			{Name: "get", Usage: "Fetch single user (zuid or email)", ArgsUsage: "<zoid> <zuid|email>",
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
						return internal.NewValidationError("zuid or emailAddress required")
					}
					return request(c, "GET", "/organization/"+zoid+"/accounts/"+zuid, nil)
				}},
			{Name: "change-role", Usage: "Change user role", ArgsUsage: "<zoid>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/accounts", &zohttp.RequestOpts{JSON: body})
				}},
			{Name: "reset-password", Usage: "Reset user password", ArgsUsage: "<zoid> <zuid>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
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
					return request(c, "PUT", "/organization/"+zoid+"/accounts/"+zuid, &zohttp.RequestOpts{JSON: body})
				}},
			{Name: "update", Usage: "Update user (alias, enable/disable, etc.; JSON body)", ArgsUsage: "<zoid> <zuid|accountId>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, id := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if id == "" {
						return internal.NewValidationError("zuid or accountId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/accounts/"+id, &zohttp.RequestOpts{JSON: body})
				}},
			{Name: "delete", Usage: "Delete user account", ArgsUsage: "<zoid>",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "DELETE", "/organization/"+zoid+"/accounts", &zohttp.RequestOpts{JSON: body})
				}},
		},
	}
}
