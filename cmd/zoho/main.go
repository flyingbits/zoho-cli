package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/flyingbits/zoho-cli/internal"
	"github.com/flyingbits/zoho-cli/internal/auth"
	"github.com/flyingbits/zoho-cli/internal/bigin"
	"github.com/flyingbits/zoho-cli/internal/billing"
	"github.com/flyingbits/zoho-cli/internal/books"
	"github.com/flyingbits/zoho-cli/internal/cliq"
	"github.com/flyingbits/zoho-cli/internal/creator"
	"github.com/flyingbits/zoho-cli/internal/crm"
	"github.com/flyingbits/zoho-cli/internal/desk"
	"github.com/flyingbits/zoho-cli/internal/drive"
	"github.com/flyingbits/zoho-cli/internal/expense"
	"github.com/flyingbits/zoho-cli/internal/inventory"
	"github.com/flyingbits/zoho-cli/internal/invoice"
	"github.com/flyingbits/zoho-cli/internal/mail"
	"github.com/flyingbits/zoho-cli/internal/output"
	"github.com/flyingbits/zoho-cli/internal/people"
	"github.com/flyingbits/zoho-cli/internal/projects"
	"github.com/flyingbits/zoho-cli/internal/recruit"
	"github.com/flyingbits/zoho-cli/internal/sheet"
	"github.com/flyingbits/zoho-cli/internal/sign"
	"github.com/flyingbits/zoho-cli/internal/sprints"
	"github.com/flyingbits/zoho-cli/internal/writer"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	app := &cli.Command{
		Name:    "zoho",
		Usage:   "CLI for Zoho REST APIs",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "help-all", Usage: "Show help for all commands recursively"},
		},
		Commands: []*cli.Command{
			auth.Commands(),
			bigin.Commands(),
			billing.Commands(),
			books.Commands(),
			cliq.Commands(),
			creator.Commands(),
			crm.Commands(),
			desk.Commands(),
			drive.Commands(),
			expense.Commands(),
			inventory.Commands(),
			invoice.Commands(),
			mail.Commands(),
			people.Commands(),
			projects.Commands(),
			recruit.Commands(),
			sheet.Commands(),
			sign.Commands(),
			sprints.Commands(),
			writer.Commands(),
		},
	}

	app.Action = func(_ context.Context, cmd *cli.Command) error {
		if cmd.Bool("help-all") {
			return output.PrintHelpAll(cmd)
		}
		cmd.Root().Run(context.Background(), []string{"zoho", "--help"})
		return nil
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		var apiErr *internal.ZohoAPIError
		var cliErr *internal.ZohoCliError
		if errors.As(err, &apiErr) {
			fmt.Fprintln(os.Stderr, apiErr.Message)
			os.Exit(apiErr.ExitCode)
		}
		if errors.As(err, &cliErr) {
			fmt.Fprintln(os.Stderr, cliErr.Message)
			os.Exit(cliErr.ExitCode)
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
