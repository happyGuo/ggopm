package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"ggopm/action"
	"ggopm/msg"
	"ggopm/repo"
)

var version = "0.1.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "ggopm"
	app.Version = version

	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q", c.Args().Get(0))
		return nil
	}

	//cli 执行的命令
	app.Commands = commands()

	app.Run(os.Args)
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "init",
			Usage: "初始化一个新项目，创建依赖配置文件",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "skip-scan",
					Usage: "不扫描项目源代码，只是生成一个依赖配置文件",
				},
			},

			Action: func(c *cli.Context) error {
				action.Create(".", c.Bool("skip-scan"))
				return nil
			},
		},
		{
			Name:  "install, c",
			Usage: "安装依赖包到指定目录，不指定则在GOAPTH的第一个目录",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "path",
					Usage: "安装到指定目录",
				},
				cli.BoolFlag{
					Name:  "skip-p",
					Usage: "跳过扫描项目依赖",
				},
			},

			Action: func(c *cli.Context) error {
				if c.String("path") != "" {
					msg.Warn("依赖包将会安装到=>" + c.String("path"))
				}
				installer := repo.NewInstaller()
				installer.Force = c.Bool("force")
				installer.Home = c.GlobalString("home")

				action.Install(installer, c.Bool("strip-vendor"))

				return nil
			},
		},

		{
			Name:  "update, up",
			Usage: "更新依赖包",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Usage: "更新指定目录的包",
				},
			},

			Action: func(c *cli.Context) error {
				if c.String("path") != "" {
					msg.Warn("将会更新下面目录的包=>" + c.String("path"))
				}
				installer := repo.NewInstaller()
				installer.Force = c.Bool("force")
				installer.Home = c.GlobalString("home")

				action.Install(installer, c.Bool("strip-vendor"))

				return nil
			},
		},

		{
			Name:  "get",
			Usage: "单独安装一个包",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Usage: "安装到指定目录",
				},
			},

			Action: func(c *cli.Context) error {
				if c.String("path") != "" {
					msg.Warn("包将会安装到=>" + c.String("path"))
				}
				installer := repo.NewInstaller()
				installer.Force = c.Bool("force")
				installer.Home = c.GlobalString("home")

				action.Install(installer, c.Bool("strip-vendor"))

				return nil
			},
		},
	}
}
