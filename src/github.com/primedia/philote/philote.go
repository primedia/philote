package main

import (
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

const philoteConfig string = "/usr/local/etc/philote"
const philoteAnsiblefile string = philoteConfig + "/Ansiblefile"

func main() {
	app := cli.NewApp()
	app.Name = "philote"
	app.Version = "1.0.0"
	app.Usage = "setup your machine with the power of Ansible"
	app.Commands = []cli.Command{
		{
			Name:  "install",
			Usage: "install specified Ansible roles",
			Action: func(c *cli.Context) {
				println("installing...")
			},
		},
		{
			Name:  "roles",
			Usage: "display roles that will be installed",
			Action: func(c *cli.Context) {
				println("roles...")
			},
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add specified role",
					Action: func(c *cli.Context) {
						println("adding roles:", c.Args().First())
						addAnsibleRole(c.Args()[0], "foo")
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "github, g",
							Usage: "github owner/repository for the role, e.g. primedia/idg",
						},
						cli.StringFlag{
							Name:  "path, p",
							Usage: "absolute path for the role, e.g. /Users/johndoe/philote/myrole",
						},
					},
				},
				{
					Name:  "remove",
					Usage: "remove specified role(s)",
					Action: func(c *cli.Context) {
						println("removing roles:", strings.Join(c.Args(), ", "))
						for _, roleName := range c.Args() {
							removeAnsiblefileRole(roleName)
						}
					},
				},
			},
		},
		{
			Name:  "setup",
			Usage: "setup philote with roles and then install them",
			Action: func(c *cli.Context) {

			},
		},
	}
	app.Run(os.Args)
}
