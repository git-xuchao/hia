package main

import (
	"fmt"
	"os"
	/*
	 *"sort"
	 */

	"gopkg.in/urfave/cli.v1"
)

/*
 *func main() {
 *    app := cli.NewApp()
 *    app.Name = "greet"
 *    app.Usage = "fight the loneliness!"
 *    app.Action = func(c *cli.Context) error {
 *        fmt.Println("Hello friend!")
 *        return nil
 *    }
 *
 *    app.Run(os.Args)
 *}
 */
/*
 *func main() {
 *    app := cli.NewApp()
 *
 *    app.Flags = []cli.Flag{
 *        cli.StringFlag{
 *            Name:  "lang",
 *            Value: "english",
 *            Usage: "language for the greeting",
 *        },
 *    }
 *
 *    app.Action = func(c *cli.Context) error {
 *        name := "Nefertiti"
 *        if c.NArg() > 0 {
 *            name = c.Args().Get(0)
 *        }
 *        if c.String("lang") == "spanish" {
 *            fmt.Println("Hola", name)
 *        } else {
 *            fmt.Println("Hello", name)
 *        }
 *        return nil
 *    }
 *
 *    app.Run(os.Args)
 *}
 */

/*
 *func main() {
 *    app := cli.NewApp()
 *
 *    app.Flags = []cli.Flag{
 *        cli.StringFlag{
 *            Name:  "lang, l",
 *            Value: "english",
 *            Usage: "Language for the greeting",
 *        },
 *        cli.StringFlag{
 *            Name:  "config, c",
 *            Usage: "Load configuration from `FILE`",
 *        },
 *    }
 *
 *    app.Commands = []cli.Command{
 *        {
 *            Name:    "complete",
 *            Aliases: []string{"c"},
 *            Usage:   "complete a task on the list",
 *            Action: func(c *cli.Context) error {
 *                fmt.Println("run complete")
 *                return nil
 *            },
 *        },
 *        {
 *            Name:    "add",
 *            Aliases: []string{"a"},
 *            Usage:   "add a task to the list",
 *            Action: func(c *cli.Context) error {
 *                fmt.Println("run add")
 *                return nil
 *            },
 *        },
 *    }
 *
 *    sort.Sort(cli.FlagsByName(app.Flags))
 *    sort.Sort(cli.CommandsByName(app.Commands))
 *
 *    app.Run(os.Args)
 *}
 */
func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
