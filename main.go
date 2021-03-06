package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/client"
	"github.com/urfave/cli/v2"
	"github.com/wesleimp/docker-clean/internal/screens"
	"github.com/wesleimp/docker-clean/internal/styles"
)

type model struct {
	containers []string
	cursor     int
	selected   []string
}

func main() {
	app := &cli.App{
		Name:  "docker-clean",
		Usage: "Clean up your docker state interactively",
		Commands: []*cli.Command{
			{
				Name:    "containers",
				Aliases: []string{"c"},
				Usage:   "Interact with containers",
				Action: func(_ *cli.Context) error {
					return run(screens.ContainersScreen)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(styles.Error.Render(err.Error()))
		os.Exit(1)
	}
}

func run(s screens.Screen) error {
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	p := tea.NewProgram(screens.NewInitialModel(client, s))
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		return err
	}

	return nil
}
