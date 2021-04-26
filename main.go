package main

import (
	"os"

	"github.com/NETWAYS/go-check"
)

const readme = `Notifications via a brevis.one gateway.
Sends SMS or rings at a given number`

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	plugin.Name = "notify-brevisone"
	plugin.Readme = readme
	plugin.Timeout = 30
	plugin.Version = buildVersion()

	config := &Config{}
	config.BindArguments(plugin.FlagSet)

	plugin.ParseArguments()

	if len(os.Args) <= 1 {
		plugin.FlagSet.Usage()
		check.Exit(check.Unknown, "No arguments given")
	}

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	plugin.SetupLogging()

	err = config.Run()
	if err != nil {
		check.ExitError(err)
	}

	check.Exit(check.OK, "done")
}
