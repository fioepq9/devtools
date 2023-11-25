package main

import (
	"os"

	"github.com/fioepq9/devtools/cmd/modelgen"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		modelgen.Command(),
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
