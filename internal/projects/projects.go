package projects

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/omin8tor/zoho-cli/internal/pagination"
	"github.com/urfave/cli/v3"
)

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

func base(c *zohttp.Client, portal, project string) string {
	return c.ProjectsBase + "/portal/" + portal + "/projects/" + project
}

var portalFlag = &cli.StringFlag{Name: "portal", Required: true, Usage: "Portal ID"}
var projectFlag = &cli.StringFlag{Name: "project", Required: true, Usage: "Project ID"}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Zoho Projects operations",
		Commands: []*cli.Command{
			projectsCoreCmd(),
			tasksCmd(),
			issuesCmd(),
			issueCommentsCmd(),
			commentsCmd(),
			tasklistsCmd(),
			timelogsCmd(),
			usersCmd(),
			milestonesCmd(),
			dependenciesCmd(),
		},
	}
}

func projectsCoreCmd() *cli.Command {
	return &cli.Command{
		Name:  "core",
		Usage: "Project CRUD operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all projects",
				Flags: []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + cmd.String("portal") + "/projects"
					items, err := pagination.PaginateProjects(c, url, "", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single project",
				ArgsUsage: "<project-id>",
				Flags:     []cli.Flag{portalFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + cmd.String("portal") + "/projects/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "search",
				Usage: "Search projects",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + cmd.String("portal") + "/search"
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"search_term": cmd.String("query"),
							"module":      "all",
							"status":      "all",
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a project",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Project name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := c.ProjectsBase + "/portal/" + cmd.String("portal") + "/projects"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a project",
				ArgsUsage: "<project-id>",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := c.ProjectsBase + "/portal/" + cmd.String("portal") + "/projects/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tasksCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasks",
		Usage: "Project task operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tasks in a project",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "status", Usage: "Filter: open, closed, in progress"},
					&cli.StringFlag{Name: "priority", Usage: "Filter: none, low, medium, high"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks"
					params := map[string]string{}
					if s := cmd.String("status"); s != "" {
						params["status"] = s
					}
					if p := cmd.String("priority"); p != "" {
						params["priority"] = p
					}
					items, err := pagination.PaginateProjects(c, url, "tasks", params, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:  "my",
				Usage: "List my tasks across all projects",
				Flags: []cli.Flag{
					portalFlag,
					&cli.StringFlag{Name: "status", Usage: "Filter: open, closed, in progress"},
					&cli.StringFlag{Name: "priority", Usage: "Filter: none, low, medium, high"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := c.ProjectsBase + "/portal/" + cmd.String("portal") + "/tasks"
					params := map[string]string{}
					if s := cmd.String("status"); s != "" {
						params["status"] = s
					}
					if p := cmd.String("priority"); p != "" {
						params["priority"] = p
					}
					items, err := pagination.PaginateProjects(c, url, "tasks", params, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single task",
				ArgsUsage: "<task-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a task",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Task name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a task",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a task",
				ArgsUsage: "<task-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "subtasks",
				Usage:     "List subtasks of a task",
				ArgsUsage: "<task-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/subtasks"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-subtask",
				Usage: "Create a subtask",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "parent", Required: true, Usage: "Parent task ID"},
					&cli.StringFlag{Name: "name", Required: true, Usage: "Subtask name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.String("parent") + "/subtasks"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issuesCmd() *cli.Command {
	return &cli.Command{
		Name:  "issues",
		Usage: "Project issue operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issues in a project",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues"
					items, err := pagination.PaginateProjects(c, url, "issues", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a single issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create an issue",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Issue title"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update an issue",
				ArgsUsage: "<issue-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete an issue",
				ArgsUsage: "<issue-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "defaults",
				Usage: "Get issue default fields (statuses, severities, etc.)",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues/defaultfields"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func issueCommentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "issue-comments",
		Usage: "Issue comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List issue comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add an issue comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "issue", Required: true, Usage: "Issue ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/issues/" + cmd.String("issue") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func commentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "comments",
		Usage: "Task comment operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List task comments",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments"
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a task comment",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a task comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
					&cli.StringFlag{Name: "comment", Required: true, Usage: "Updated comment text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: map[string]string{"comment": cmd.String("comment")}})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a task comment",
				ArgsUsage: "<comment-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "task", Required: true, Usage: "Task ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.String("task") + "/comments/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tasklistsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tasklists",
		Usage: "Project tasklist operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List tasklists",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasklists"
					items, err := pagination.PaginateProjects(c, url, "tasklists", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:  "create",
				Usage: "Create a tasklist",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Tasklist name"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{"name": cmd.String("name")}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasklists"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasklists/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a tasklist",
				ArgsUsage: "<tasklist-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasklists/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func timelogsCmd() *cli.Command {
	return &cli.Command{
		Name:  "timelogs",
		Usage: "Project timelog operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project timelogs",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "module", Value: "general", Usage: "task, issue, or general"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/timelogs"
					moduleJSON := fmt.Sprintf(`{"type":"%s"}`, cmd.String("module"))
					raw, err := c.Request("GET", url, &zohttp.RequestOpts{
						Params: map[string]string{
							"module":    moduleJSON,
							"view_type": "projectspan",
						},
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add a timelog",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "date", Required: true, Usage: "Date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "hours", Required: true, Usage: "Hours (e.g. 2, 1.5, 0:30)"},
					&cli.StringFlag{Name: "task", Usage: "Task ID"},
					&cli.StringFlag{Name: "bill-status", Value: "Billable", Usage: "Billable or Non Billable"},
					&cli.StringFlag{Name: "notes", Usage: "Notes for time entry"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{
						"date":        cmd.String("date"),
						"hours":       cmd.String("hours"),
						"bill_status": cmd.String("bill-status"),
						"log_name":    "Time log",
					}
					if n := cmd.String("notes"); n != "" {
						body["notes"] = n
						body["log_name"] = n
					}
					if t := cmd.String("task"); t != "" {
						body["module"] = map[string]string{"type": "task", "id": t}
					} else {
						body["module"] = map[string]string{"type": "general"}
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/log"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func usersCmd() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Project user operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List project users",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/users"
					items, err := pagination.PaginateProjects(c, url, "users", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
		},
	}
}

func milestonesCmd() *cli.Command {
	return &cli.Command{
		Name:  "milestones",
		Usage: "Project milestone operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List milestones",
				Flags: []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/milestones"
					items, err := pagination.PaginateProjects(c, url, "milestones", nil, 0)
					if err != nil {
						return err
					}
					return output.JSON(items)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a milestone",
				ArgsUsage: "<milestone-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request("GET", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create a milestone",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Milestone name"},
					&cli.StringFlag{Name: "start", Required: true, Usage: "Start date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "end", Required: true, Usage: "End date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "json", Usage: "Additional fields as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body := map[string]any{
						"name":       cmd.String("name"),
						"start_date": cmd.String("start"),
						"end_date":   cmd.String("end"),
					}
					if j := cmd.String("json"); j != "" {
						var extra map[string]any
						json.Unmarshal([]byte(j), &extra)
						for k, v := range extra {
							body[k] = v
						}
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/milestones"
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "update",
				Usage:     "Update a milestone",
				ArgsUsage: "<milestone-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "json", Required: true, Usage: "Fields to update as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var parsed map[string]any
					json.Unmarshal([]byte(cmd.String("json")), &parsed)
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request("PATCH", url, &zohttp.RequestOpts{JSON: parsed})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a milestone",
				ArgsUsage: "<milestone-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/milestones/" + cmd.Args().First()
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func dependenciesCmd() *cli.Command {
	return &cli.Command{
		Name:  "dependencies",
		Usage: "Task dependency operations",
		Commands: []*cli.Command{
			{
				Name:      "add",
				Usage:     "Add a task dependency",
				ArgsUsage: "<task-id>",
				Flags: []cli.Flag{
					portalFlag, projectFlag,
					&cli.StringFlag{Name: "depends-on", Required: true, Usage: "Dependency task ID"},
					&cli.StringFlag{Name: "type", Value: "FS", Usage: "FS, SS, FF, or SF"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.Args().First() + "/dependencies"
					body := map[string]any{
						"predecessor": map[string]string{
							"id":   cmd.String("depends-on"),
							"type": cmd.String("type"),
						},
					}
					raw, err := c.Request("POST", url, &zohttp.RequestOpts{JSON: body})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "remove",
				Usage:     "Remove a task dependency",
				ArgsUsage: "<task-id> <dependency-id>",
				Flags:     []cli.Flag{portalFlag, projectFlag},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					url := base(c, cmd.String("portal"), cmd.String("project")) + "/tasks/" + cmd.Args().Get(0) + "/dependencies/" + cmd.Args().Get(1)
					raw, err := c.Request("DELETE", url, nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
