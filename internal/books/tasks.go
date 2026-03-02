package books

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func tasksCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasks",
		Usage: "Books task operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List tasks", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}, &cli.StringFlag{Name: "page", Usage: "Page"}, &cli.StringFlag{Name: "per-page", Usage: "Per page"}}, Action: runListWithProject("tasks")},
			{Name: "get", Usage: "Get a task", ArgsUsage: "<task-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: runGetWithProject("tasks")},
			{Name: "add", Usage: "Add a task", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}, &cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateWithProject("tasks")},
			{Name: "update", Usage: "Update a task", ArgsUsage: "<task-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}, &cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdateWithProject("tasks")},
			{Name: "delete", Usage: "Delete a task", ArgsUsage: "<task-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: runDeleteWithProject("tasks")},
		},
	}
}

func runListWithProject(path string) func(context.Context, *cli.Command) error {
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
		raw, err := req(c, orgID, "GET", "/projects/"+cmd.String("project-id")+"/"+path, &zohttp.RequestOpts{Params: params})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runGetWithProject(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		params := mergeParams(orgID, map[string]string{"project_id": cmd.String("project-id")})
		raw, err := req(c, orgID, "GET", "/"+path+"/"+cmd.Args().First(), &zohttp.RequestOpts{Params: params})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runCreateWithProject(path string) func(context.Context, *cli.Command) error {
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
		raw, err := req(c, orgID, "POST", "/projects/"+cmd.String("project-id")+"/"+path, &zohttp.RequestOpts{JSON: body})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runUpdateWithProject(path string) func(context.Context, *cli.Command) error {
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
		params := mergeParams(orgID, map[string]string{"project_id": cmd.String("project-id")})
		raw, err := req(c, orgID, "PUT", "/"+path+"/"+cmd.Args().First(), &zohttp.RequestOpts{JSON: body, Params: params})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func runDeleteWithProject(path string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		params := mergeParams(orgID, map[string]string{"project_id": cmd.String("project-id")})
		raw, err := req(c, orgID, "DELETE", "/"+path+"/"+cmd.Args().First(), &zohttp.RequestOpts{Params: params})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}
