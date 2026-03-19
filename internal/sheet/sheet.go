package sheet

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/omin8tor/zoho-cli/internal/auth"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/omin8tor/zoho-cli/internal/output"
	"github.com/urfave/cli/v3"
)

var (
	cellRefRe   = regexp.MustCompile(`(?i)^([A-Z]+)([0-9]+)$`)
	rangeRefRe  = regexp.MustCompile(`(?i)^\s*([A-Z]+[0-9]+)\s*:\s*([A-Z]+[0-9]+)\s*$`)
	cellRangeRe = regexp.MustCompile(`(?i)^\s*([A-Z]+)([0-9]+)\s*:\s*([A-Z]+)([0-9]+)\s*$`)
)

func colLettersToIndex(col string) (int, error) {
	col = strings.ToUpper(strings.TrimSpace(col))
	if col == "" {
		return 0, fmt.Errorf("empty column")
	}
	idx := 0
	for _, r := range col {
		if r < 'A' || r > 'Z' {
			return 0, fmt.Errorf("invalid column letter %q", r)
		}
		idx = idx*26 + int(r-'A'+1)
	}
	return idx, nil
}

func parseCellRef(a1 string) (row int, col int, err error) {
	a1 = strings.ToUpper(strings.TrimSpace(a1))
	m := cellRefRe.FindStringSubmatch(a1)
	if len(m) != 3 {
		return 0, 0, fmt.Errorf("invalid cell reference %q (expected like A1)", a1)
	}
	col, err = colLettersToIndex(m[1])
	if err != nil {
		return 0, 0, err
	}
	row, err = strconv.Atoi(m[2])
	if err != nil || row <= 0 {
		return 0, 0, fmt.Errorf("invalid row in cell reference %q", a1)
	}
	return row, col, nil
}

func parseRangeRef(a1range string) (startRow int, startCol int, endRow int, endCol int, err error) {
	a1range = strings.ToUpper(strings.TrimSpace(a1range))
	// Fast path: "A1:B2"
	m := rangeRefRe.FindStringSubmatch(a1range)
	if len(m) == 3 {
		r1, c1, err := parseCellRef(m[1])
		if err != nil {
			return 0, 0, 0, 0, err
		}
		r2, c2, err := parseCellRef(m[2])
		if err != nil {
			return 0, 0, 0, 0, err
		}
		if r1 <= r2 {
			return r1, c1, r2, c2, nil
		}
		return r2, c2, r1, c1, nil
	}

	// Some docs/examples include ranges with no whitespace normalization; keep a fallback.
	m2 := cellRangeRe.FindStringSubmatch(a1range)
	if len(m2) == 5 {
		c1, err := colLettersToIndex(m2[1])
		if err != nil {
			return 0, 0, 0, 0, err
		}
		r1, err := strconv.Atoi(m2[2])
		if err != nil || r1 <= 0 {
			return 0, 0, 0, 0, fmt.Errorf("invalid row in range %q", a1range)
		}
		c2, err := colLettersToIndex(m2[3])
		if err != nil {
			return 0, 0, 0, 0, err
		}
		r2, err := strconv.Atoi(m2[4])
		if err != nil || r2 <= 0 {
			return 0, 0, 0, 0, fmt.Errorf("invalid row in range %q", a1range)
		}
		if r1 <= r2 {
			return r1, c1, r2, c2, nil
		}
		return r2, c2, r1, c1, nil
	}

	return 0, 0, 0, 0, fmt.Errorf("invalid range reference %q (expected like A1:B5)", a1range)
}

func getClient() (*zohttp.Client, error) {
	config, err := auth.ResolveAuth()
	if err != nil {
		return nil, err
	}
	return zohttp.NewClient(config)
}

var workbookFlag = &cli.StringFlag{Name: "workbook", Required: true, Usage: "Workbook resource ID"}

var worksheetFlag = &cli.StringFlag{Name: "worksheet", Usage: "Worksheet name"}

func Commands() *cli.Command {
	return &cli.Command{
		Name:  "sheet",
		Usage: "Zoho Sheet operations",
		Commands: []*cli.Command{
			workbooksCmd(),
			worksheetsCmd(),
			tablesCmd(),
			recordsCmd(),
			cellsCmd(),
			contentCmd(),
			formatCmd(),
			namedRangesCmd(),
			mergeCmd(),
			premiumCmd(),
			utilityCmd(),
		},
	}
}

