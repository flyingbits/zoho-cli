package books

import (
	"github.com/urfave/cli/v3"
)

func openingBalanceCmd() *cli.Command {
	return &cli.Command{
		Name:  "openingbalance",
		Usage: "Opening balance operations",
		Commands: []*cli.Command{
			{Name: "get", Usage: "Get opening balance", Action: runGetNoID("openingbalance")},
			{Name: "create", Usage: "Create opening balance", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runCreate("openingbalance")},
			{Name: "update", Usage: "Update opening balance", Flags: []cli.Flag{&cli.StringFlag{Name: "json", Required: true, Usage: "JSON body"}}, Action: runUpdateNoID("openingbalance")},
			{Name: "delete", Usage: "Delete opening balance", Action: runDeleteNoID("openingbalance")},
		},
	}
}
