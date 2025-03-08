package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dcommisso/img/internal/mgparser"
)

const (
	maxNumberOfPanels = 4
)

type model struct {
	panels       [maxNumberOfPanels]panel
	focused      int
	windowWidth  int
	windowHeight int
}

func (m *model) setSize(width, height int) {
	m.windowWidth = width
	m.windowHeight = height
}

func (m *model) getActivePanels() []panel {
	var activePanels []panel
	for _, panel := range m.panels {
		if panel.active {
			activePanels = append(activePanels, panel)
		}
	}
	return activePanels
}

func (m *model) increaseFocused() {
	// Do nothing if we're on the last panel or if the next panel is empty
	if m.focused == len(m.panels)-1 || !m.panels[m.focused+1].active {
		return
	}

	m.focused++
}

func (m *model) decreaseFocused() {
	if m.focused == 0 {
		return
	}

	m.focused--
	panelToDelete := m.focused + 2
	if panelToDelete < len(m.panels) {
		m.deletePanel(panelToDelete)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		for i, p := range m.getActivePanels() {
			res, cmd := p.Update(msg)
			m.panels[i] = res.(panel)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "right":
			m.increaseFocused()
		case "left":
			m.decreaseFocused()
		}
	}

	// send the msg to the focused panel. This must be done before getting the
	// selected item.
	var res tea.Model
	res, cmd = m.panels[m.focused].Update(msg)
	m.panels[m.focused] = res.(panel)

	m.updateNextPanel()

	return m, cmd
}

func (m model) View() string {
	activePanels := []string{}
	// TODO: use GetActivePanels funcion
	for _, panel := range m.panels {
		if panel.active {
			activePanels = append(activePanels, panel.View())
		} else {
			break
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, activePanels...)
}

func main() {
	mgPathToLoad := "./internal/mgparser/testdata/mgs/validMg"
	mgToLoad, err := mgparser.NewMg(mgPathToLoad)
	if err != nil {
		fmt.Println("Error loading mg:", err)
		os.Exit(1)
	}

	ocpResources := []ActionableElement{
		// a pointer is needed since aeLogs Init method has a pointer receiver
		new(aeLogs),
	}

	// initialize the ActionableElement and add them to firstPanelItems
	for _, elem := range ocpResources {
		elem.Init(mgToLoad)
	}
	m := model{}
	m.addNewPanel(0, ocpResources)
	m.panels[0].active = true

	m.updateNextPanel()
	//	m := model{
	//		focused: 0,
	//		panels: [maxNumberOfPanels]panel{
	//			panel{
	//				Model:  list.New(aeSliceToItem(ocpResources), list.NewDefaultDelegate(), 50, 10),
	//				active: true,
	//			},
	//		},
	//		windowWidth:  500,
	//		windowHeight: 100,
	//	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