func workbooksCmd() *cli.Command {
	return &cli.Command{
		Name:  "workbooks",
		Usage: "Workbook operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all workbooks",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "start-index", Usage: "Start index"},
					&cli.IntFlag{Name: "count", Usage: "Number of workbooks"},
					&cli.StringFlag{Name: "sort-option", Usage: "Sort option"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.list"}
					if v := cmd.Int("start-index"); v > 0 {
						params["start_index"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.String("sort-option"); v != "" {
						params["sort_option"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/workbooks", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "templates",
				Usage: "List all templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "template.list"}
					raw, err := c.Request("POST", c.SheetBase+"/templates", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "versions",
				Usage: "List all versions",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.version.list"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create workbook",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "workbook-name", Required: true, Usage: "Workbook name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "workbook.create",
						"workbook_name": cmd.String("workbook-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/workbooks", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-from-template",
				Usage: "Create workbook from template",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "workbook-name", Required: true, Usage: "New workbook name"},
					&cli.StringFlag{Name: "parent-id", Usage: "Parent folder ID for the new workbook"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "workbook.createfromtemplate",
						"resource_id":    cmd.String("workbook"),
						"workbook_name":  cmd.String("workbook-name"),
					}
					if v := cmd.String("parent-id"); v != "" {
						params["parent_id"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/createfromtemplate", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "upload",
				Usage: "Upload workbook",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Required: true, Usage: "Path to file"},
					&cli.StringFlag{Name: "workbook-name", Usage: "Workbook name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					filePath := cmd.String("file")
					data, err := os.ReadFile(filePath)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					name := filepath.Base(filePath)
					params := map[string]string{"method": "workbook.upload"}
					form := map[string]string{}
					if v := cmd.String("workbook-name"); v != "" {
						form["workbook_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/upload", &zohttp.RequestOpts{
						Params: params,
						Files:  map[string]zohttp.FileUpload{"file": {Filename: name, Data: data}},
						Form:   form,
					})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "download",
				Usage: "Download workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "format", Required: true, Usage: "Download format (xlsx/csv/tsv/ods/pdf/html)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "workbook.download",
						"format": cmd.String("format"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/download/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "external-share-link",
				Usage: "Create an external share link",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "link-name", Required: true, Usage: "Link name for easy reference"},
					&cli.StringFlag{Name: "access-level", Required: true, Usage: "Access level: 1.edit | 2.view"},
					&cli.StringFlag{Name: "allow-download", Usage: "Whether shared user can download (true/false)"},
					&cli.StringFlag{Name: "password", Usage: "Optional password for the share link"},
					&cli.StringFlag{Name: "expiration-date", Usage: "Optional expiry date (YYYY-MM-DD)"},
					&cli.StringFlag{Name: "request-user-data", Usage: "Optional requested user data (NAME,PHONE,EMAIL)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}

					params := map[string]string{
						"method":      "workbook.externalsharelink",
						"resource_id": cmd.String("workbook"),
						"link_name":   cmd.String("link-name"),
						"access_level": cmd.String("access-level"),
					}
					if v := cmd.String("allow-download"); v != "" {
						params["allow_download"] = v
					}
					if v := cmd.String("password"); v != "" {
						params["password"] = v
					}
					if v := cmd.String("expiration-date"); v != "" {
						params["expiration_date"] = v
					}
					if v := cmd.String("request-user-data"); v != "" {
						params["request_user_data"] = v
					}

					raw, err := c.Request("POST", c.SheetBase+"/externalsharelink", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-images",
				Usage: "Insert images into workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "image-json", Required: true, Usage: "Image JSON configuration"},
					&cli.StringFlag{Name: "file", Usage: "Path to image file"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "workbook.images.insert",
						"image_json": cmd.String("image-json"),
					}
					opts := &zohttp.RequestOpts{Params: params}
					if filePath := cmd.String("file"); filePath != "" {
						data, err := os.ReadFile(filePath)
						if err != nil {
							return fmt.Errorf("failed to read file: %w", err)
						}
						name := filepath.Base(filePath)
						opts.Files = map[string]zohttp.FileUpload{"imagefiles": {Filename: name, Data: data}}
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), opts)
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy",
				Usage: "Copy workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "new-workbook-name", Required: true, Usage: "New workbook name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":            "workbook.copy",
						"new_workbook_name": cmd.String("new-workbook-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "share",
				Usage: "Share workbook",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "email", Required: true, Usage: "Email address to share with"},
					&cli.StringFlag{Name: "role", Usage: "Role for shared user"},
					&cli.StringFlag{Name: "notify", Usage: "Notify user (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":   "workbook.share",
						"email_id": cmd.String("email"),
					}
					if v := cmd.String("role"); v != "" {
						params["role"] = v
					}
					if v := cmd.String("notify"); v != "" {
						params["notify"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create-version",
				Usage: "Create a version",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "version-name", Required: true, Usage: "Version name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "workbook.version.create",
						"version_name": cmd.String("version-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "revert-version",
				Usage: "Revert to a version",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "version-id", Required: true, Usage: "Version ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "workbook.version.revert",
						"version_id": cmd.String("version-id"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "trash",
				Usage: "Trash workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.trash"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "restore",
				Usage: "Restore workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.restore"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.delete"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "publish",
				Usage: "Publish workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.publish"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "unpublish",
				Usage: "Remove publish from workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.publish.remove"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "lock",
				Usage: "Lock worksheet or range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "scope", Required: true, Usage: "Lock scope: worksheet|range"},
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Usage: "Start row index (range scope only)"},
					&cli.IntFlag{Name: "start-column", Usage: "Start column index (range scope only)"},
					&cli.IntFlag{Name: "end-row", Usage: "End row index (range scope only)"},
					&cli.IntFlag{Name: "end-column", Usage: "End column index (range scope only)"},
					&cli.StringFlag{Name: "user-emails", Usage: "JSON array of user email ids (at least one of user-emails/external-share-links required)"},
					&cli.StringFlag{Name: "external-share-links", Usage: "JSON array of external share links (at least one of user-emails/external-share-links required)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}

					userEmails := cmd.String("user-emails")
					externalShareLinks := cmd.String("external-share-links")
					if userEmails == "" && externalShareLinks == "" {
						return fmt.Errorf("at least one of --user-emails or --external-share-links is required")
					}

					params := map[string]string{
						"method": "lock",
						"scope":  cmd.String("scope"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-column"); v > 0 {
						params["start_column"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-row"); v > 0 {
						params["end_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-column"); v > 0 {
						params["end_column"] = fmt.Sprintf("%d", v)
					}
					if userEmails != "" {
						params["user_emails"] = userEmails
					}
					if externalShareLinks != "" {
						params["external_share_links"] = externalShareLinks
					}

					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "unlock",
				Usage: "Unlock worksheet or range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "scope", Required: true, Usage: "Unlock scope: worksheet|range"},
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Usage: "Start row index (range scope only)"},
					&cli.IntFlag{Name: "start-column", Usage: "Start column index (range scope only)"},
					&cli.IntFlag{Name: "end-row", Usage: "End row index (range scope only)"},
					&cli.IntFlag{Name: "end-column", Usage: "End column index (range scope only)"},
					&cli.StringFlag{Name: "user-emails", Usage: "JSON array of user email ids (at least one of user-emails/external-share-links required)"},
					&cli.StringFlag{Name: "external-share-links", Usage: "JSON array of external share links (at least one of user-emails/external-share-links required)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					userEmails := cmd.String("user-emails")
					externalShareLinks := cmd.String("external-share-links")
					if userEmails == "" && externalShareLinks == "" {
						return fmt.Errorf("at least one of --user-emails or --external-share-links is required")
					}

					params := map[string]string{
						"method": "unlock",
						"scope":  cmd.String("scope"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-column"); v > 0 {
						params["start_column"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-row"); v > 0 {
						params["end_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-column"); v > 0 {
						params["end_column"] = fmt.Sprintf("%d", v)
					}
					if userEmails != "" {
						params["user_emails"] = userEmails
					}
					if externalShareLinks != "" {
						params["external_share_links"] = externalShareLinks
					}

					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func worksheetsCmd() *cli.Command {
	return &cli.Command{
		Name:  "worksheets",
		Usage: "Worksheet operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all worksheets",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.list"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "worksheet-name", Usage: "Worksheet name (optional; default Sheet1, Sheet2, ...)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.insert"}
					if v := cmd.String("worksheet-name"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy",
				Usage: "Copy worksheet within same workbook",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "new-worksheet-name", Required: true, Usage: "New worksheet name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":             "worksheet.copy",
						"worksheet_name":     cmd.String("worksheet"),
						"new_worksheet_name": cmd.String("new-worksheet-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy-to",
				Usage: "Copy worksheet from another workbook",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "source-workbook", Required: true, Usage: "Source workbook resource ID"},
					&cli.StringFlag{Name: "source-worksheet-id", Usage: "Source worksheet ID (optional)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":                "worksheet.copy.otherdoc",
						"source_resource_id":   cmd.String("source-workbook"),
						"source_worksheet_name": cmd.String("worksheet"),
					}
					if v := cmd.String("source-worksheet-id"); v != "" {
						params["source_worksheet_id"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "copy-create-new",
				Usage: "Copy worksheet and create a new workbook",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "workbook-name", Required: true, Usage: "New workbook name"},
					&cli.StringFlag{Name: "include-referred-sheets", Usage: "Include referred sheets (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "worksheet.copy.createnew",
						"worksheet_name": cmd.String("worksheet"),
						"workbook_name": cmd.String("workbook-name"),
					}
					if v := cmd.String("include-referred-sheets"); v != "" {
						params["include_referred_sheets"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "freeze",
				Usage: "Freeze worksheet panes",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "freeze-type", Required: true, Usage: "freeze_firstrow|freeze_firstcolumn|freeze_rows|freeze_columns|freeze_panes|unfreeze_panes"},
					&cli.IntFlag{Name: "start-row", Usage: "Start row index (freeze_rows/freeze_panes)"},
					&cli.IntFlag{Name: "start-column", Usage: "Start column index (freeze_columns/freeze_panes)"},
					&cli.IntFlag{Name: "end-row", Usage: "End row index (freeze_rows/freeze_panes)"},
					&cli.IntFlag{Name: "end-column", Usage: "End column index (freeze_columns/freeze_panes)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "worksheet.freeze",
						"freeze_type": cmd.String("freeze-type"),
						"worksheet_name": cmd.String("worksheet"),
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-column"); v > 0 {
						params["start_column"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-row"); v > 0 {
						params["end_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-column"); v > 0 {
						params["end_column"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "rename",
				Usage: "Rename worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "new-worksheet-name", Required: true, Usage: "New worksheet name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":             "worksheet.rename",
						"worksheet_name":     cmd.String("worksheet"),
						"new_worksheet_name": cmd.String("new-worksheet-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "worksheet.delete",
						"worksheet_name": cmd.String("worksheet"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-multiple",
				Usage: "Delete multiple worksheets",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "worksheet-names", Usage: "JSON array of worksheet names"},
					&cli.StringFlag{Name: "worksheet-ids", Usage: "JSON array of worksheet ids"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					names := cmd.String("worksheet-names")
					ids := cmd.String("worksheet-ids")
					if names == "" && ids == "" {
						return fmt.Errorf("at least one of --worksheet-names or --worksheet-ids is required")
					}
					params := map[string]string{"method": "worksheets.delete"}
					if names != "" {
						params["worksheet_names"] = names
					}
					if ids != "" {
						params["worksheet_ids"] = ids
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func tablesCmd() *cli.Command {
	return &cli.Command{
		Name:  "tables",
		Usage: "Table operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all tables",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "table.list"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create table",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row"},
					&cli.IntFlag{Name: "end-column", Required: true, Usage: "End column"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "table.create",
						"table_name":   cmd.String("table-name"),
						"start_row":    fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("start-column")),
						"end_row":      fmt.Sprintf("%d", cmd.Int("end-row")),
						"end_column":   fmt.Sprintf("%d", cmd.Int("end-column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "remove",
				Usage: "Remove table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.remove",
						"table_name": cmd.String("table-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "rename-headers",
				Usage: "Rename headers of table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Header rename data JSON array"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.header.rename",
						"table_name": cmd.String("table-name"),
						"data":       cmd.String("data"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fetch-records",
				Usage: "Fetch records from table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "criteria", Usage: "Filter criteria"},
					&cli.IntFlag{Name: "start-index", Usage: "Start index"},
					&cli.IntFlag{Name: "count", Usage: "Number of records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.fetch",
						"table_name": cmd.String("table-name"),
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.Int("start-index"); v > 0 {
						params["start_index"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-records",
				Usage: "Add records to table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.add",
						"table_name": cmd.String("table-name"),
						"json_data":  cmd.String("json"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-records",
				Usage: "Update records in table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.update",
						"table_name": cmd.String("table-name"),
						"criteria":   cmd.String("criteria"),
						"json_data":  cmd.String("json"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-records",
				Usage: "Delete records from table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "delete-rows", Usage: "Delete entire rows (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.records.delete",
						"table_name": cmd.String("table-name"),
						"criteria":   cmd.String("criteria"),
					}
					if v := cmd.String("delete-rows"); v != "" {
						params["delete_rows"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-columns",
				Usage: "Insert columns to table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "columns", Required: true, Usage: "Columns JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.columns.insert",
						"table_name": cmd.String("table-name"),
						"columns":    cmd.String("columns"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-columns",
				Usage: "Delete columns from table",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "table-name", Required: true, Usage: "Table name"},
					&cli.StringFlag{Name: "columns", Required: true, Usage: "Columns JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "table.columns.delete",
						"table_name": cmd.String("table-name"),
						"columns":    cmd.String("columns"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func recordsCmd() *cli.Command {
	return &cli.Command{
		Name:  "records",
		Usage: "Worksheet record operations",
		Commands: []*cli.Command{
			{
				Name:  "fetch",
				Usage: "Fetch records from worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.IntFlag{Name: "start-row", Usage: "Start row number"},
					&cli.IntFlag{Name: "count", Usage: "Number of records"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.records.fetch"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("count"); v > 0 {
						params["count"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add",
				Usage: "Add records to worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.records.add",
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update records in worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.records.update",
						"criteria":  cmd.String("criteria"),
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete records from worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "delete-rows", Usage: "Delete entire rows (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":   "worksheet.records.delete",
						"criteria": cmd.String("criteria"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.String("delete-rows"); v != "" {
						params["delete_rows"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-columns",
				Usage: "Insert columns to worksheet records",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "insert-column-after", Usage: "Insert columns after this header name"},
					&cli.StringFlag{Name: "column-names", Required: true, Usage: "JSON array of new column header names"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "records.columns.insert",
						"column_names":  cmd.String("column-names"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.String("insert-column-after"); v != "" {
						params["insert_column_after"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func cellsCmd() *cli.Command {
	return &cli.Command{
		Name:  "cells",
		Usage: "Cell and range content operations",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get content of cell",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "cell", Required: true, Usage: "Cell reference (e.g. A1)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					row, col, err := parseCellRef(cmd.String("cell"))
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "cell.content.get",
						"row":    fmt.Sprintf("%d", row),
						"column": fmt.Sprintf("%d", col),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-range",
				Usage: "Get content of range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference (e.g. A1:B5)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					startRow, startCol, endRow, endCol, err := parseRangeRef(cmd.String("range"))
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.content.get",
						"start_row":   fmt.Sprintf("%d", startRow),
						"start_column": fmt.Sprintf("%d", startCol),
						"end_row":     fmt.Sprintf("%d", endRow),
						"end_column":  fmt.Sprintf("%d", endCol),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-named-range",
				Usage: "Get content of named range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "named-range", Required: true, Usage: "Named range"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "namedrange.content.get",
						"name_of_range": cmd.String("named-range"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-worksheet",
				Usage: "Get content of worksheet area",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Usage: "Start row"},
					&cli.IntFlag{Name: "start-column", Usage: "Start column"},
					&cli.IntFlag{Name: "end-row", Usage: "End row"},
					&cli.IntFlag{Name: "end-column", Usage: "End column"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.content.get"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("start-row"); v > 0 {
						params["start_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("start-column"); v > 0 {
						params["start_column"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-row"); v > 0 {
						params["end_row"] = fmt.Sprintf("%d", v)
					}
					if v := cmd.Int("end-column"); v > 0 {
						params["end_column"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "get-used-area",
				Usage: "Get used area of worksheet",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.usedarea"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set",
				Usage: "Set content to cell",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "cell", Required: true, Usage: "Cell reference (e.g. A1)"},
					&cli.StringFlag{Name: "value", Required: true, Usage: "Cell value"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					row, col, err := parseCellRef(cmd.String("cell"))
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":  "cell.content.set",
						"row":     fmt.Sprintf("%d", row),
						"column":  fmt.Sprintf("%d", col),
						"content": cmd.String("value"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-multiple",
				Usage: "Set content to multiple cells",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "cell-data", Required: true, Usage: "Cell data JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "cells.content.set",
						"data":   cmd.String("cell-data"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-row",
				Usage: "Set content to row",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
					&cli.StringFlag{Name: "column-array", Required: true, Usage: "Column array JSON"},
					&cli.StringFlag{Name: "data-array", Required: true, Usage: "Data array JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "row.content.set",
						"row":          fmt.Sprintf("%d", cmd.Int("row")),
						"column_array": cmd.String("column-array"),
						"data_array":   cmd.String("data-array"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-range",
				Usage: "Set content to range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference (e.g. A1:B5)"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "Data JSON 2D array"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "range.set",
						"range":  cmd.String("range"),
						"data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func contentCmd() *cli.Command {
	return &cli.Command{
		Name:  "content",
		Usage: "Content operations",
		Commands: []*cli.Command{
			{
				Name:  "append-csv",
				Usage: "Append rows with CSV data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "data", Required: true, Usage: "CSV data string"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "worksheet.csvdata.append",
						"data":   cmd.String("data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-csv",
				Usage: "Set rows with CSV data (worksheet.csvdata.set)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Start row index"},
					&cli.IntFlag{Name: "column", Required: true, Usage: "Start column index"},
					&cli.StringFlag{Name: "ignore-empty", Usage: "Whether to ignore empty cells (true/false)"},
					&cli.StringFlag{Name: "data", Required: true, Usage: "CSV data string"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "worksheet.csvdata.set",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
						"column": fmt.Sprintf("%d", cmd.Int("column")),
						"data":   cmd.String("data"),
					}
					if v := cmd.String("ignore-empty"); v != "" {
						params["ignore_empty"] = v
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "append-json",
				Usage: "Append rows with JSON data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.jsondata.append",
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-json",
				Usage: "Update rows with JSON data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "worksheet.jsondata.set",
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-json",
				Usage: "Insert row with JSON data",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.IntFlag{Name: "row-index", Required: true, Usage: "Row index to insert at"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "worksheet.jsondata.insert",
						"row":    fmt.Sprintf("%d", cmd.Int("row-index")),
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "clear-contents",
				Usage: "Clear contents of range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					startRow, startCol, endRow, endCol, err := parseRangeRef(cmd.String("range"))
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.content.clear",
						"start_row":   fmt.Sprintf("%d", startRow),
						"start_column": fmt.Sprintf("%d", startCol),
						"end_row":     fmt.Sprintf("%d", endRow),
						"end_column":  fmt.Sprintf("%d", endCol),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "clear-range",
				Usage: "Clear range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					startRow, startCol, endRow, endCol, err := parseRangeRef(cmd.String("range"))
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":       "range.clear",
						"start_row":   fmt.Sprintf("%d", startRow),
						"start_column": fmt.Sprintf("%d", startCol),
						"end_row":     fmt.Sprintf("%d", endRow),
						"end_column":  fmt.Sprintf("%d", endCol),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "clear-filters",
				Usage: "Clear filters",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "worksheet.filter.clear"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "find",
				Usage: "Find content",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "search-value", Required: true, Usage: "Value to search for"},
					&cli.StringFlag{Name: "match-type", Usage: "Match type"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "find",
						"search": cmd.String("search-value"),
						"scope":  "worksheet",
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := strings.TrimSpace(strings.ToLower(cmd.String("match-type"))); v != "" {
						// best-effort mapping: pass-through boolean for exact match
						if v == "true" || v == "false" {
							params["is_exact_match"] = v
						}
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "find-replace",
				Usage: "Find and replace content",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "search-value", Required: true, Usage: "Value to search for"},
					&cli.StringFlag{Name: "replace-value", Required: true, Usage: "Replacement value"},
					&cli.StringFlag{Name: "match-type", Usage: "Match type"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "replace",
						"search":        cmd.String("search-value"),
						"replace_with": cmd.String("replace-value"),
						"scope":         "worksheet",
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := strings.TrimSpace(strings.ToLower(cmd.String("match-type"))); v != "" {
						if v == "true" || v == "false" {
							params["is_exact_match"] = v
						}
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "recalculate",
				Usage: "Recalculate workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "recalculate"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func formatCmd() *cli.Command {
	return &cli.Command{
		Name:  "format",
		Usage: "Formatting and structure operations",
		Commands: []*cli.Command{
			{
				Name:  "ranges",
				Usage: "Format ranges",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "format-json", Required: true, Usage: "Format JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "ranges.format.set",
						"format_json": cmd.String("format-json"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "image-fit",
				Usage: "Image fit options",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row index"},
					&cli.IntFlag{Name: "start-column", Required: true, Usage: "Start column index"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row index"},
					&cli.IntFlag{Name: "end-column", Required: true, Usage: "End column index"},
					&cli.StringFlag{Name: "image-fit-option", Required: true, Usage: "fit|stretch|cover"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "range.images.fit",
						"start_row":      fmt.Sprintf("%d", cmd.Int("start-row")),
						"start_column":   fmt.Sprintf("%d", cmd.Int("start-column")),
						"end_row":        fmt.Sprintf("%d", cmd.Int("end-row")),
						"end_column":     fmt.Sprintf("%d", cmd.Int("end-column")),
						"image_fit_option": cmd.String("image-fit-option"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "row-height",
				Usage: "Set row height",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
					&cli.IntFlag{Name: "height", Required: true, Usage: "Height in pixels"},
					&cli.StringFlag{Name: "auto-fit", Usage: "Auto fit (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "worksheet.rows.height",
						"row_index_array": fmt.Sprintf(`[{"start_row":%d,"end_row":%d}]`, cmd.Int("row"), cmd.Int("row")),
						"row_height":    fmt.Sprintf("%d", cmd.Int("height")),
					}
					if v := cmd.String("auto-fit"); v != "" {
						params["auto_fit"] = v
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "column-width",
				Usage: "Set column width",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column number"},
					&cli.IntFlag{Name: "width", Required: true, Usage: "Width in pixels"},
					&cli.StringFlag{Name: "auto-fit", Usage: "Auto fit (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":            "worksheet.columns.width",
						"column_index_array": fmt.Sprintf(`[{"start_column":%d,"end_column":%d}]`, cmd.Int("column"), cmd.Int("column")),
						"column_width":     fmt.Sprintf("%d", cmd.Int("width")),
					}
					if v := cmd.String("auto-fit"); v != "" {
						params["auto_fit"] = v
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-row",
				Usage: "Insert row",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
					&cli.StringFlag{Name: "json-data", Required: true, Usage: "JSON array for the inserted row"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "row.insert",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
						"json_data": cmd.String("json-data"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "insert-column",
				Usage: "Insert column",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column number"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "column.insert",
						"column": fmt.Sprintf("%d", cmd.Int("column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-row",
				Usage: "Delete row",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "row.delete",
						"row":    fmt.Sprintf("%d", cmd.Int("row")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-rows",
				Usage: "Delete multiple rows",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "start-row", Required: true, Usage: "Start row"},
					&cli.IntFlag{Name: "end-row", Required: true, Usage: "End row"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":            "worksheet.rows.delete",
						"row_index_array":  fmt.Sprintf(`[{"start_row":%d,"end_row":%d}]`, cmd.Int("start-row"), cmd.Int("end-row")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete-column",
				Usage: "Delete column",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column number"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "column.delete",
						"column": fmt.Sprintf("%d", cmd.Int("column")),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "set-note",
				Usage: "Set note to cell",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "cell", Required: true, Usage: "Cell reference (e.g. A1)"},
					&cli.StringFlag{Name: "note", Required: true, Usage: "Note text"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					row, col, err := parseCellRef(cmd.String("cell"))
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "cell.note.set",
						"row":     fmt.Sprintf("%d", row),
						"column":  fmt.Sprintf("%d", col),
						"note":    cmd.String("note"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func namedRangesCmd() *cli.Command {
	return &cli.Command{
		Name:  "named-ranges",
		Usage: "Named range operations",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all named ranges",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "namedrange.list"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "create",
				Usage: "Create named range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Named range name"},
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "namedrange.create",
						"name":   cmd.String("name"),
						"range":  cmd.String("range"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update",
				Usage: "Update named range",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Named range name"},
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "namedrange.update",
						"name":   cmd.String("name"),
						"range":  cmd.String("range"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete named range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "name", Required: true, Usage: "Named range name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method": "namedrange.delete",
						"name":   cmd.String("name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func mergeCmd() *cli.Command {
	return &cli.Command{
		Name:  "merge",
		Usage: "Merge template operations",
		Commands: []*cli.Command{
			{
				Name:  "templates",
				Usage: "Get merge templates",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "mergetemplate.list"}
					raw, err := c.Request("POST", c.SheetBase+"/mergetemplates", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fields",
				Usage: "Get merge fields",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.mergefield.list"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "jobs",
				Usage: "Get merge jobs",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "workbook.mergejob.list"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "job-detail",
				Usage: "Get merge job details",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "job-id", Required: true, Usage: "Job ID"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "workbook.mergejob.details",
						"mergejob_id": cmd.String("job-id"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "save",
				Usage: "Merge and save",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "merge-data", Required: true, Usage: "Merge data JSON"},
					&cli.StringFlag{Name: "output-settings", Required: true, Usage: "Output settings JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":          "merge.save",
						"merge_data":      cmd.String("merge-data"),
						"output_settings": cmd.String("output-settings"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "email",
				Usage: "Merge and email",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "merge-data", Required: true, Usage: "Merge data JSON"},
					&cli.StringFlag{Name: "email-settings", Required: true, Usage: "Email settings JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":         "merge.email.attachment",
						"merge_data":     cmd.String("merge-data"),
						"email_settings": cmd.String("email-settings"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "field-list",
				Usage: "List fields of a workbook",
				Flags: []cli.Flag{
					workbookFlag,
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "field.list"}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "field-create",
				Usage: "Create a merge field",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "field-names", Required: true, Usage: "JSON array of field names"},
					&cli.StringFlag{Name: "cells", Required: true, Usage: "Array of discontiguous ranges (JSON)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "field.create",
						"field_names": cmd.String("field-names"),
						"cells":       cmd.String("cells"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "field-update",
				Usage: "Update a merge field",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "field-name", Required: true, Usage: "Existing field name"},
					&cli.StringFlag{Name: "cells", Required: true, Usage: "Array of discontiguous ranges (JSON)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "field.update",
						"field_name": cmd.String("field-name"),
						"cells":      cmd.String("cells"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "field-delete",
				Usage: "Delete a merge field",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "field-name", Required: true, Usage: "Existing field name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":     "field.delete",
						"field_name": cmd.String("field-name"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fields-create-pdfs",
				Usage: "Create PDFs from fields",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "fields-data", Required: true, Usage: "Fields data JSON array"},
					&cli.StringFlag{Name: "pdf-file-name", Usage: "Optional pdf file name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "fields.create.pdfs",
						"fields_data": cmd.String("fields-data"),
					}
					if v := cmd.String("pdf-file-name"); v != "" {
						params["pdf_file_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fields-create-workbooks",
				Usage: "Create workbooks from fields",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "fields-data", Required: true, Usage: "Fields data JSON array"},
					&cli.StringFlag{Name: "new-workbook-name", Usage: "Optional new workbook name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "fields.create.workbooks",
						"fields_data": cmd.String("fields-data"),
					}
					if v := cmd.String("new-workbook-name"); v != "" {
						params["new_workbook_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fields-create-worksheets",
				Usage: "Create worksheets from fields",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "worksheet-name", Usage: "Target worksheet name"},
					&cli.StringFlag{Name: "worksheet-id", Usage: "Alternatively target worksheet ID"},
					&cli.StringFlag{Name: "fields-data", Required: true, Usage: "Fields data JSON array"},
					&cli.StringFlag{Name: "new-worksheet-name", Usage: "Optional new worksheet name"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "fields.create.worksheets",
						"fields_data": cmd.String("fields-data"),
					}
					if v := cmd.String("worksheet-name"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.String("worksheet-id"); v != "" {
						params["worksheet_id"] = v
					}
					if v := cmd.String("new-worksheet-name"); v != "" {
						params["new_worksheet_name"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "fields-mail-attachment",
				Usage: "Send field-based mail attachments",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "fields-data", Required: true, Usage: "Fields data JSON array"},
					&cli.StringFlag{Name: "file-format", Usage: "Optional output file format"},
					&cli.StringFlag{Name: "recipients", Required: true, Usage: "Recipients email ids (JSON array or comma-separated string)"},
					&cli.StringFlag{Name: "subject", Required: true, Usage: "Email subject"},
					&cli.StringFlag{Name: "message", Required: true, Usage: "Email message/body"},
					&cli.StringFlag{Name: "send-me-a-copy", Usage: "Optional send me a copy (true/false)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "fields.mail.attachment",
						"fields_data": cmd.String("fields-data"),
						"recipients":  cmd.String("recipients"),
						"subject":     cmd.String("subject"),
						"message":     cmd.String("message"),
					}
					if v := cmd.String("file-format"); v != "" {
						params["file_format"] = v
					}
					if v := cmd.String("send-me-a-copy"); v != "" {
						params["send_me_a_copy"] = v
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func premiumCmd() *cli.Command {
	return &cli.Command{
		Name:  "premium",
		Usage: "Premium API operations",
		Commands: []*cli.Command{
			{
				Name:  "fetch-records",
				Usage: "Fetch records (premium)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.StringFlag{Name: "criteria", Usage: "Filter criteria"},
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{"method": "premium.records.fetch"}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.String("criteria"); v != "" {
						params["criteria"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "add-records",
				Usage: "Add records (premium)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "premium.records.add",
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "update-records",
				Usage: "Update records (premium)",
				Flags: []cli.Flag{
					workbookFlag,
					worksheetFlag,
					&cli.IntFlag{Name: "header-row", Usage: "Header row number"},
					&cli.StringFlag{Name: "criteria", Required: true, Usage: "Filter criteria"},
					&cli.StringFlag{Name: "json", Required: true, Usage: "JSON data"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":    "premium.records.update",
						"criteria":  cmd.String("criteria"),
						"json_data": cmd.String("json"),
					}
					if v := cmd.String("worksheet"); v != "" {
						params["worksheet_name"] = v
					}
					if v := cmd.Int("header-row"); v > 0 {
						params["header_row"] = fmt.Sprintf("%d", v)
					}
					raw, err := c.Request("POST", c.SheetBase+"/"+cmd.String("workbook"), &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}

func utilityCmd() *cli.Command {
	return &cli.Command{
		Name:  "utility",
		Usage: "Utility operations",
		Commands: []*cli.Command{
			{
				Name:  "range-to-index",
				Usage: "Convert range to index",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.StringFlag{Name: "range", Required: true, Usage: "Range reference"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":        "range.index.get",
						"range_address": cmd.String("range"),
					}
					raw, err := c.Request("POST", c.SheetBase+"/utils", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
			{
				Name:  "index-to-range",
				Usage: "Convert index to range",
				Flags: []cli.Flag{
					workbookFlag,
					&cli.IntFlag{Name: "row", Required: true, Usage: "Row number"},
					&cli.IntFlag{Name: "column", Required: true, Usage: "Column number"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					params := map[string]string{
						"method":      "range.address.get",
						"start_row":   fmt.Sprintf("%d", cmd.Int("row")),
						"start_column": fmt.Sprintf("%d", cmd.Int("column")),
						"end_row":     fmt.Sprintf("%d", cmd.Int("row")),
						"end_column":  fmt.Sprintf("%d", cmd.Int("column")),
					}
					raw, err := c.Request("POST", c.SheetBase+"/utils", &zohttp.RequestOpts{Params: params})
					if err != nil {
						return err
					}
					return output.JSONRaw(raw)
				},
			},
		},
	}
}
