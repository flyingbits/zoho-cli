package mail

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func threadsCmd() *cli.Command {
	return &cli.Command{
		Name:  "threads",
		Usage: "Threads API",
		Commands: []*cli.Command{
			{
				Name:      "update",
				Usage:     "Flag, move, label, mark read/unread, spam (use JSON body)",
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
					return request(c, "PUT", "/accounts/"+accountId+"/updatethread", &zohttp.RequestOpts{JSON: body})
				},
			},
		},
	}
}
