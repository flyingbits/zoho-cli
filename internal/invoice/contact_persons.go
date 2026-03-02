package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func contactPersonsCmd() *cli.Command {
	return &cli.Command{
		Name:  "contact-persons",
		Usage: "Contact person operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List contact persons for a contact",
				ArgsUsage: "<contact-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					raw, err := req(c, orgID, "GET", "/contacts/"+cmd.Args().First()+"/contact_persons", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a contact person",
				ArgsUsage: "<contact-person-id>",
				Action:    runGet("contact_persons"),
			},
			{
				Name:      "create",
				Usage:     "Create a contact person",
				ArgsUsage: "<contact-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().First()+"/contact_persons", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a contact person",
				ArgsUsage: "<contact-person-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action:    runUpdate("contact_persons"),
			},
			{
				Name:      "delete",
				Usage:     "Delete a contact person",
				ArgsUsage: "<contact-person-id>",
				Action:    runDelete("contact_persons"),
			},
			{
				Name:      "mark-primary",
				Usage:     "Mark as primary contact person",
				ArgsUsage: "<contact-person-id>",
				Action:    runPost("contact_persons", "primary"),
			},
		},
	}
}
