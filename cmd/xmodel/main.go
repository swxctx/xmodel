package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/xiaoenai/xmodel/cmd/create"
	"github.com/xiaoenai/xmodel/cmd/info"
)

func main() {
	app := cli.NewApp()
	app.Name = "Micro project aids"
	app.Version = "6.1.1"
	app.Author = "xiaoenai"
	app.Usage = "a deployment tools of tp-micro frameware"

	// new a project
	newCom := cli.Command{
		Name:  "gen",
		Usage: "Generate a tp-micro project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
			cli.BoolFlag{
				Name:  "force, f",
				Usage: "Forced to rebuild the whole project",
			},
			cli.BoolFlag{
				Name:  "newdoc",
				Usage: "Rebuild the README.md",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.CreateProject(c.Bool("force"), c.Bool("newdoc"))
			return nil
		},
	}

	// new a README.md
	newdocCom := cli.Command{
		Name:  "newdoc",
		Usage: "Generate a tp-micro project README.md",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.CreateDoc()
			return nil
		},
	}

	// run a project
	runCom := cli.Command{
		Name:  "run",
		Usage: "Compile and run gracefully (monitor changes) an any existing go project",
		UsageText: `micro run [options] [arguments...]
 or
   micro run [options except -app_path] [arguments...] {app_path}`,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "watch_exts, x",
				Value: (*cli.StringSlice)(&[]string{".go", ".ini", ".yaml", ".toml", ".xml"}),
				Usage: "Specified to increase the listening file suffix",
			},
			cli.StringSliceFlag{
				Name:  "notwatch, n",
				Value: (*cli.StringSlice)(&[]string{}),
				Usage: "Not watch files or directories",
			},
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			run.RunProject(c.StringSlice("watch_exts"), c.StringSlice("notwatch"))
			return nil
		},
	}

	// add mysql model struct code to project template
	tplCom := cli.Command{
		Name:  "tpl",
		Usage: "Add mysql model struct code to project template",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
			cli.StringFlag{
				Name:  "host",
				Value: "localhost",
				Usage: "mysql host ip",
			},
			cli.StringFlag{
				Name:  "port",
				Value: "3306",
				Usage: "mysql host port",
			},
			cli.StringFlag{
				Name:  "username, user",
				Value: "root",
				Usage: "mysql username",
			},
			cli.StringFlag{
				Name:  "password, pwd",
				Value: "",
				Usage: "mysql password",
			},
			cli.StringFlag{
				Name:  "db",
				Value: "test",
				Usage: "mysql database",
			},
			cli.StringSliceFlag{
				Name:  "table",
				Usage: "mysql table",
			},
			cli.StringFlag{
				Name:  "ssh_user",
				Value: "",
				Usage: "ssh user",
			},
			cli.StringFlag{
				Name:  "ssh_host",
				Value: "",
				Usage: "ssh host ip",
			},
			cli.StringFlag{
				Name:  "ssh_port",
				Value: "",
				Usage: "ssh host port",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.AddTableStructToTpl(create.ConnConfig{
				MysqlConfig: create.MysqlConfig{
					Host:     c.String("host"),
					Port:     c.String("port"),
					User:     c.String("user"),
					Password: c.String("password"),
					Db:       c.String("db"),
				},
				Tables:  c.StringSlice("table"),
				SshHost: c.String("ssh_host"),
				SshPort: c.String("ssh_port"),
				SshUser: c.String("ssh_user"),
			})
			return nil
		},
	}

	app.Commands = []cli.Command{newCom, newdocCom, runCom, tplCom}
	app.Run(os.Args)
}

func initProject(c *cli.Context) error {
	appPath := c.String("app_path")
	if len(appPath) == 0 {
		appPath = c.Args().First()
	}
	if len(appPath) == 0 {
		appPath = "./"
	}
	return info.Init(appPath)
}
