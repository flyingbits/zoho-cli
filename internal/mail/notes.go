package mail

import (
	"context"
	"encoding/json"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func notesCmd() *cli.Command {
	return &cli.Command{
		Name:  "notes",
		Usage: "Notes API",
		Commands: []*cli.Command{
			{
				Name:      "create-group",
				Usage:     "Create a note (group)",
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
					return request(c, "POST", "/notes/groups/"+groupId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-personal",
				Usage: "Create a note (personal)",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/notes/me", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "create-book-group",
				Usage:     "Create a book (group)",
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
					return request(c, "POST", "/notes/groups/"+groupId+"/books", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-book-personal",
				Usage: "Create a book (personal)",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/notes/me/books", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "add-attachment-group",
				Usage:     "Add an attachment to a note (group)",
				ArgsUsage: "<groupId> <noteId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "file", Required: true, Usage: "File path"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return err
					}
					return request(c, "POST", "/notes/groups/"+groupId+"/"+noteId+"/attachments", &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{"attachment": {Filename: cmd.String("file"), Data: data}},
					})
				},
			},
			{
				Name:      "add-attachment-personal",
				Usage:     "Add an attachment to a note (personal)",
				ArgsUsage: "<noteId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "file", Required: true, Usage: "File path"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					data, err := os.ReadFile(cmd.String("file"))
					if err != nil {
						return err
					}
					return request(c, "POST", "/notes/me/"+noteId+"/attachments", &zohttp.RequestOpts{
						Files: map[string]zohttp.FileUpload{"attachment": {Filename: cmd.String("file"), Data: data}},
					})
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
					return request(c, "GET", "/notes/groups", nil)
				},
			},
			{
				Name:      "list-group",
				Usage:     "Get all notes (group)",
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
					return request(c, "GET", "/notes/groups/"+groupId, nil)
				},
			},
			{
				Name:  "list-personal",
				Usage: "Get all notes (personal)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/notes/me", nil)
				},
			},
			{
				Name:      "books-group",
				Usage:     "Get all books (group)",
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
					return request(c, "GET", "/notes/groups/"+groupId+"/books", nil)
				},
			},
			{
				Name:  "books-personal",
				Usage: "Get all books (personal)",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/notes/me/books", nil)
				},
			},
			{
				Name:  "favorites",
				Usage: "Get all favourite notes",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/notes/favorites", nil)
				},
			},
			{
				Name:  "shared",
				Usage: "Get all shared notes",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					return request(c, "GET", "/notes/sharedtome", nil)
				},
			},
			{
				Name:      "notes-in-book-group",
				Usage:     "Get all notes in a book (group)",
				ArgsUsage: "<groupId> <bookId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookId == "" {
						return internal.NewValidationError("groupId and bookId required")
					}
					return request(c, "GET", "/notes/groups/"+groupId+"/books/"+bookId, nil)
				},
			},
			{
				Name:      "notes-in-book-personal",
				Usage:     "Get all notes in a book (personal)",
				ArgsUsage: "<bookId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookId := cmd.Args().First()
					if bookId == "" {
						return internal.NewValidationError("bookId required")
					}
					return request(c, "GET", "/notes/me/books/"+bookId, nil)
				},
			},
			{
				Name:      "attachments-group",
				Usage:     "Get all attachments in a note (group)",
				ArgsUsage: "<groupId> <noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					return request(c, "GET", "/notes/groups/"+groupId+"/"+noteId+"/attachments", nil)
				},
			},
			{
				Name:      "attachments-personal",
				Usage:     "Get all attachments in a note (personal)",
				ArgsUsage: "<noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					return request(c, "GET", "/notes/me/"+noteId+"/attachments", nil)
				},
			},
			{
				Name:      "get-group",
				Usage:     "Get a note (group)",
				ArgsUsage: "<groupId> <noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					return request(c, "GET", "/notes/groups/"+groupId+"/"+noteId, nil)
				},
			},
			{
				Name:      "get-personal",
				Usage:     "Get a note (personal)",
				ArgsUsage: "<noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					return request(c, "GET", "/notes/me/"+noteId, nil)
				},
			},
			{
				Name:      "attachment-group",
				Usage:     "Get an attachment in a note (group)",
				ArgsUsage: "<groupId> <noteId> <attachmentId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId, attachmentId := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					if groupId == "" || noteId == "" || attachmentId == "" {
						return internal.NewValidationError("groupId, noteId, attachmentId required")
					}
					return request(c, "GET", "/notes/groups/"+groupId+"/"+noteId+"/attachments/"+attachmentId, nil)
				},
			},
			{
				Name:      "attachment-personal",
				Usage:     "Get an attachment in a note (personal)",
				ArgsUsage: "<noteId> <attachmentId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId, attachmentId := cmd.Args().Get(0), cmd.Args().Get(1)
					if noteId == "" || attachmentId == "" {
						return internal.NewValidationError("noteId and attachmentId required")
					}
					return request(c, "GET", "/notes/me/"+noteId+"/attachments/"+attachmentId, nil)
				},
			},
			{
				Name:      "update-group",
				Usage:     "Edit a note (group)",
				ArgsUsage: "<groupId> <noteId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/notes/groups/"+groupId+"/"+noteId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update-personal",
				Usage:     "Edit a note (personal)",
				ArgsUsage: "<noteId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/notes/me/"+noteId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update-book-group",
				Usage:     "Edit a book (group)",
				ArgsUsage: "<groupId> <bookId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookId == "" {
						return internal.NewValidationError("groupId and bookId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/notes/groups/"+groupId+"/books/"+bookId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "update-book-personal",
				Usage:     "Edit a book (personal)",
				ArgsUsage: "<bookId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookId := cmd.Args().First()
					if bookId == "" {
						return internal.NewValidationError("bookId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/notes/me/books/"+bookId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "favorite-group",
				Usage:     "Mark a note as favorite (group)",
				ArgsUsage: "<groupId> <noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					return request(c, "PUT", "/notes/groups/"+groupId+"/"+noteId+"/favorite", nil)
				},
			},
			{
				Name:      "favorite-personal",
				Usage:     "Mark a note as favorite (personal)",
				ArgsUsage: "<noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					return request(c, "PUT", "/notes/me/"+noteId+"/favorite", nil)
				},
			},
			{
				Name:      "unfavorite-group",
				Usage:     "Unmark a note as favorite (group)",
				ArgsUsage: "<groupId> <noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					return request(c, "DELETE", "/notes/groups/"+groupId+"/"+noteId+"/favorite", nil)
				},
			},
			{
				Name:      "unfavorite-personal",
				Usage:     "Unmark a note as favorite (personal)",
				ArgsUsage: "<noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					return request(c, "DELETE", "/notes/me/"+noteId+"/favorite", nil)
				},
			},
			{
				Name:      "delete-attachment-group",
				Usage:     "Delete an attachment in a note (group)",
				ArgsUsage: "<groupId> <noteId> <attachmentId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId, attachmentId := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					if groupId == "" || noteId == "" || attachmentId == "" {
						return internal.NewValidationError("groupId, noteId, attachmentId required")
					}
					return request(c, "DELETE", "/notes/groups/"+groupId+"/"+noteId+"/attachments/"+attachmentId, nil)
				},
			},
			{
				Name:      "delete-attachment-personal",
				Usage:     "Delete an attachment in a note (personal)",
				ArgsUsage: "<noteId> <attachmentId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId, attachmentId := cmd.Args().Get(0), cmd.Args().Get(1)
					if noteId == "" || attachmentId == "" {
						return internal.NewValidationError("noteId and attachmentId required")
					}
					return request(c, "DELETE", "/notes/me/"+noteId+"/attachments/"+attachmentId, nil)
				},
			},
			{
				Name:      "delete-book-group",
				Usage:     "Delete a book (group)",
				ArgsUsage: "<groupId> <bookId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, bookId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || bookId == "" {
						return internal.NewValidationError("groupId and bookId required")
					}
					return request(c, "DELETE", "/notes/groups/"+groupId+"/books/"+bookId, nil)
				},
			},
			{
				Name:      "delete-book-personal",
				Usage:     "Delete a book (personal)",
				ArgsUsage: "<bookId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					bookId := cmd.Args().First()
					if bookId == "" {
						return internal.NewValidationError("bookId required")
					}
					return request(c, "DELETE", "/notes/me/books/"+bookId, nil)
				},
			},
			{
				Name:      "delete-group",
				Usage:     "Delete a note (group)",
				ArgsUsage: "<groupId> <noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					groupId, noteId := cmd.Args().Get(0), cmd.Args().Get(1)
					if groupId == "" || noteId == "" {
						return internal.NewValidationError("groupId and noteId required")
					}
					return request(c, "DELETE", "/notes/groups/"+groupId+"/"+noteId, nil)
				},
			},
			{
				Name:      "delete-personal",
				Usage:     "Delete a note (personal)",
				ArgsUsage: "<noteId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					noteId := cmd.Args().First()
					if noteId == "" {
						return internal.NewValidationError("noteId required")
					}
					return request(c, "DELETE", "/notes/me/"+noteId, nil)
				},
			},
		},
	}
}
