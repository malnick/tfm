package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var (
	VERSION  = "UNSET"
	REVISION = "UNSET"
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s @ %s", VERSION, REVISION)
	app.Name = "tfm"
	app.Usage = "Terraform wrapper for implementing deployment and component abstraction"
	app.Authors = []cli.Author{
		{
			Name:  "Jeff Malnick",
			Email: "malnick at google mail",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Run in debug mode",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "Make a new component for a deployment",
		},
		{
			Name:  "plan",
			Usage: "Run a terraform plan for a component",
		},
		{
			Name:  "apply",
			Usage: "Run a terraform apply for a component",
		},
		{
			Name:  "list",
			Usage: "List all components",
		},
	}

	app.Run(os.Args)
}
