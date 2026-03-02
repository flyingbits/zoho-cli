package mail

import (
	"context"
	"encoding/json"

	"github.com/omin8tor/zoho-cli/internal"
	zohttp "github.com/omin8tor/zoho-cli/internal/http"
	"github.com/urfave/cli/v3"
)

func policyCmd() *cli.Command {
	return &cli.Command{
		Name:  "policy",
		Usage: "Mail Policy API",
		Commands: []*cli.Command{
			{
				Name:  "create-org",
				Usage: "Create org policy",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/policy", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-email-restriction",
				Usage: "Create email restriction",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/policy", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-account-restriction",
				Usage: "Create account restriction",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/policy", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-access-restriction",
				Usage: "Create access restriction",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/policy", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:  "create-forward-restriction",
				Usage: "Create forward restriction",
				Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, err := mustOrg(cmd)
					if err != nil {
						return err
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "POST", "/organization/"+zoid+"/policy", &zohttp.RequestOpts{JSON: body})
				},
			},
			{
				Name:      "list",
				Usage:     "Get all policies",
				ArgsUsage: "[zoid]",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/policy", nil)
				},
			},
			{
				Name:      "email-restrictions",
				Usage:     "Get all email restrictions",
				ArgsUsage: "[zoid]",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/policy/mailRestriction", nil)
				},
			},
			{
				Name:      "account-restrictions",
				Usage:     "Get all account restrictions",
				ArgsUsage: "[zoid]",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/policy/accountRestriction", nil)
				},
			},
			{
				Name:      "access-restrictions",
				Usage:     "Get all access restrictions",
				ArgsUsage: "[zoid]",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/policy/accessRestriction", nil)
				},
			},
			{
				Name:      "forward-restrictions",
				Usage:     "Get all forward restrictions",
				ArgsUsage: "[zoid]",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid := cmd.Args().First()
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					return request(c, "GET", "/organization/"+zoid+"/policy/mailForwardPolicy", nil)
				},
			},
			{
				Name:      "policy-users",
				Usage:     "Get policy users",
				ArgsUsage: "<zoid> <policyId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, policyID := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if policyID == "" {
						return internal.NewValidationError("policyId required")
					}
					return request(c, "GET", "/organization/"+zoid+"/policy/"+policyID+"/getUsers", nil)
				},
			},
			{
				Name:      "policy-groups",
				Usage:     "Get policy groups",
				ArgsUsage: "<zoid> <policyId>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, policyID := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if policyID == "" {
						return internal.NewValidationError("policyId required")
					}
					return request(c, "GET", "/organization/"+zoid+"/policy/"+policyID+"/getGroups", nil)
				},
			},
			{
				Name:      "apply",
				Usage:     "Apply policy to users/groups or assign restrictions (use JSON body)",
				ArgsUsage: "<zoid> <policyId>",
				Flags:     []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}},
				Action: func(_ context.Context, cmd *cli.Command) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					zoid, policyID := cmd.Args().Get(0), cmd.Args().Get(1)
					if zoid == "" {
						zoid, err = mustOrg(cmd)
						if err != nil {
							return err
						}
					}
					if policyID == "" {
						return internal.NewValidationError("policyId required")
					}
					var body any
					json.Unmarshal([]byte(cmd.String("json")), &body)
					return request(c, "PUT", "/organization/"+zoid+"/policy/"+policyID, &zohttp.RequestOpts{JSON: body})
				},
			},
		},
	}
}
