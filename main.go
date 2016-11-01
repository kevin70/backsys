package main

import (
	"github.com/kevin70/backsys/cmd"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := &cli.App{}
	cmd.AddWeixinCmd(app)

	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Kevin Zou",
			Email: "kevinz@weghst.com",
		},
	}
	app.Name = "backsys"
	app.Usage = "系统工具包"
	app.Version = "1.0.0"
	app.Run(os.Args)
}
