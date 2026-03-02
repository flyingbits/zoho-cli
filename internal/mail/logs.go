package mail

import (
	"context"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func logsCmd() *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Logs API",
		Commands: []*cli.Command{
			{
				Name:      "login-history",
				Usage:     "Get login history (admin)",
				ArgsUsage: "<zoid>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "zuid", Usage: "User ID filter"},
					&cli.StringFlag{Name: "start", Usage: "Start index"},
					&cli.StringFlag{Name: "limit", Usage: "Limit"},
				},
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
					params := map[string]string{}
					if v := cmd.String("zuid"); v != "" {
						params["zuid"] = v
					}
					if v := cmd.String("start"); v != "" {
						params["start"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					return request(c, "GET", "/organization/"+zoid+"/accounts/reports/loginHistory", &zohttp.RequestOpts{Params: params})
				},
			},
			{
				Name:      "audit",
				Usage:     "Get audit records",
				ArgsUsage: "<zoid>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "start", Usage: "Start index"},
					&cli.StringFlag{Name: "limit", Usage: "Limit"},
				},
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
					params := map[string]string{}
					if v := cmd.String("start"); v != "" {
						params["start"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					return request(c, "GET", "/organization/"+zoid+"/activity", &zohttp.RequestOpts{Params: params})
				},
			},
			{
				Name:      "smtp",
				Usage:     "Get SMTP logs",
				ArgsUsage: "<zoid>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "start", Usage: "Start index"},
					&cli.StringFlag{Name: "limit", Usage: "Limit"},
				},
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
					params := map[string]string{}
					if v := cmd.String("start"); v != "" {
						params["start"] = v
					}
					if v := cmd.String("limit"); v != "" {
						params["limit"] = v
					}
					return request(c, "GET", "/organization/"+zoid+"/smtplogs", &zohttp.RequestOpts{Params: params})
				},
			},
		},
	}
}
