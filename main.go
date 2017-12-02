package main

import (
	"fmt"
	"hia/http"
	"hia/lab"
	"hia/redis"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func helloCommand(c *cli.Context) error {
	fmt.Println("run hello command: ", c.Args().First())
	return nil
}

func lab1Command(c *cli.Context) error {
	fmt.Println("lab1Command: ", c.Args().First())
	lab.Lab1Command()
	return nil
}

func lab2Command(c *cli.Context) error {
	fmt.Println("lab1Command: ", c.Args().First())
	lab.Lab2Command()
	return nil
}

func lab3Command(c *cli.Context) error {
	fmt.Println("lab1Command: ", c.Args().First())
	lab.Lab3Command()
	return nil
}

func testNetCommand(c *cli.Context) error {
	fmt.Println("test net: ", c.Args().First())
	return nil
}

func testRedisCommand(c *cli.Context) error {
	fmt.Println("test redis: ", c.Args().First())
	redis.TestRedis()
	return nil
}

func testHttpcCommand(c *cli.Context) error {
	http.TestHttpcCommand()
	return nil
}

func testHttpdCommand(c *cli.Context) error {
	http.TestHttpdCommand()
	return nil
}

var (
	app = cli.NewApp()

	addCommand = cli.Command{
		Name:    "hello",
		Aliases: []string{"a"},
		Usage:   "hello test command",
		Action:  helloCommand,
	}

	labCommand = cli.Command{
		Name:    "lab",
		Aliases: []string{"l"},
		Usage:   "options for task templates",
		Subcommands: []cli.Command{
			{
				Name:   "lab1",
				Usage:  "lab net",
				Action: lab1Command,
			},
			{
				Name:   "lab2",
				Usage:  "lab net",
				Action: lab2Command,
			},
			{
				Name:   "lab3",
				Usage:  "lab net",
				Action: lab3Command,
			},
		},
	}

	testCommand = cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "options for task templates",
		Subcommands: []cli.Command{
			{
				Name:   "net",
				Usage:  "test net",
				Action: testNetCommand,
			},
			{
				Name:   "redis",
				Usage:  "test redis",
				Action: testRedisCommand,
			},
			{
				Name:   "httpc",
				Usage:  "test httpc",
				Action: testHttpcCommand,
			},
			{
				Name:   "httpd",
				Usage:  "test httpd",
				Action: testHttpdCommand,
			},
		},
	}
)

func init() {
	app.Commands = []cli.Command{
		addCommand,
		labCommand,
		testCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
