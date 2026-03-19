package writer

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "writer",
		Usage: "Zoho Writer operations",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List documents in your Writer account",
				ArgsUsage: "",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "offset", Value: 0, Usage: "Pagination offset"},
					&cli.IntFlag{Name: "limit", Value: 10, Usage: "Page size (number of documents)"},
					&cli.StringFlag{Name: "sortby", Value: "modified_time", Usage: "Sort by: created_time, modified_time, or name"},
					&cli.StringFlag{Name: "sort-order-by", Value: "descending", Usage: "Sort order: ascending or descending"},
					&cli.StringFlag{Name: "category", Value: "all", Usage: "all, shared_to_me, or owned_by_me"},
					&cli.StringFlag{Name: "type", Usage: "Document type filter: fillable, merge, or sign"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}

					params := map[string]string{
						"offset":          strconv.Itoa(cmd.Int("offset")),
						"limit":           strconv.Itoa(cmd.Int("limit")),
						"sortby":          cmd.String("sortby"),
						"sort_order_by":  cmd.String("sort-order-by"),
						"category":       cmd.String("category"),
					}
					if t := cmd.String("type"); t != "" {
						params["resource_type"] = t
					}

					raw, err := c.Request("GET", c.WriterBase+"/documents", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create/upload a new Writer document",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "filename", Usage: "Document name (optional). If omitted, Zoho creates an untitled blank document."},
					&cli.StringFlag{Name: "type", Usage: "Document type: fillable, merge, or sign"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}

					form := map[string]string{}
					if f := cmd.String("filename"); f != "" {
						form["filename"] = f
					}
					if t := cmd.String("type"); t != "" {
						form["resource_type"] = t
					}
					raw, err := c.Request("POST", c.WriterBase+"/documents", &zohttp.RequestOpts{
						Form: form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "details",
				Usage:     "Get document metadata",
				ArgsUsage: "<doc-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WriterBase+"/documents/"+cmd.Args().First(), nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "fields",
				Usage:     "List merge fields in a document",
				ArgsUsage: "<doc-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("GET", c.WriterBase+"/documents/"+cmd.Args().First()+"/fields", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "merge",
				Usage:     "Merge data into a document template",
				ArgsUsage: "<doc-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "json", Required: true, Usage: "Merge data as JSON"},
					&cli.StringFlag{Name: "format", Value: "pdf", Usage: "pdf, pdfform, docx, html, zfdoc, zip (use inline as alias for html)"},
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					docID := cmd.Args().First()

					var mergeData any
					if err := json.Unmarshal([]byte(cmd.String("json")), &mergeData); err != nil {
						return err
					}

					format := cmd.String("format")
					if format == "inline" {
						format = "html"
					}

					outputSettings := map[string]any{
						"format": format,
					}
					outputSettingsJSON, err := json.Marshal(outputSettings)
					if err != nil {
						return err
					}

					mergeJSON, err := json.Marshal(mergeData)
					if err != nil {
						return err
					}

					raw, err := c.Request("POST", c.WriterBase+"/documents/"+docID+"/merge", &zohttp.RequestOpts{
						Form: map[string]string{
							"output_settings": string(outputSettingsJSON),
							"merge_data":      string(mergeJSON),
						},
					})
					if err != nil {
						if strings.Contains(err.Error(), "R3002") {
							return output.JSON(map[string]string{"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
						}
						return err
					}

					body := []byte(raw)
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0644); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}

					os.Stdout.Write(body)
					return nil
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a document",
				ArgsUsage: "<doc-id>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					raw, err := c.Request("DELETE", c.WriterBase+"/documents/"+cmd.Args().First()+"/delete", nil)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:      "read",
				Usage:     "Read document content as text",
				ArgsUsage: "<doc-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Value: "txt", Usage: "txt or html"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw("GET", c.WriterBase+"/download/"+cmd.Args().First(), map[string]string{"format": cmd.String("format")})
					if err != nil {
						if strings.Contains(err.Error(), "R3002") {
							return output.JSON(map[string]string{"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
						}
						return err
					}
					if len(body) == 0 {
						return output.JSON(map[string]string{"error": "Document is empty or could not be read"})
					}
					return output.JSON(map[string]any{
						"document_id": cmd.Args().First(),
						"format":      cmd.String("format"),
						"content":     string(body),
					})
				},
			},
			{
				Name:      "download",
				Usage:     "Download a document",
				ArgsUsage: "<doc-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "format", Value: "docx", Usage: "zdoc, docx, odt, rtf, txt, html, pdf, zip, epub, pdfform"},
					&cli.StringFlag{Name: "output", Usage: "Output file path"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					body, _, _, err := c.RequestRaw("GET", c.WriterBase+"/download/"+cmd.Args().First(), map[string]string{"format": cmd.String("format")})
					if err != nil {
						if strings.Contains(err.Error(), "R3002") {
							return output.JSON(map[string]string{"error": "Document is empty — Zoho cannot export empty documents (R3002)"})
						}
						return err
					}
					if out := cmd.String("output"); out != "" {
						if err := os.WriteFile(out, body, 0644); err != nil {
							return err
						}
						return output.JSON(map[string]any{"ok": true, "path": out, "size": len(body)})
					}
					os.Stdout.Write(body)
					return nil
				},
			},
		},
	}
}
