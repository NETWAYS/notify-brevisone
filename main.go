package main

import (
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
	plugin.Version = "0.1"

	config := &Config{}
	config.BindArguments(plugin.FlagSet)

	plugin.ParseArguments()

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	err = config.Run()
	if err != nil {
		check.ExitError(err)
	}

	check.Exit(check.OK, "done")
}
