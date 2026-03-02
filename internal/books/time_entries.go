package books

import (
	"context"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func timeEntriesCmd() *cli.Command {
	return &cli.Command{
		Name:  "timeentries",
		Usage: "Time entry operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List time entries", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}, &cli.StringFlag{Name: "page", Usage: "Page"}, &cli.StringFlag{Name: "per-page", Usage: "Per page"}}, Action: runListWithProject("timeentries")},
			{Name: "get", Usage: "Get a time entry", ArgsUsage: "<timeentry-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: runGetWithProject("timeentries")},
			{Name: "log", Usage: "Log time entry", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}, &cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateWithProject("timeentries")},
			{Name: "update", Usage: "Update a time entry", ArgsUsage: "<timeentry-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}, &cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdateWithProject("timeentries")},
			{Name: "delete", Usage: "Delete a time entry", ArgsUsage: "<timeentry-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: runDeleteWithProject("timeentries")},
			{Name: "start-timer", Usage: "Start timer", ArgsUsage: "<timeentry-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: runPostWithProject("timeentries", "timer/start")},
			{Name: "stop-timer", Usage: "Stop timer", ArgsUsage: "<timeentry-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: runPostWithProject("timeentries", "timer/stop")},
			{Name: "get-timer", Usage: "Get timer", Flags: []cli.Flag{&cli.StringFlag{Name: "project-id", Required: true, Usage: "Project ID"}}, Action: func(_ context.Context, cmd *cli.Command) error {
				c, err := getClient()
				if err != nil {
					return err
				}
				orgID, err := resolveOrgID(cmd)
				if err != nil {
					return err
				}
				params := mergeParams(orgID, map[string]string{"project_id": cmd.String("project-id")})
				raw, err := req(c, orgID, "GET", "/projects/"+cmd.String("project-id")+"/timeentries/timer", &zohttp.RequestOpts{Params: params})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			}},
		},
	}
}

func runPostWithProject(path, action string) func(context.Context, *cli.Command) error {
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
		raw, err := req(c, orgID, "POST", "/"+path+"/"+cmd.Args().First()+"/"+action, &zohttp.RequestOpts{Params: params})
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}
