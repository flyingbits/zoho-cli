package invoice

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
		Usage: "Task operations (require project-id)",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List tasks", ArgsUsage: "<project-id>", Flags: listFlags(), Action: tasksList},
			{Name: "get", Usage: "Get a task", ArgsUsage: "<project-id> <task-id>", Action: tasksGet},
			{Name: "add", Usage: "Add a task", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: tasksCreate},
			{Name: "update", Usage: "Update a task", ArgsUsage: "<project-id> <task-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: tasksUpdate},
			{Name: "delete", Usage: "Delete a task", ArgsUsage: "<project-id> <task-id>", Action: tasksDelete},
		},
	}
}

func tasksList(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "GET", "/projects/"+cmd.Args().First()+"/tasks", &zohttp.RequestOpts{Params: params})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func tasksGet(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/projects/"+cmd.Args().First()+"/tasks/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func tasksCreate(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/projects/"+cmd.Args().First()+"/tasks", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func tasksUpdate(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/projects/"+cmd.Args().First()+"/tasks/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func tasksDelete(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/projects/"+cmd.Args().First()+"/tasks/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
