package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func groupsCmd() *cli.Command {
	return &cli.Command{
		Name:  "groups",
		Usage: "Groups API",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a group",
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
					return request(c, "POST", "/organization/"+zoid+"/groups", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list",
				Usage:     "Get all group details",
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
					return request(c, "GET", "/organization/"+zoid+"/groups", nil)
				},
			},
			{
				Name:      "get",
				Usage:     "Get specific group details",
				ArgsUsage: "<zoid> <zgid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zgid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					return request(c, "GET", "/organization/"+zoid+"/groups/"+zgid, nil)
				},
			},
			{
				Name:      "moderation-list",
				Usage:     "Get all emails held for moderation",
				ArgsUsage: "<zoid> <zgid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zgid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					return request(c, "GET", "/organization/"+zoid+"/groups/"+zgid+"/messages", nil)
				},
			},
			{
				Name:      "moderation-get",
				Usage:     "Get moderation email content",
				ArgsUsage: "<zoid> <zgid> <messageId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zgid, msgID := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zgid == "" || msgID == "" {
						return internal.NewValidationError("zgid and messageId required")
					}
					return request(c, "GET", "/organization/"+zoid+"/groups/"+zgid+"/messages/"+msgID, nil)
				},
			},
			{
				Name:      "moderate",
				Usage:     "Moderate emails in a group",
				ArgsUsage: "<zoid> <zgid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zgid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/groups/"+zgid+"/messages", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update",
				Usage:     "Update group (name, members, settings, etc.; use JSON body)",
				ArgsUsage: "<zoid> <zgid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zgid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/groups/"+zgid, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete group",
				ArgsUsage: "<zoid> <zgid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, zgid := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					return request(c, "DELETE", "/organization/"+zoid+"/groups/"+zgid, nil)
				},
			},
		},
	}
}
