package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func tasksCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasks",
		Usage: "Tasks API (group and personal)",
		Commands: []*cli.Command{
			{
				Name:      "add-group",
				Usage:     "Add a new group task",
				ArgsUsage: "<zgid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid := cmd.Args().First()
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/tasks/groups/"+zgid, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "add-personal",
				Usage: "Add a new personal task",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/tasks/me", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "add-project",
				Usage:     "Add a new project (group)",
				ArgsUsage: "<zgid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid := cmd.Args().First()
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/tasks/groups/"+zgid+"/projects", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list-group",
				Usage:     "Get all tasks in a group",
				ArgsUsage: "<zgid>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "status", Usage: "Filter by status"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid := cmd.Args().First()
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					params := map[string]string{}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					return request(c, "GET", "/tasks/groups/"+zgid, &zohttp.RequestOpts{Params: params})
				},
			},
			{
				Name:  "list-personal",
				Usage: "Get all personal tasks",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/tasks/me", nil)
				},
			},
			{
				Name:  "list-assigned",
				Usage: "Get all tasks assigned to you",
				Flags: []cli.Flag{&cli.StringFlag{Name: "filter", Usage: "Filter (e.g. assignee)"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{}
					if v := cmd.String("filter"); v != "" {
						params["filter"] = v
					}
					return request(c, "GET", "/tasks", &zohttp.RequestOpts{Params: params})
				},
			},
			{
				Name:  "list-created",
				Usage: "Get all tasks created by you",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/tasks", &zohttp.RequestOpts{Params: map[string]string{"filter": "createdByMe"}})
				},
			},
			{
				Name:      "get-group",
				Usage:     "Get a specific group task",
				ArgsUsage: "<zgid> <taskId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, taskId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || taskId == "" {
						return internal.NewValidationError("zgid and taskId required")
					}
					return request(c, "GET", "/tasks/groups/"+zgid+"/"+taskId, nil)
				},
			},
			{
				Name:      "get-personal",
				Usage:     "Get a specific personal task",
				ArgsUsage: "<taskId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					taskId := cmd.Args().First()
					if taskId == "" {
						return internal.NewValidationError("taskId required")
					}
					return request(c, "GET", "/tasks/me/"+taskId, nil)
				},
			},
			{
				Name:      "subtasks-group",
				Usage:     "Get all subtasks under a group task",
				ArgsUsage: "<zgid> <taskId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, taskId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || taskId == "" {
						return internal.NewValidationError("zgid and taskId required")
					}
					return request(c, "GET", "/tasks/groups/"+zgid+"/"+taskId+"/subtasks", nil)
				},
			},
			{
				Name:      "subtasks-personal",
				Usage:     "Get all subtasks under a personal task",
				ArgsUsage: "<taskId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					taskId := cmd.Args().First()
					if taskId == "" {
						return internal.NewValidationError("taskId required")
					}
					return request(c, "GET", "/tasks/me/"+taskId+"/subtasks", nil)
				},
			},
			{
				Name:      "projects-list",
				Usage:     "Get all projects in a group",
				ArgsUsage: "<zgid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid := cmd.Args().First()
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					return request(c, "GET", "/tasks/groups/"+zgid+"/projects", nil)
				},
			},
			{
				Name:      "project-tasks",
				Usage:     "Get all tasks in a project",
				ArgsUsage: "<zgid> <projectId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "status", Usage: "Filter by status"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, projectId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || projectId == "" {
						return internal.NewValidationError("zgid and projectId required")
					}
					params := map[string]string{}
					if v := cmd.String("status"); v != "" {
						params["status"] = v
					}
					return request(c, "GET", "/tasks/groups/"+zgid+"/projects/"+projectId, &zohttp.RequestOpts{Params: params})
				},
			},
			{
				Name:  "groups-list",
				Usage: "Get all groups",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/tasks/groups", nil)
				},
			},
			{
				Name:      "members",
				Usage:     "Get member details in a group",
				ArgsUsage: "<zgid>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid := cmd.Args().First()
					if zgid == "" {
						return internal.NewValidationError("zgid required")
					}
					return request(c, "GET", "/tasks/groups/"+zgid+"/members", nil)
				},
			},
			{
				Name:      "update-group",
				Usage:     "Update group task (title, description, priority, status, etc.; JSON body)",
				ArgsUsage: "<zgid> <taskId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, taskId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || taskId == "" {
						return internal.NewValidationError("zgid and taskId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/tasks/groups/"+zgid+"/"+taskId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update-personal",
				Usage:     "Update personal task (JSON body)",
				ArgsUsage: "<taskId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					taskId := cmd.Args().First()
					if taskId == "" {
						return internal.NewValidationError("taskId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/tasks/me/"+taskId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "edit-project",
				Usage:     "Edit a project",
				ArgsUsage: "<zgid> <projectId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, projectId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || projectId == "" {
						return internal.NewValidationError("zgid and projectId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/tasks/groups/"+zgid+"/projects/"+projectId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "delete-project",
				Usage:     "Delete a project",
				ArgsUsage: "<zgid> <projectId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, projectId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || projectId == "" {
						return internal.NewValidationError("zgid and projectId required")
					}
					return request(c, "DELETE", "/tasks/groups/"+zgid+"/projects/"+projectId, nil)
				},
			},
			{
				Name:      "delete-group",
				Usage:     "Delete a group task",
				ArgsUsage: "<zgid> <taskId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zgid, taskId := cmd.Args().Get(0), cmd.Args().Get(1)
					if zgid == "" || taskId == "" {
						return internal.NewValidationError("zgid and taskId required")
					}
					return request(c, "DELETE", "/tasks/groups/"+zgid+"/"+taskId, nil)
				},
			},
			{
				Name:      "delete-personal",
				Usage:     "Delete a personal task",
				ArgsUsage: "<taskId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					taskId := cmd.Args().First()
					if taskId == "" {
						return internal.NewValidationError("taskId required")
					}
					return request(c, "DELETE", "/tasks/me/"+taskId, nil)
				},
			},
		},
	}
}
