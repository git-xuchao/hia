package main

import (
	"fmt"
	"os"
	/*
	 *"sort"
	 */

	"gopkg.in/urfave/cli.v1"
)

var (
	addCommand = cli.Command{
		Name:    "hello",
		Aliases: []string{"a"},
		Usage:   "hello test command",
		Action: func(c *cli.Context) error {
			fmt.Println("run hello command: ", c.Args().First())
			return nil
		},
	}

	completeCommand = cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "complete a task on the list",
		Action: func(c *cli.Context) error {
			fmt.Println("completed task: ", c.Args().First())
			return nil
		},
	}

	testCommand = cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "options for task templates",
		Subcommands: []cli.Command{
			{
				Name:  "net",
				Usage: "test net",
				Action: func(c *cli.Context) error {
					fmt.Println("test net: ", c.Args().First())
					return nil
				},
			},
			{
				Name:  "redis",
				Usage: "test redis",
				Action: func(c *cli.Context) error {
					fmt.Println("test redis: ", c.Args().First())
					return nil
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		addCommand,
		completeCommand,
		testCommand,
	}

	app.Run(os.Args)
}
