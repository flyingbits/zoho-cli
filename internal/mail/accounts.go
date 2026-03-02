package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func accountsCmd() *cli.Command {
	return &cli.Command{
		Name:  "accounts",
		Usage: "Accounts API (user mail account settings)",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "Get all accounts of a user",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/accounts", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific account details",
				ArgsUsage: "<accountId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId := cmd.Args().First()
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" {
						return internal.NewValidationError(accountRequiredMsg)
					}
					return request(c, "GET", "/accounts/"+accountId, nil)
				},
			},
			{
				Name:      "update",
				Usage:     "Update account (sequence, reply-to, forwarding, vacation, etc.; use JSON body)",
				ArgsUsage: "<accountId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, err := mustAccount(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/accounts/"+accountId, &zohttp.RequestOpts{JSON: body})
				},
			},
		},
	}
}
