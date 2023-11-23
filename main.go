package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed files/*
var stackFile embed.FS

const (
	SCREEN_MAIN  int = 1
	SCREEN_CHECK int = 2
	SCREEN_STACK int = 3
)

const (
	COMMAND_FAILED string = "failed"
	COMMAND_OK     string = "ok"
)

type checks struct {
	dockerInstalled   bool
	dockerSwarmActive bool
}

type model struct {
	table table.Model

	stackTable    table.Model
	acceptedSetup bool
	cursor        int
	choices       []string
	selected      map[int]struct{}

	activeScreen int

	prompting bool
	checking  bool

	finishedChecks bool
	currentMessage string

	err error

	textInput textinput.Model

	checks checks
}

var cleanStyle = lipgloss.NewStyle()
var coloredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Bold(true)

var instructionCols = []table.Column{
	{Title: "Command", Width: 20},
	{Title: "Description", Width: 40},
}

var instructionRows = []table.Row{
	{"which docker", "Check if docker is installed"},
	{"docker info", "Check if docker swarm is active"},
}

var servicesCols = []table.Column{
	{Title: "Service", Width: 20},
	{Title: "Description", Width: 40},
}

var servicesRows = []table.Row{
	{"traefik", "Proxy, auto service discovery and more"},
	{"portainer", "Agent/UI to manage and deploy services"},
	{"vector", "Agent to capture all container logs"},
	{"loki", "Log index to visualize logs"},
	{"grafana", "All the graphs you will ever need"},
}

// Messages
type msgSetupAccepted struct{}

func initialModel(t table.Model, st table.Model) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()

	return model{
		cursor:       0,
		choices:      []string{"A", "B", "C"},
		selected:     make(map[int]struct{}),
		table:        t,
		prompting:    false,
		activeScreen: SCREEN_MAIN,
		textInput:    ti,
		checks:       checks{},
		stackTable:   st,
	}
}

func (m model) Init() tea.Cmd {
	return m.textInput.Cursor.BlinkCmd()
}

func (m *model) setPrompting(p bool) *model {
	m.prompting = true
	return m
}

func (m model) View() string {
	s := ""
	status := "⚪️"
	statusSwarm := "⚪️"

	if m.activeScreen == SCREEN_MAIN {
		s += fmt.Sprintf("%s\n\n%s\n",
			cleanStyle.Bold(true).Padding(1, 0).Border(lipgloss.RoundedBorder()).Padding(1, 2).
				Render(fmt.Sprintf("%s\n%s",
					cleanStyle.Foreground(lipgloss.Color("45")).Render("One Line Command Tool for Self Hosting"),
					cleanStyle.Render("v0.1.0"))),
			cleanStyle.
				Render("The script needs to check that you have all\nrequirements in place before the setup.\n\n"))
		s += m.table.View()
		s += "\n\n\n"

		if !m.acceptedSetup {
			s += cleanStyle.Bold(true).Foreground(lipgloss.Color("220")).Render("Proceed? (Y/N)")
			s += "\n"

			if m.prompting {
				s += m.textInput.View()
			}

		}
	}

	if m.activeScreen == SCREEN_CHECK {
		s += cleanStyle.Foreground(lipgloss.Color("155")).Border(lipgloss.NormalBorder(), false, false, true, false).Render("Running checks...\n" + "\t\t\t\t\t\t")
		s += "\n"

		// if m.checking && !m.checks.dockerInstalled {
		// 	status = "❌"
		// }
		if m.checking && m.checks.dockerInstalled {
			status = "✅"
		}

		// if m.checking && !m.checks.dockerSwarmActive {
		// 	statusSwarm = "❌"
		// }

		if m.checking && m.checks.dockerSwarmActive {
			statusSwarm = "✅"
		}

		s += fmt.Sprintf("[%s] Docker is installed\n[%s] Docker Swarm is active", status, statusSwarm)
		s += "\n"

		if m.err != nil {
			s += cleanStyle.Foreground(lipgloss.Color("210")).Render(fmt.Sprintf("\n\nError: %v", m.err))
		}

		if m.checking && m.checks.dockerInstalled && m.checks.dockerSwarmActive {
			s += cleanStyle.Bold(true).Foreground(lipgloss.Color("35")).Render("\n\n\n\nEverything looks good. \nNext we will install all the required services into your docker swarm.")
			s += "\n\n"
			s += m.stackTable.View()
			s += "\n\n"
			s += cleanStyle.Bold(false).Foreground(lipgloss.Color("220")).Render("\nProceed with next steps? (Y/N)")
			s += "\n"
			s += m.textInput.View()
		}
	}

	if m.activeScreen == SCREEN_STACK && m.finishedChecks {
		s += "Stack here"
	}

	// Iterate over our choices
	// for i, choice := range m.choices {

	// 	// Is the cursor pointing at this choice?
	// 	cursor := " " // no cursor
	// 	if m.cursor == i {
	// 		cursor = ">" // cursor!
	// 	}

	// 	// Is this choice selected?
	// 	checked := " " // not selected
	// 	if _, ok := m.selected[i]; ok {
	// 		checked = "x" // selected!
	// 	}

	// 	// Render the row
	// 	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	// }

	// The footer
	s += cleanStyle.Foreground(lipgloss.Color("240")).Render("\n\n\nPress q to quit.\n")

	// Send the UI for rendering
	return s
}

type cmdDockerMsg struct{ err error }
type cmdDockerSwarmInfoMsg struct{ err error }

func cmdDockerExec() tea.Cmd {
	r, s := exec.Command("/usr/bin/which", "docker"), new(strings.Builder)
	r.Stdout = s
	r.Stderr = s
	r.Run()

	return func() tea.Msg {
		if !strings.Contains(s.String(), "docker") {
			return cmdDockerMsg{fmt.Errorf("could not find docker executable. Please install docker. %s", s.String())}
		}
		return cmdDockerMsg{}
	}
}

func cmdDockerSwarmInfo() tea.Cmd {

	dc, b := exec.Command("/usr/bin/which", "docker"), new(strings.Builder)
	dc.Stdout = b
	dc.Run()

	r, sb := exec.Command("docker", "info", "--format", "'{{.Swarm.LocalNodeState}}'"), new(strings.Builder)
	r.Stdout = sb
	r.Stderr = sb
	r.Run()

	return func() tea.Msg {
		if strings.Contains(sb.String(), "active") {
			return cmdDockerSwarmInfoMsg{}
		}
		return cmdDockerSwarmInfoMsg{fmt.Errorf("docker swarm is not active. \nPlease initialize one with `docker swarm init`\n %s", sb.String())}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case cmdDockerMsg:
		if msg.err != nil {
			m.checks.dockerInstalled = false
			m.err = msg.err
		} else {
			m.checks.dockerInstalled = true
			return m, cmdDockerSwarmInfo()
		}

		return m, nil

	case cmdDockerSwarmInfoMsg:
		if msg.err != nil {
			m.checks.dockerSwarmActive = false
			m.err = msg.err
		} else {
			m.checks.dockerSwarmActive = true
		}

		return m, nil

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.activeScreen == SCREEN_MAIN && m.prompting && strings.ToLower(m.textInput.Value()) == "y" {
				m.acceptedSetup = true
				m.activeScreen = SCREEN_CHECK
				m.prompting = false
				m.textInput.Reset()

			} else if m.activeScreen == SCREEN_MAIN && m.prompting && strings.ToLower(m.textInput.Value()) == "n" {
				return m, tea.Quit
			}

			if m.activeScreen == SCREEN_CHECK && strings.ToLower(m.textInput.Value()) == "y" {
				m.finishedChecks = true
				m.activeScreen = SCREEN_STACK
				m.textInput.Reset()

			} else if m.activeScreen == SCREEN_CHECK && strings.ToLower(m.textInput.Value()) == "n" {
				return m, tea.Quit
			}

			return m, nil
		case "N", "n":
			if m.activeScreen == SCREEN_MAIN {
				return m, tea.Quit
			}
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		}
	}

	if m.activeScreen == SCREEN_CHECK && !m.checking {
		// Run checks
		m.checking = true
		return m, cmdDockerExec()
	}

	m.setPrompting(true)
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd

}

func main() {
	t := table.New(
		table.WithColumns(instructionCols),
		table.WithRows(instructionRows),
		table.WithFocused(false),
		table.WithHeight(2),
	)

	stackTable := table.New(
		table.WithColumns(servicesCols),
		table.WithRows(servicesRows),
		table.WithFocused(false),
		table.WithHeight(6),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("250")).Italic(true)
	s.Selected.UnsetForeground().UnsetBold()
	s.Cell.Foreground(lipgloss.Color("250")).Italic(true)
	t.SetStyles(s)
	stackTable.SetStyles(s)

	p := tea.NewProgram(initialModel(t, stackTable))
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
