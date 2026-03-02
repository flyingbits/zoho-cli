package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func bookmarksCmd() *cli.Command {
	return &cli.Command{
		Name:  "bookmarks",
		Usage: "Bookmarks API (links)",
		Commands: []*cli.Command{
			{
				Name:      "create-group",
				Usage:     "Create a group bookmark",
				ArgsUsage: "<groupId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId := cmd.Args().First()
					if groupId == "" {
						return internal.NewValidationError("groupId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/links/groups/"+groupId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-personal",
				Usage: "Create a personal bookmark",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/links/me", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "create-collection-group",
				Usage:     "Create a collection (group)",
				ArgsUsage: "<groupId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId := cmd.Args().First()
					if groupId == "" {
						return internal.NewValidationError("groupId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/links/groups/"+groupId+"/collections", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-collection-personal",
				Usage: "Create a collection (personal)",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/links/me/collections", &zohttp.RequestOpts{JSON: body})
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
					return request(c, "GET", "/links/groups", nil)
				},
			},
			{
				Name:      "list-group",
				Usage:     "Get all bookmarks (group)",
				ArgsUsage: "<groupId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId := cmd.Args().First()
					if groupId == "" {
						return internal.NewValidationError("groupId required")
					}
					return request(c, "GET", "/links/groups/"+groupId, nil)
				},
			},
			{
				Name:  "list-personal",
				Usage: "Get all bookmarks (personal)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/links/me", nil)
				},
			},
			{
				Name:  "favorites",
				Usage: "Get all favourite bookmarks",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/links/favorites", nil)
				},
			},
			{
				Name:  "shared",
				Usage: "Get all shared bookmarks",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/links", nil)
				},
			},
			{
				Name:      "trash-group",
				Usage:     "Get all bookmarks in trash (group)",
				ArgsUsage: "<groupId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId := cmd.Args().First()
					if groupId == "" {
						return internal.NewValidationError("groupId required")
					}
					return request(c, "GET", "/links/groups/"+groupId+"/trash", nil)
				},
			},
			{
				Name:  "trash-personal",
				Usage: "Get all bookmarks in trash (personal)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/links/me/trash", nil)
				},
			},
			{
				Name:      "collections-group",
				Usage:     "Get all collections (group)",
				ArgsUsage: "<groupId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId := cmd.Args().First()
					if groupId == "" {
						return internal.NewValidationError("groupId required")
					}
					return request(c, "GET", "/links/groups/"+groupId+"/collections", nil)
				},
			},
			{
				Name:  "collections-personal",
				Usage: "Get all collections (personal)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/links/me/collections", nil)
				},
			},
			{
				Name:  "collections-all-groups",
				Usage: "Get all collections in groups",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/links/groups/collections", nil)
				},
			},
			{
				Name:      "collection-bookmarks-group",
				Usage:     "Get all bookmarks in a collection (group)",
				ArgsUsage: "<groupId> <collectionId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, collectionId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || collectionId == "" {
						return internal.NewValidationError("groupId and collectionId required")
					}
					return request(c, "GET", "/links/groups/"+groupId+"/collections/"+collectionId, nil)
				},
			},
			{
				Name:      "collection-bookmarks-personal",
				Usage:     "Get all bookmarks in a collection (personal)",
				ArgsUsage: "<collectionId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					collectionId := cmd.Args().First()
					if collectionId == "" {
						return internal.NewValidationError("collectionId required")
					}
					return request(c, "GET", "/links/me/collections/"+collectionId, nil)
				},
			},
			{
				Name:      "get-group",
				Usage:     "Get a bookmark (group)",
				ArgsUsage: "<groupId> <bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookmarkId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookmarkId == "" {
						return internal.NewValidationError("groupId and bookmarkId required")
					}
					return request(c, "GET", "/links/groups/"+groupId+"/"+bookmarkId, nil)
				},
			},
			{
				Name:      "get-personal",
				Usage:     "Get a bookmark (personal)",
				ArgsUsage: "<bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookmarkId := cmd.Args().First()
					if bookmarkId == "" {
						return internal.NewValidationError("bookmarkId required")
					}
					return request(c, "GET", "/links/me/"+bookmarkId, nil)
				},
			},
			{
				Name:      "update-group",
				Usage:     "Edit a bookmark (group)",
				ArgsUsage: "<groupId> <bookmarkId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookmarkId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookmarkId == "" {
						return internal.NewValidationError("groupId and bookmarkId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/links/groups/"+groupId+"/"+bookmarkId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update-personal",
				Usage:     "Edit a bookmark (personal)",
				ArgsUsage: "<bookmarkId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookmarkId := cmd.Args().First()
					if bookmarkId == "" {
						return internal.NewValidationError("bookmarkId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/links/me/"+bookmarkId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "restore",
				Usage:     "Restore a bookmark from trash",
				ArgsUsage: "<groupId> <bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookmarkId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookmarkId == "" {
						return internal.NewValidationError("groupId and bookmarkId required")
					}
					return request(c, "PUT", "/links/groups/"+groupId+"/"+bookmarkId+"/restore", nil)
				},
			},
			{
				Name:      "edit-collection-group",
				Usage:     "Edit a collection (group)",
				ArgsUsage: "<groupId> <collectionId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, collectionId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || collectionId == "" {
						return internal.NewValidationError("groupId and collectionId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/links/groups/"+groupId+"/collections/"+collectionId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "edit-collection-personal",
				Usage:     "Edit a collection (personal)",
				ArgsUsage: "<collectionId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					collectionId := cmd.Args().First()
					if collectionId == "" {
						return internal.NewValidationError("collectionId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/links/me/collections/"+collectionId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "favorite-group",
				Usage:     "Mark a bookmark as favorite (group)",
				ArgsUsage: "<groupId> <bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookmarkId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookmarkId == "" {
						return internal.NewValidationError("groupId and bookmarkId required")
					}
					return request(c, "PUT", "/links/groups/"+groupId+"/"+bookmarkId+"/favorite", nil)
				},
			},
			{
				Name:      "favorite-personal",
				Usage:     "Mark a bookmark as favorite (personal)",
				ArgsUsage: "<bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookmarkId := cmd.Args().First()
					if bookmarkId == "" {
						return internal.NewValidationError("bookmarkId required")
					}
					return request(c, "PUT", "/links/me/"+bookmarkId+"/favorite", nil)
				},
			},
			{
				Name:      "unfavorite-group",
				Usage:     "Unmark a bookmark as favorite (group)",
				ArgsUsage: "<groupId> <bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookmarkId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookmarkId == "" {
						return internal.NewValidationError("groupId and bookmarkId required")
					}
					return request(c, "DELETE", "/links/groups/"+groupId+"/"+bookmarkId+"/favorite", nil)
				},
			},
			{
				Name:      "unfavorite-personal",
				Usage:     "Unmark a bookmark as favorite (personal)",
				ArgsUsage: "<bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookmarkId := cmd.Args().First()
					if bookmarkId == "" {
						return internal.NewValidationError("bookmarkId required")
					}
					return request(c, "DELETE", "/links/me/"+bookmarkId+"/favorite", nil)
				},
			},
			{
				Name:      "delete-group",
				Usage:     "Delete a bookmark (group)",
				ArgsUsage: "<groupId> <bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookmarkId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookmarkId == "" {
						return internal.NewValidationError("groupId and bookmarkId required")
					}
					return request(c, "DELETE", "/links/groups/"+groupId+"/"+bookmarkId, nil)
				},
			},
			{
				Name:      "delete-personal",
				Usage:     "Delete a bookmark (personal)",
				ArgsUsage: "<bookmarkId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookmarkId := cmd.Args().First()
					if bookmarkId == "" {
						return internal.NewValidationError("bookmarkId required")
					}
					return request(c, "DELETE", "/links/me/"+bookmarkId, nil)
				},
			},
			{
				Name:      "delete-collection-group",
				Usage:     "Delete a collection (group)",
				ArgsUsage: "<groupId> <collectionId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, collectionId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || collectionId == "" {
						return internal.NewValidationError("groupId and collectionId required")
					}
					return request(c, "DELETE", "/links/groups/"+groupId+"/collections/"+collectionId, nil)
				},
			},
			{
				Name:      "delete-collection-personal",
				Usage:     "Delete a collection (personal)",
				ArgsUsage: "<collectionId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					collectionId := cmd.Args().First()
					if collectionId == "" {
						return internal.NewValidationError("collectionId required")
					}
					return request(c, "DELETE", "/links/me/collections/"+collectionId, nil)
				},
			},
		},
	}
}
