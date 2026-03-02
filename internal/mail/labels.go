package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func labelsCmd() *cli.Command {
	return &cli.Command{
		Name:  "labels",
		Usage: "Labels API",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a new label",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
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
					return request(c, "POST", "/accounts/"+accountId+"/labels", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list",
				Usage:     "Get all label details",
				ArgsUsage: "[accountId]",
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
					return request(c, "GET", "/accounts/"+accountId+"/labels", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific label details",
				ArgsUsage: "<accountId> <labelId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, labelId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || labelId == "" {
						return internal.NewValidationError("accountId and labelId required")
					}
					return request(c, "GET", "/accounts/"+accountId+"/labels/"+labelId, nil)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a label",
				ArgsUsage: "<accountId> <labelId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, labelId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || labelId == "" {
						return internal.NewValidationError("accountId and labelId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/accounts/"+accountId+"/labels/"+labelId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a label",
				ArgsUsage: "<accountId> <labelId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, labelId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || labelId == "" {
						return internal.NewValidationError("accountId and labelId required")
					}
					return request(c, "DELETE", "/accounts/"+accountId+"/labels/"+labelId, nil)
				},
			},
		},
	}
}
