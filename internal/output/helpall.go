package output

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func PrintHelpAll(app *cli.Command) error {
	origPrinter := cli.HelpPrinter
	cli.HelpPrinter = filterHelpPrinter
	defer func() { cli.HelpPrinter = origPrinter }()

	printHelp(app)
	walkCommands(app)
	return nil
}

func walkCommands(cmd *cli.Command) {
	for _, sub := range cmd.Commands {
		if sub.Name == "help" {
			continue
		}
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "────────────────────────────────────────")
		printHelp(sub)
		walkCommands(sub)
	}
}

func printHelp(cmd *cli.Command) {
	cmd.Run(context.Background(), []string{cmd.Name, "--help"})
}

func filterHelpPrinter(w io.Writer, templ string, data any) {
	var buf strings.Builder
	cli.DefaultPrintHelp(&buf, templ, data)
	raw := buf.String()

	lines := strings.Split(raw, "\n")
	var out []string
	skip := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "GLOBAL OPTIONS:" {
			skip = true
			continue
		}
		if skip {
			if trimmed == "" || strings.HasPrefix(trimmed, "--") {
				continue
			}
			skip = false
		}

		if strings.HasPrefix(trimmed, "--help,") && strings.HasSuffix(trimmed, "show help") {
			continue
		}
		if strings.HasPrefix(trimmed, "help,") && strings.Contains(trimmed, "Shows a list of commands") {
			continue
		}

		out = append(out, line)
	}

	var cleaned []string
	for i, line := range out {
		trimmed := strings.TrimSpace(line)
		if trimmed == "OPTIONS:" || trimmed == "COMMANDS:" {
			rest := ""
			for j := i + 1; j < len(out); j++ {
				rest = strings.TrimSpace(out[j])
				if rest != "" {
					break
				}
			}
			if rest == "" || strings.HasSuffix(rest, ":") {
				continue
			}
		}
		cleaned = append(cleaned, line)
	}

	for len(cleaned) > 0 && strings.TrimSpace(cleaned[len(cleaned)-1]) == "" {
		cleaned = cleaned[:len(cleaned)-1]
	}

	fmt.Fprintln(w, strings.Join(cleaned, "\n"))
}
