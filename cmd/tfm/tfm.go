package main

import (
	"fmt"
	"os"

	"github.com/malnick/tfm/component"
	logging "github.com/op/go-logging"
	"github.com/urfave/cli"
)

var (
	VERSION  = "UNSET"
	REVISION = "UNSET"
	log      = logging.MustGetLogger("tfm.main")
	format   = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{longfunc} ▶ %{level:.6s} %{id:03x}%{color:reset} %{message}`,
	)
)

type Config struct {
	Debug         bool
	ComponentName string
	Workspace     string
	TerraformPath string
	ListDetails   bool
	FlagDependsOn string
}

func main() {
	logging.SetFormatter(format)

	tfm := &Config{
		TerraformPath: "terraform/",
		Workspace:     component.WorkspaceDevelopment,
	}

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
			Name:        "debug",
			Usage:       "Run in debug mode",
			Destination: &tfm.Debug,
		},
		cli.StringFlag{
			Name:        "workspace, w",
			Usage:       "Choose the workspace to work in",
			Value:       tfm.Workspace,
			Destination: &tfm.Workspace,
		},
		cli.StringFlag{
			Name:        "terraform-path, tp",
			Usage:       "Choose the top level path where your Terraform directory structure lives",
			Value:       tfm.TerraformPath,
			Destination: &tfm.TerraformPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "component",
			Usage: "Component master command",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:        "long,l",
							Usage:       "Include component metadata details in listing",
							Destination: &tfm.ListDetails,
						},
					},
					Usage: "List terraform components",
					Action: func(cliContext *cli.Context) error {
						name := cliContext.Args().Get(0)
						c, err := makeComponent(name, tfm)
						if err != nil {
							log.Error(err)
							return err
						}
						return c.List(tfm.ListDetails)
					},
				},
				{
					Name:  "create",
					Usage: "Make a new component for a deployment with defaults",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "depends",
							Usage:       "Comma separated list of dependant component names",
							Value:       tfm.FlagDependsOn,
							Destination: &tfm.FlagDependsOn,
						},
					},
					Action: func(cliContext *cli.Context) error {
						name := cliContext.Args().Get(0)
						c, err := makeComponent(name, tfm)
						if err != nil {
							log.Error(err)
							return err
						}

						return component.Create(c)
					},
				},
				{
					Name:  "plan",
					Usage: "Run a terraform plan for a component",
					Action: func(cliContext *cli.Context) error {
						name := cliContext.Args().Get(0)
						c, err := makeComponent(name, tfm)
						if err != nil {
							return err
						}

						return component.Run(component.ActionPlan, c)
					},
				},
				{
					Name:  "apply",
					Usage: "Run a terraform apply for a component",
					Action: func(cliContext *cli.Context) error {
						name := cliContext.Args().Get(0)
						c, err := makeComponent(name, tfm)
						if err != nil {
							return err
						}

						return component.Run(component.ActionApply, c)
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func makeComponent(name string, tfm *Config) (*component.Component, error) {
	return component.New(
		component.OptionDependsOn(tfm.TerraformPath, tfm.FlagDependsOn),
		component.OptionName(name),
		component.OptionWorkingDirectory(tfm.TerraformPath),
		component.OptionWorkspace(tfm.Workspace),
		component.OptionAllowedWorkspaces([]string{tfm.Workspace}))
}
