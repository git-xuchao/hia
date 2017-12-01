package main

import (
	"fmt"
	"hia/http/httpc"
	"hia/http/httpd"
	"hia/redis"
	"os"
	/*
	 *"sort"
	 */

	"gopkg.in/urfave/cli.v1"
)

func hello_command(c *cli.Context) error {
	fmt.Println("run hello command: ", c.Args().First())
	return nil
}

func test_net_command(c *cli.Context) error {
	fmt.Println("test net: ", c.Args().First())
	return nil
}

func test_redis_command(c *cli.Context) error {
	fmt.Println("test redis: ", c.Args().First())
	redis.TestRedis()
	return nil
}

func test_httpc_command(c *cli.Context) error {
	httpc.TestHttpcCommand()
	return nil
}

func test_httpd_command(c *cli.Context) error {
	httpd.TestHttpdCommand()
	return nil
}

var (
	app = cli.NewApp()

	addCommand = cli.Command{
		Name:    "hello",
		Aliases: []string{"a"},
		Usage:   "hello test command",
		Action:  hello_command,
	}

	testCommand = cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "options for task templates",
		Subcommands: []cli.Command{
			{
				Name:   "net",
				Usage:  "test net",
				Action: test_net_command,
			},
			{
				Name:   "redis",
				Usage:  "test redis",
				Action: test_redis_command,
			},
			{
				Name:   "httpc",
				Usage:  "test httpc",
				Action: test_httpc_command,
			},
			{
				Name:   "httpd",
				Usage:  "test httpd",
				Action: test_httpd_command,
			},
		},
	}
)

func init() {
	app.Commands = []cli.Command{
		addCommand,
		testCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
