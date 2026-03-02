package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func domainsCmd() *cli.Command {
	return &cli.Command{
		Name:  "domains",
		Usage: "Domain API",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a domain to an organization",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
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
					return request(c, "POST", "/organization/"+zoid+"/domains", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list",
				Usage:     "Fetch all domain details",
				ArgsUsage: "[zoid]",
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
					return request(c, "GET", "/organization/"+zoid+"/domains", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Fetch a specific domain details",
				ArgsUsage: "<zoid> <domainname>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().Get(0)
					domain := cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if domain == "" {
						return internal.NewValidationError("domainname required")
					}
					return request(c, "GET", "/organization/"+zoid+"/domains/"+domain, nil)
				},
			},
			{
				Name:      "update",
				Usage:     "Update domain (verify, set primary, DKIM, etc.; use JSON body)",
				ArgsUsage: "<zoid> <domainname>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, domain := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if domain == "" {
						return internal.NewValidationError("domainname required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/domains/"+domain, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a domain from the organization",
				ArgsUsage: "<zoid> <domainname>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, domain := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if domain == "" {
						return internal.NewValidationError("domainname required")
					}
					return request(c, "DELETE", "/organization/"+zoid+"/domains/"+domain, nil)
				},
			},
		},
	}
}
