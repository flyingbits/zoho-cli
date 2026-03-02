package invoice

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
			{Name: "list", Usage: "List contacts", Flags: listFlags(), Action: runList("contacts")},
			{Name: "get", Usage: "Get a contact", ArgsUsage: "<contact-id>", Action: runGet("contacts")},
			{Name: "create", Usage: "Create a contact", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("contacts")},
			{Name: "update", Usage: "Update a contact", ArgsUsage: "<contact-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("contacts")},
			{Name: "delete", Usage: "Delete a contact", ArgsUsage: "<contact-id>", Action: runDelete("contacts")},
			{Name: "mark-active", Usage: "Mark contact as active", ArgsUsage: "<contact-id>", Action: runPost("contacts", "active")},
			{Name: "mark-inactive", Usage: "Mark contact as inactive", ArgsUsage: "<contact-id>", Action: runPost("contacts", "inactive")},
			{Name: "enable-portal", Usage: "Enable portal access", ArgsUsage: "<contact-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runContactPostJSON("contacts", "portal/enable")},
			{Name: "view-client-reviews", Usage: "View all client reviews", ArgsUsage: "<contact-id>", Action: runGetSub("contacts", "reviews")},
			{Name: "get-client-review", Usage: "Details of a particular client review", ArgsUsage: "<contact-id> <review-id>", Action: contactGetReview},
			{Name: "reply-client-review", Usage: "Reply to a client review", ArgsUsage: "<contact-id> <review-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: contactReplyReview},
			{Name: "enable-payment-reminders", Usage: "Enable payment reminders", ArgsUsage: "<contact-id>", Action: runPost("contacts", "paymentreminder/enable")},
			{Name: "disable-payment-reminders", Usage: "Disable payment reminders", ArgsUsage: "<contact-id>", Action: runPost("contacts", "paymentreminder/disable")},
			{Name: "email-statement", Usage: "Email statement", ArgsUsage: "<contact-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runContactPostJSON("contacts", "statements/email")},
			{Name: "get-statement-mail-content", Usage: "Get statement mail content", ArgsUsage: "<contact-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: contactGetStatementMail},
			{Name: "email", Usage: "Email contact", ArgsUsage: "<contact-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: contactEmail},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<contact-id>", Action: runGetSub("contacts", "comments")},
			{Name: "add-address", Usage: "Add additional address", ArgsUsage: "<contact-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: contactAddAddress},
			{Name: "get-addresses", Usage: "Get contact addresses", ArgsUsage: "<contact-id>", Action: runGetSub("contacts", "address")},
			{Name: "edit-address", Usage: "Edit additional address", ArgsUsage: "<contact-id> <address-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: contactEditAddress},
			{Name: "delete-address", Usage: "Delete additional address", ArgsUsage: "<contact-id> <address-id>", Action: contactDeleteAddress},
			{Name: "list-refunds", Usage: "List refunds", ArgsUsage: "<contact-id>", Action: runGetSub("contacts", "refunds")},
		},
	}
}

func listFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "page", Usage: "Page"},
		&cli.StringFlag{Name: "per-page", Usage: "Per page"},
		&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
		&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
		&cli.StringFlag{Name: "status", Usage: "active or inactive"},
	}
}

func contactGetReview(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/contacts/"+cmd.Args().First()+"/reviews/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func contactReplyReview(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/contacts/"+cmd.Args().First()+"/reviews/"+cmd.Args().Get(1)+"/reply", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func contactGetStatementMail(_ context.Context, cmd *cli.Command) error {
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
}

func contactEmail(_ context.Context, cmd *cli.Command) error {
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
}

func contactAddAddress(_ context.Context, cmd *cli.Command) error {
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
}

func contactEditAddress(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/contacts/"+cmd.Args().First()+"/address/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func contactDeleteAddress(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/contacts/"+cmd.Args().First()+"/address/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func runContactPostJSON(path, action string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
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
		raw, err := req(c, orgID, "POST", "/"+path+"/"+cmd.Args().First()+"/"+action, opts)
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
