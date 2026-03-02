package inventory

import (
	"context"
	"encoding/json"

	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func reportingTagsCmd() *cli.Command {
	base := crudSubcommands("reportingtags")
	path := "/reportingtags/%s"
	extra := []*cli.Command{
		{Name: "mark-default", Usage: "Mark option as default", ArgsUsage: "<reporting-tag-id> <option-id>", Action: invReq2("POST", "/reportingtags/%s/options/%s/default")},
		{Name: "update-options", Usage: "Update reporting tag options", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, path+"/options")},
		{Name: "update-visibility", Usage: "Update visibility conditions", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}}, Action: invReq("PUT", true, path+"/visibility")},
		{Name: "mark-active", Usage: "Mark reporting tag as active", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/active")},
		{Name: "mark-inactive", Usage: "Mark reporting tag as inactive", ArgsUsage: "<id>", Action: invReq("POST", false, path+"/inactive")},
		{Name: "mark-option-active", Usage: "Mark option as active", ArgsUsage: "<reporting-tag-id> <option-id>", Action: invReq2("POST", "/reportingtags/%s/options/%s/active")},
		{Name: "mark-option-inactive", Usage: "Mark option as inactive", ArgsUsage: "<reporting-tag-id> <option-id>", Action: invReq2("POST", "/reportingtags/%s/options/%s/inactive")},
		{Name: "get-options-detail", Usage: "Get reporting tags options detail page", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/options")},
		{Name: "get-all-options", Usage: "Get all options", ArgsUsage: "<id>", Action: invReq("GET", false, path+"/options")},
		{
			Name:  "reorder",
			Usage: "Reorder reporting tags",
			Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true}},
			Action: func(_ context.Context, cmd *cli.Command) error {
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
				raw, err := req(c, orgID, "PUT", "/reportingtags/reorder", &zohttp.RequestOpts{JSON: body})
				if err != nil {
					return err
				}
				return output.JSONRaw(raw)
			},
		},
	}
	return &cli.Command{
		Name:     "reporting-tags",
		Usage:    "Reporting tag operations",
		Commands: append(base, extra...),
	}
}
