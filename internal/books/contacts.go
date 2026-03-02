package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func contactsCmd() *cli.Command {
	return &cli.Command{
		Name:  "contacts",
		Usage: "Contact operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List contacts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page"},
					&cli.StringFlag{Name: "per-page", Usage: "Per page"},
					&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
					&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
				},
				Action: runList("contacts"),
			},
			{
				Name:      "get",
				Usage:     "Get a contact",
				ArgsUsage: "<contact-id>",
				Action:    runGet("contacts"),
			},
			{
				Name:  "create",
				Usage: "Create a contact",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: runCreate("contacts"),
			},
			{
				Name:      "update",
				Usage:     "Update a contact",
				ArgsUsage: "<contact-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action:    runUpdate("contacts"),
			},
			{
				Name:      "delete",
				Usage:     "Delete a contact",
				ArgsUsage: "<contact-id>",
				Action:    runDelete("contacts"),
			},
			{
				Name:      "mark-active",
				Usage:     "Mark contact as active",
				ArgsUsage: "<contact-id>",
				Action:    runPost("contacts", "active"),
			},
			{
				Name:      "mark-inactive",
				Usage:     "Mark contact as inactive",
				ArgsUsage: "<contact-id>",
				Action:    runPost("contacts", "inactive"),
			},
			{
				Name:      "enable-portal",
				Usage:     "Enable portal access",
				ArgsUsage: "<contact-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().First()+"/portal/enable", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "enable-payment-reminders",
				Usage:     "Enable payment reminders",
				ArgsUsage: "<contact-id>",
				Action:    runPost("contacts", "paymentreminder/enable"),
			},
			{
				Name:      "disable-payment-reminders",
				Usage:     "Disable payment reminders",
				ArgsUsage: "<contact-id>",
				Action:    runPost("contacts", "paymentreminder/disable"),
			},
			{
				Name:      "email-statement",
				Usage:     "Email statement",
				ArgsUsage: "<contact-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					opts := &zohttp.RequestOpts{}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().First()+"/statements/email", opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-statement-mail-content",
				Usage:     "Get statement mail content",
				ArgsUsage: "<contact-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					orgID, err := resolveOrgID(cmd)
					if err != nil {
						return err
					}
					params := orgParam(orgID)
					if j := cmd.String("json"); j != "" {
						var m map[string]string
						json.Unmarshal([]byte(j), &m)
						for k, v := range m {
							params[k] = v
						}
					}
					raw, err := req(c, orgID, "GET", "/contacts/"+cmd.Args().First()+"/statements/email", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "email",
				Usage:     "Email contact",
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
					raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().First()+"/email", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "list-comments",
				Usage:     "List comments",
				ArgsUsage: "<contact-id>",
				Action:    runGetSub("contacts", "comments"),
			},
			{
				Name:      "add-address",
				Usage:     "Add additional address",
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
					raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().First()+"/address", &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "get-addresses",
				Usage:     "Get contact addresses",
				ArgsUsage: "<contact-id>",
				Action:    runGetSub("contacts", "address"),
			},
			{
				Name:      "list-refunds",
				Usage:     "List refunds",
				ArgsUsage: "<contact-id>",
				Action:    runGetSub("contacts", "refunds"),
			},
			{
				Name:      "track-1099",
				Usage:     "Track 1099",
				ArgsUsage: "<contact-id>",
				Action:    runPost("contacts", "tracks1099"),
			},
			{
				Name:      "untrack-1099",
				Usage:     "Untrack 1099",
				ArgsUsage: "<contact-id>",
				Action:    runPost("contacts", "untracks1099"),
			},
			{
				Name:      "get-unused-retainer-payments",
				Usage:     "Get unused retainer payments",
				ArgsUsage: "<contact-id>",
				Action:    runGetSub("contacts", "unusedretainerpayments"),
			},
		},
	}
}

func runList(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		params := mergeParams(orgID, nil)
		if v := cmd.String("page"); v != "" {
			params["page"] = v
		}
		if v := cmd.String("per-page"); v != "" {
			params["per_page"] = v
		}
		if v := cmd.String("sort-column"); v != "" {
			params["sort_column"] = v
		}
		if v := cmd.String("sort-order"); v != "" {
			params["sort_order"] = v
		}
		if v := cmd.String("status"); v != "" {
			params["status"] = v
		}
		raw, err := req(c, orgID, "GET", "/"+path, &zohttp.RequestOpts{Params: params})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runGet(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "GET", "/"+path+"/"+cmd.Args().First(), nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runGetSub(path, sub string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "GET", "/"+path+"/"+cmd.Args().First()+"/"+sub, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runCreate(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
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
		raw, err := req(c, orgID, "POST", "/"+path, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runUpdate(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
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
		raw, err := req(c, orgID, "PUT", "/"+path+"/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runDelete(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "DELETE", "/"+path+"/"+cmd.Args().First(), nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runPost(path, action string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "POST", "/"+path+"/"+cmd.Args().First()+"/"+action, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runGetNoID(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "GET", "/"+path, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runUpdateNoID(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
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
		raw, err := req(c, orgID, "PUT", "/"+path, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runDeleteNoID(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "DELETE", "/"+path, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runDeleteSub(path, sub string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		raw, err := req(c, orgID, "DELETE", "/"+path+"/"+cmd.Args().First()+"/"+sub, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}
