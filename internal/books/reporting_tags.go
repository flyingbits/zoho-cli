package books

import (
	"context"

	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func runPostOption(status string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		orgID, err := resolveOrgID(cmd)
		if err != nil {
			return err
		}
		tagID := cmd.Args().First()
		optionID := cmd.Args().Get(1)
		raw, err := req(c, orgID, "POST", "/settings/reportingtags/"+tagID+"/option/"+optionID+"/"+status, nil)
		if err != nil {
			return err
		}
		return output.JSONRaw(raw)
	}
}

func reportingTagsCmd() *cli.Command {
	return &cli.Command{
		Name:  "reportingtags",
		Usage: "Reporting tag operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List all reporting tags", Flags: listFlags(), Action: runList("settings/reportingtags")},
			{Name: "get", Usage: "Get reporting tag options detail", ArgsUsage: "<reportingtag-id>", Action: runGet("settings/reportingtags")},
			{Name: "create", Usage: "Create reporting tag", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("settings/reportingtags")},
			{Name: "update", Usage: "Update reporting tag", ArgsUsage: "<reportingtag-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("settings/reportingtags")},
			{Name: "delete", Usage: "Delete a reporting tag", ArgsUsage: "<reportingtag-id>", Action: runDelete("settings/reportingtags")},
			{Name: "mark-option-default", Usage: "Mark an option as default", ArgsUsage: "<reportingtag-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("settings/reportingtags", "default")},
			{Name: "update-options", Usage: "Update reporting tag options", ArgsUsage: "<reportingtag-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("settings/reportingtags", "options")},
			{Name: "update-visibility", Usage: "Update visibility conditions", ArgsUsage: "<reportingtag-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreateSub("settings/reportingtags", "visibility")},
			{Name: "mark-active", Usage: "Mark reporting tag as active", ArgsUsage: "<reportingtag-id>", Action: runPost("settings/reportingtags", "active")},
			{Name: "mark-inactive", Usage: "Mark reporting tag as inactive", ArgsUsage: "<reportingtag-id>", Action: runPost("settings/reportingtags", "inactive")},
			{Name: "mark-option-active", Usage: "Mark an option as active", ArgsUsage: "<reportingtag-id> <option-id>", Action: runPostOption("active")},
			{Name: "mark-option-inactive", Usage: "Mark an option as inactive", ArgsUsage: "<reportingtag-id> <option-id>", Action: runPostOption("inactive")},
			{Name: "get-options-detail", Usage: "Get reporting tag options detail page", ArgsUsage: "<reportingtag-id>", Action: runGetSub("settings/reportingtags", "options")},
			{Name: "get-all-options", Usage: "Get all options", Action: runGetNoID("settings/reportingtags/options")},
			{Name: "reorder", Usage: "Reorder reporting tags", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("settings/reportingtags/reorder")},
		},
	}
}
