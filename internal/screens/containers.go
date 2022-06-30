package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wesleimp/docker-clean/internal/components/help"
	"github.com/wesleimp/docker-clean/internal/styles"
	"github.com/wesleimp/docker-clean/pkg/containers"
)

const (
	iconSelected    = "●"
	iconNotSelected = "○"
)

type ContainersModel struct {
	err        error
	client     *client.Client
	containers []types.Container
	selected   map[int]struct{}
	cursor     int
	all        bool
}

func NewContainersModel(client *client.Client, containers []types.Container) ContainersModel {
	return ContainersModel{
		client:     client,
		containers: containers,
		selected:   map[int]struct{}{},
	}
}
func (m ContainersModel) Init() tea.Cmd {
	return nil
}

func (m ContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.containers)-1 {
				m.cursor++
			}
		case "s":
			for i := range m.containers {
				m.selected[i] = struct{}{}
			}
		case "n":
			for i := range m.selected {
				delete(m.selected, i)
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "a":
			m.selected = map[int]struct{}{} // reset selection
			m.cursor = 0                    // reset cursor
			m.all = !m.all
			contaiers, err := containers.List(m.client, m.all)
			if err != nil {
				m.err = err
			}
			m.containers = contaiers
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m ContainersModel) View() string {
	if m.err != nil {
		errs := styles.Error.Copy()
		return errs.MarginLeft(2).Render("Failed to load containers.")
	}

	marker := lipgloss.NewStyle().Width(styles.MARKER_COLUMN_SIZE)
	id := lipgloss.NewStyle().Width(styles.ID_COLUMN_SIZE)
	name := lipgloss.NewStyle().Width(styles.NAME_COLUMN_SIZE)
	imageName := lipgloss.NewStyle().Width(styles.IMAGE_NAME_COLUMN_SIZE)
	status := lipgloss.NewStyle().Width(styles.STATUS_COLUMN_SIZE)
	state := lipgloss.NewStyle().Width(styles.STATE_COLUMN_SIZE)

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		marker.Render(""),
		id.Bold(true).Render("ID"),
		name.Bold(true).Render("NAME"),
		imageName.Bold(true).Render("IMAGE"),
		status.Bold(true).Render("STATUS"),
		state.Bold(true).Render("STATE"),
	)

	var rows []string
	for i, container := range m.containers {
		style := styles.Row
		if i == m.cursor {
			style = styles.SelectedRow
		}

		icon := iconNotSelected
		if _, ok := m.selected[i]; ok {
			icon = iconSelected
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			marker.Render(icon),
			id.Bold(false).Render(container.ID[:16]),
			name.Bold(false).Render(strings.Join(container.Names, ",")),
			imageName.Bold(false).Render(container.Image),
			status.Bold(false).Render(container.Status),
			state.Bold(false).Render(container.State),
		)

		rows = append(rows, style.Render(row))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Header.Render(header),
		lipgloss.JoinVertical(lipgloss.Top, rows...),
	) + helpView()
}

func helpView() string {
	return help.View([]help.HelpOption{
		{Key: "up/down j/k", Help: "navigate"},
		{Key: "space/enter", Help: "toggle selection"},
		{Key: "a", Help: "toggle list all"},
		{Key: "d", Help: "delete all selected"},
		{Key: "s", Help: "select all"},
		{Key: "n", Help: "deselect all"},
		{Key: "q/esc", Help: "quit"},
	})
}
