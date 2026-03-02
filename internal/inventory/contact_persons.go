package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func contactPersonsCmd() *cli.Command {
	path := "contact_persons"
	base := []*cli.Command{
		{
			Name:  "list",
			Usage: "List contact persons",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "contact-id", Required: true, Usage: "Contact ID"},
				&cli.StringFlag{Name: "page", Usage: "Page"},
				&cli.StringFlag{Name: "per-page", Usage: "Per page"},
			},
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				params := mergeParams(orgID, map[string]string{"contact_id": cmd.String("contact-id")})
				if v := cmd.String("page"); v != "" {
					params["page"] = v
				}
				if v := cmd.String("per-page"); v != "" {
					params["per_page"] = v
				}
				raw, err := req(c, orgID, "GET", "/contacts/"+cmd.String("contact-id")+"/contact_persons", &zohttp.RequestOpts{Params: params})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "get",
			Usage:     "Get contact person",
			ArgsUsage: "<contact-id> <contact-person-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "GET", "/contacts/"+cmd.Args().Get(0)+"/contact_persons/"+cmd.Args().Get(1), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "create",
			Usage:     "Create contact person",
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
			Usage:     "Update contact person",
			ArgsUsage: "<contact-id> <contact-person-id>",
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
				raw, err := req(c, orgID, "PUT", "/contacts/"+cmd.Args().Get(0)+"/contact_persons/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete contact person",
			ArgsUsage: "<contact-id> <contact-person-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "DELETE", "/contacts/"+cmd.Args().Get(0)+"/contact_persons/"+cmd.Args().Get(1), nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
		{
			Name:      "mark-primary",
			Usage:     "Mark as primary contact person",
			ArgsUsage: "<contact-id> <contact-person-id>",
			Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().Get(0)+"/contact_persons/"+cmd.Args().Get(1)+"/primary", nil)
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	_ = path
	return &cli.Command{
		Name:     "contact-persons",
		Usage:    "Contact person operations",
		Commands: base,
	}
}
