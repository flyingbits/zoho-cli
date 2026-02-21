package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

type CommandSchema struct {
	Command     string          `json:"command"`
	Usage       string          `json:"usage"`
	ArgsUsage   string          `json:"args_usage,omitempty"`
	Flags       []FlagSchema    `json:"flags,omitempty"`
	Subcommands []CommandSchema `json:"subcommands,omitempty"`
}

type FlagSchema struct {
	Name     string `json:"name"`
	Usage    string `json:"usage"`
	Required bool   `json:"required,omitempty"`
	Default  string `json:"default,omitempty"`
}

func CollectSchemas(cmd *cli.Command, prefix string) []CommandSchema {
	var schemas []CommandSchema
	for _, sub := range cmd.Commands {
		if sub.Name == "help" {
			continue
		}
		fullName := sub.Name
		if prefix != "" {
			fullName = prefix + " " + sub.Name
		}
		userCmds := 0
		for _, c := range sub.Commands {
			if c.Name != "help" {
				userCmds++
			}
		}
		if userCmds > 0 {
			schemas = append(schemas, CollectSchemas(sub, fullName)...)
		} else {
			schema := CommandSchema{
				Command:   fullName,
				Usage:     sub.Usage,
				ArgsUsage: sub.ArgsUsage,
			}
			for _, f := range sub.Flags {
				names := f.Names()
				if len(names) == 0 {
					continue
				}
				isHelp := false
				for _, n := range names {
					if n == "help" || n == "h" {
						isHelp = true
					}
				}
				if isHelp {
					continue
				}
				name := names[0]
				for _, n := range names {
					if len(n) > len(name) {
						name = n
					}
				}
				fs := FlagSchema{Name: "--" + name}
				switch ff := f.(type) {
				case *cli.StringFlag:
					fs.Usage = ff.Usage
					fs.Required = ff.Required
					if ff.Value != "" {
						fs.Default = ff.Value
					}
				case *cli.IntFlag:
					fs.Usage = ff.Usage
					fs.Required = ff.Required
					if ff.Value != 0 {
						fs.Default = fmt.Sprintf("%d", ff.Value)
					}
				case *cli.BoolFlag:
					fs.Usage = ff.Usage
					fs.Required = ff.Required
				}
				schema.Flags = append(schema.Flags, fs)
			}
			schemas = append(schemas, schema)
		}
	}
	return schemas
}

func PrintHelpAll(app *cli.Command) error {
	schemas := CollectSchemas(app, "")
	grouped := map[string][]CommandSchema{}
	for _, s := range schemas {
		parts := strings.SplitN(s.Command, " ", 2)
		group := parts[0]
		grouped[group] = append(grouped[group], s)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(grouped)
}
