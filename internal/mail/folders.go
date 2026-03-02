package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func foldersCmd() *cli.Command {
	return &cli.Command{
		Name:  "folders",
		Usage: "Folders API",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a new folder",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"},
				},
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
					return request(c, "POST", "/accounts/"+accountId+"/folders", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list",
				Usage:     "Get all folders",
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
					return request(c, "GET", "/accounts/"+accountId+"/folders", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get specific folder",
				ArgsUsage: "<accountId> <folderId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, folderId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || folderId == "" {
						return internal.NewValidationError("accountId and folderId required")
					}
					return request(c, "GET", "/accounts/"+accountId+"/folders/"+folderId, nil)
				},
			},
			{
				Name:      "update",
				Usage:     "Rename, move, IMAP, mark read, empty (use JSON body)",
				ArgsUsage: "<accountId> <folderId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, folderId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || folderId == "" {
						return internal.NewValidationError("accountId and folderId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/accounts/"+accountId+"/folders/"+folderId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete folder",
				ArgsUsage: "<accountId> <folderId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, folderId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || folderId == "" {
						return internal.NewValidationError("accountId and folderId required")
					}
					return request(c, "DELETE", "/accounts/"+accountId+"/folders/"+folderId, nil)
				},
			},
		},
	}
}
