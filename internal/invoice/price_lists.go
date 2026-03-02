package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func priceListsCmd() *cli.Command {
	return &cli.Command{
		Name:  "price-lists",
		Usage: "Price list operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List price lists",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "page", Usage: "Page"},
					&cli.StringFlag{Name: "per-page", Usage: "Per page"},
				},
				Action: runList("price-lists"),
			},
			{
				Name:      "get",
				Usage:     "Retrieve a price list",
				ArgsUsage: "<price-list-id>",
				Action:    runGet("price-lists"),
			},
			{
				Name:  "create",
				Usage: "Create a price list",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: runCreate("price-lists"),
			},
			{
				Name:      "update",
				Usage:     "Update a price list",
				ArgsUsage: "<price-list-id>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action:    runUpdate("price-lists"),
			},
			{
				Name:      "delete",
				Usage:     "Delete a price list",
				ArgsUsage: "<price-list-id>",
				Action:    runDelete("price-lists"),
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
		if v := cmd.String("name"); v != "" {
			params["name"] = v
		}
		if v := cmd.String("description"); v != "" {
			params["description"] = v
		}
		if v := cmd.String("item-id"); v != "" {
			params["item_id"] = v
		}
		if v := cmd.String("sku"); v != "" {
			params["sku"] = v
		}
		if v := cmd.String("type"); v != "" {
			params["type"] = v
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
