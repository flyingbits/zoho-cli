package invoice

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func timeEntriesCmd() *cli.Command {
	return &cli.Command{
		Name:  "time-entries",
		Usage: "Time entry operations (require project-id)",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List time entries", ArgsUsage: "<project-id>", Flags: listFlags(), Action: timeEntriesList},
			{Name: "log", Usage: "Log time entries", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: timeEntriesLog},
			{Name: "get", Usage: "Get a time entry", ArgsUsage: "<project-id> <time-entry-id>", Action: timeEntriesGet},
			{Name: "update", Usage: "Update time entry", ArgsUsage: "<project-id> <time-entry-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: timeEntriesUpdate},
			{Name: "delete", Usage: "Delete time entry", ArgsUsage: "<project-id> <time-entry-id>", Action: timeEntriesDelete},
			{Name: "delete-bulk", Usage: "Delete time entries", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body with time_entry_ids"}}, Action: timeEntriesDeleteBulk},
			{Name: "start-timer", Usage: "Start timer", ArgsUsage: "<project-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: timeEntriesStartTimer},
			{Name: "stop-timer", Usage: "Stop timer", ArgsUsage: "<project-id>", Action: timeEntriesStopTimer},
			{Name: "get-timer", Usage: "Get timer", ArgsUsage: "<project-id>", Action: timeEntriesGetTimer},
		},
	}
}

func timeEntriesList(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "GET", "/projects/"+cmd.Args().First()+"/timeentries", &zohttp.RequestOpts{Params: params})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesLog(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/projects/"+cmd.Args().First()+"/timeentries", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesGet(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/projects/"+cmd.Args().First()+"/timeentries/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesUpdate(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "PUT", "/projects/"+cmd.Args().First()+"/timeentries/"+cmd.Args().Get(1), &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesDelete(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "DELETE", "/projects/"+cmd.Args().First()+"/timeentries/"+cmd.Args().Get(1), nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesDeleteBulk(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "DELETE", "/projects/"+cmd.Args().First()+"/timeentries", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesStartTimer(_ context.Context, cmd *cli.Command) error {
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
	raw, err := req(c, orgID, "POST", "/projects/"+cmd.Args().First()+"/timeentries/timer/start", &zohttp.RequestOpts{JSON: body})
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesStopTimer(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "POST", "/projects/"+cmd.Args().First()+"/timeentries/timer/stop", nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}

func timeEntriesGetTimer(_ context.Context, cmd *cli.Command) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	orgID, err := resolveOrgID(cmd)
	if err != nil {
		return err
	}
	raw, err := req(c, orgID, "GET", "/projects/"+cmd.Args().First()+"/timeentries/timer", nil)
	if err != nil {
		return err
	}
	return output.JSONRaw(raw)
}
