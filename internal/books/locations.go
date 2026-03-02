package books

import (
	"github.com/urfave/cli/v3"
)

func locationsCmd() *cli.Command {
	return &cli.Command{
		Name:  "locations",
		Usage: "Location operations",
		Commands: []*cli.Command{
			{Name: "list", Usage: "List all locations", Flags: listFlags(), Action: runList("settings/locations")},
			{Name: "get", Usage: "Get a location", ArgsUsage: "<location-id>", Action: runGet("settings/locations")},
			{Name: "create", Usage: "Create a location", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("settings/locations")},
			{Name: "update", Usage: "Update a location", ArgsUsage: "<location-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdate("settings/locations")},
			{Name: "delete", Usage: "Delete a location", ArgsUsage: "<location-id>", Action: runDelete("settings/locations")},
			{Name: "mark-active", Usage: "Mark location as active", ArgsUsage: "<location-id>", Action: runPost("settings/locations", "active")},
			{Name: "mark-inactive", Usage: "Mark location as inactive", ArgsUsage: "<location-id>", Action: runPost("settings/locations", "inactive")},
			{Name: "mark-primary", Usage: "Mark as primary", ArgsUsage: "<location-id>", Action: runPost("settings/locations", "primary")},
		},
	}
}
