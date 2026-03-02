package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func messagesCmd() *cli.Command {
	return &cli.Command{
		Name:  "messages",
		Usage: "Email Messages API",
		Commands: []*cli.Command{
			{
				Name:  "send",
				Usage: "Send an email (or with attachment via --json)",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: runAccountJSON("POST", "/accounts/%s/messages"),
			},
			{
				Name:      "upload-attachments",
				Usage:     "Upload attachments",
				ArgsUsage: "<accountId>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "JSON body"},
					&cli.StringSliceFlag{Name: "file", Usage: "File path(s) to upload"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, err := mustAccount(cmd)
					if err != nil {
						return err
					}
					files := cmd.StringSlice("file")
					if len(files) == 0 {
						return internal.NewValidationError("at least one --file required")
					}
					fileMap := make(map[string]zohttp.FileUpload)
					for i, p := range files {
						data, err := os.ReadFile(p)
						if err != nil {
							return err
						}
						fileMap[fmt.Sprintf("attachment%d", i)] = zohttp.FileUpload{Filename: p, Data: data}
					}
					opts := &zohttp.RequestOpts{Files: fileMap}
					if j := cmd.String("json"); j != "" {
						var body any
						json.Unmarshal([]byte(j), &body)
						opts.JSON = body
					}
					return request(c, "POST", "/accounts/"+accountId+"/messages/attachments", opts)
				},
			},
			{
				Name:  "save-draft",
				Usage: "Save draft or template",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: runAccountJSON("POST", "/accounts/%s/messages"),
			},
			{
				Name:      "reply",
				Usage:     "Reply to an email",
				ArgsUsage: "<accountId> <messageId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, messageId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || messageId == "" {
						return internal.NewValidationError("accountId and messageId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/accounts/"+accountId+"/messages/"+messageId, &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list",
				Usage:     "Get list of emails in a folder",
				ArgsUsage: "<accountId>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "folderId", Required: true, Usage: "Folder ID"},
					&cli.StringFlag{Name: "start", Usage: "Start index"},
					&cli.StringFlag{Name: "limit", Usage: "Limit"},
				},
				Action: runMessagesListOrSearch("/accounts/%s/messages/view", "folderId"),
			},
			{
				Name:      "search",
				Usage:     "Get list of search results",
				ArgsUsage: "<accountId>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "query", Required: true, Usage: "Search query"},
					&cli.StringFlag{Name: "start", Usage: "Start index"},
					&cli.StringFlag{Name: "limit", Usage: "Limit"},
				},
				Action: runMessagesListOrSearch("/accounts/%s/messages/search", "query"),
			},
			{
				Name:      "header",
				Usage:     "Get email headers",
				ArgsUsage: "<accountId> <folderId> <messageId>",
				Action:    runMsgPath3("/accounts/%s/folders/%s/messages/%s/header"),
			},
			{
				Name:      "content",
				Usage:     "Get email content",
				ArgsUsage: "<accountId> <folderId> <messageId>",
				Action:    runMsgPath3("/accounts/%s/folders/%s/messages/%s/content"),
			},
			{
				Name:      "original",
				Usage:     "Get original message",
				ArgsUsage: "<accountId> <messageId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, messageId := cmd.Args().Get(0), cmd.Args().Get(1)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || messageId == "" {
						return internal.NewValidationError("accountId and messageId required")
					}
					return request(c, "GET", "/accounts/"+accountId+"/messages/"+messageId+"/originalmessage", nil)
				},
			},
			{
				Name:      "details",
				Usage:     "Get metadata of an email",
				ArgsUsage: "<accountId> <folderId> <messageId>",
				Action:    runMsgPath3("/accounts/%s/folders/%s/messages/%s/details"),
			},
			{
				Name:      "attachment-info",
				Usage:     "Get attachment info",
				ArgsUsage: "<accountId> <folderId> <messageId>",
				Action:    runMsgPath3("/accounts/%s/folders/%s/messages/%s/attachmentinfo"),
			},
			{
				Name:      "attachment",
				Usage:     "Get attachment content",
				ArgsUsage: "<accountId> <folderId> <messageId> <attachmentId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					a1, a2, a3, a4 := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2), cmd.Args().Get(3)
					if a1 == "" {
						a1, _ = resolveAccount(cmd)
					}
					if a1 == "" || a2 == "" || a3 == "" || a4 == "" {
						return internal.NewValidationError("accountId, folderId, messageId, attachmentId required")
					}
					return request(c, "GET", "/accounts/"+a1+"/folders/"+a2+"/messages/"+a3+"/attachments/"+a4, nil)
				},
			},
			{
				Name:      "inline",
				Usage:     "Download inline image",
				ArgsUsage: "<accountId> <folderId> <messageId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "contentId", Usage: "Content-ID of inline image"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, folderId, messageId := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || folderId == "" || messageId == "" {
						return internal.NewValidationError("accountId, folderId, messageId required")
					}
					params := map[string]string{}
					if v := cmd.String("contentId"); v != "" {
						params["contentId"] = v
					}
					return request(c, "GET", "/accounts/"+accountId+"/folders/"+folderId+"/messages/"+messageId+"/inline", &zohttp.RequestOpts{Params: params})
				},
			},
			{
				Name:      "update",
				Usage:     "Mark read/unread, move, flag, labels, archive, spam (use JSON body)",
				ArgsUsage: "<accountId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action:    runAccountJSON("PUT", "/accounts/%s/updatemessage"),
			},
			{
				Name:      "delete",
				Usage:     "Delete email",
				ArgsUsage: "<accountId> <folderId> <messageId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					accountId, folderId, messageId := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
					if accountId == "" {
						accountId, _ = resolveAccount(cmd)
					}
					if accountId == "" || folderId == "" || messageId == "" {
						return internal.NewValidationError("accountId, folderId, messageId required")
					}
					return request(c, "DELETE", "/accounts/"+accountId+"/folders/"+folderId+"/messages/"+messageId, nil)
				},
			},
		},
	}
}

func runAccountJSON(method, pathFmt string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		accountId, err := mustAccount(cmd)
		if err != nil {
			return err
		}
		var body any
		json.Unmarshal([]byte(cmd.String("json")), &body)
		return request(c, method, fmt.Sprintf(pathFmt, accountId), &zohttp.RequestOpts{JSON: body})
	}
}

func runMessagesListOrSearch(pathFmt, keyParam string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		accountId := cmd.Args().First()
		if accountId == "" {
			accountId, _ = resolveAccount(cmd)
		}
		if accountId == "" {
			return internal.NewValidationError(accountRequiredMsg)
		}
		params := map[string]string{}
		if keyParam == "folderId" {
			params["folderId"] = cmd.String("folderId")
		} else {
			params["query"] = cmd.String("query")
		}
		if v := cmd.String("start"); v != "" {
			params["start"] = v
		}
		if v := cmd.String("limit"); v != "" {
			params["limit"] = v
		}
		return request(c, "GET", fmt.Sprintf(pathFmt, accountId), &zohttp.RequestOpts{Params: params})
	}
}

func runMsgPath3(pathFmt string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		c, err := getClient()
		if err != nil {
			return err
		}
		a1, a2, a3 := cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2)
		if a1 == "" {
			a1, _ = resolveAccount(cmd)
		}
		if a1 == "" || a2 == "" || a3 == "" {
			return internal.NewValidationError("accountId, folderId, messageId required")
		}
		return request(c, "GET", fmt.Sprintf(pathFmt, a1, a2, a3), nil)
	}
}
