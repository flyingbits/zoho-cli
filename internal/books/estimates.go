package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func estimatesCmd() *cli.Command {
	return &cli.Command{
		Name:  "estimates",
		Usage: "Estimate operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List estimates", Flags: listFlags(), Action: runList("estimates")},
			{Name: "get", Usage: "Get an estimate", ArgsUsage: "<estimate-id>", Action: runGet("estimates")},
			{Name: "create", Usage: "Create an estimate", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("estimates")},
			{Name: "update", Usage: "Update an estimate", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("estimates")},
			{Name: "delete", Usage: "Delete an estimate", ArgsUsage: "<estimate-id>", Action: runDelete("estimates")},
			{Name: "mark-sent", Usage: "Mark as sent", ArgsUsage: "<estimate-id>", Action: runPost("estimates", "status/sent")},
			{Name: "mark-accepted", Usage: "Mark as accepted", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("estimates", "status/accepted")},
			{Name: "mark-declined", Usage: "Mark as declined", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("estimates", "status/declined")},
			{Name: "submit", Usage: "Submit for approval", ArgsUsage: "<estimate-id>", Action: runPost("estimates", "submit")},
			{Name: "approve", Usage: "Approve estimate", ArgsUsage: "<estimate-id>", Action: runPost("estimates", "approve")},
			{Name: "email", Usage: "Email estimate", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Usage: "JSON body"}}, Action: runPostJSON("estimates", "email")},
			{Name: "get-email-content", Usage: "Get estimate email content", ArgsUsage: "<estimate-id>", Action: runGetSub("estimates", "email")},
			{Name: "list-comments", Usage: "List comments", ArgsUsage: "<estimate-id>", Action: runGetSub("estimates", "comments")},
			{Name: "add-comment", Usage: "Add comment", ArgsUsage: "<estimate-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("estimates", "comments")},
		},
	}
}

func listFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "page", Usage: "Page"},
		&cli.StringFlag{Name: "per-page", Usage: "Per page"},
		&cli.StringFlag{Name: "sort-column", Usage: "Sort column"},
		&cli.StringFlag{Name: "sort-order", Usage: "asc or desc"},
		&cli.StringFlag{Name: "status", Usage: "Filter by status"},
	}
}

func runCreateSub(path, sub string) func(context.Context, *cli.Command) error {
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
		raw, err := req(c, orgID, "POST", "/"+path+"/"+cmd.Args().First()+"/"+sub, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runUpdateSub(path, sub string) func(context.Context, *cli.Command) error {
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
		currencyID := cmd.Args().First()
		rateID := cmd.Args().Get(1)
		raw, err := req(c, orgID, "PUT", "/"+path+"/"+currencyID+"/"+sub+"/"+rateID, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runPostJSON(path, action string) func(context.Context, *cli.Command) error {
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
