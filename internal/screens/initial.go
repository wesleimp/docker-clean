package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/wesleimp/docker-clean/internal/styles"
	"github.com/wesleimp/docker-clean/pkg/containers"
)

// InitialModel is the UI when the CLI starts, basically loading the repos.
type InitialModel struct {
	client  *client.Client
	err     error
	spinner spinner.Model
	loading bool
}

// NewInitialModel creates a new InitialModel with required fields.
func NewInitialModel(c *client.Client) InitialModel {
	s := spinner.NewModel()
	s.Spinner = spinner.MiniDot

	return InitialModel{
		client:  c,
		spinner: s,
		loading: true,
	}
}

func (m InitialModel) Init() tea.Cmd {
	return tea.Batch(fetchContainers(m.client), spinner.Tick)
}

func (m InitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case containerListMsg:
		cm := NewContainersModel(m.client, msg.containers)
		return cm, cm.Init()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m InitialModel) View() string {
	if m.loading {
		return styles.Primary.Render(m.spinner.View()) + " Gathering containers..."
	}

	if m.err != nil {
		return styles.Error.Render(fmt.Sprintf("Error gathering the containers list: %s", m.err))
	}

	return ""
}

type errorMsg struct {
	err error
}

type containerListMsg struct {
	containers []types.Container
}

func fetchContainers(client *client.Client) tea.Cmd {
	return func() tea.Msg {
		cc, err := containers.List(client, false)
		if err != nil {
			return errorMsg{err: err}
		}

		return containerListMsg{containers: cc}
	}
}
