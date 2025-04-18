package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dcommisso/mugshine/internal/mgparser"
)

const (
	maxNumberOfPanels = 4
)

type mgBoard struct {
	clusterInfoPanel ClusterInfoPanel
	panels           [maxNumberOfPanels]panel
	focused          int
	windowWidth      int
	windowHeight     int
	keys             keyMap
	help             help.Model
}

func NewMgBoard(mustGatherPath string) (mgBoard, error) {
	// load must-gather
	mg, err := mgparser.NewMg(mustGatherPath)
	if err != nil {
		return mgBoard{}, err
	}

	var (
		ocpSupportedResources = []ActionableElement{
			// a pointer is needed since aeLogs Init method has a pointer receiver
			new(aeLogs),
			new(aeNodesResource),
			new(aeCoResource),
		}
	)

	newMgBoard := mgBoard{
		clusterInfoPanel: NewClusterInfoPanel(mg),
		keys:             keys,
		help:             help.New(),
	}

	resourcesToShow := []ActionableElement{}
	// initialize the ActionableElement with mg and add them to firstPanelItems
	// if they return at least one elements (e.g. inspect doesn't return nodes)
	for _, elem := range ocpSupportedResources {
		elem.Init(mg)
		if len(elem.Selected()) > 0 {
			resourcesToShow = append(resourcesToShow, elem)
		}
	}
	newMgBoard.AddNewPanel(0, resourcesToShow)
	newMgBoard.UpdatePanelsAfterMoving()

	return newMgBoard, nil
}

func (m *mgBoard) AddNewPanel(index int, actElements []ActionableElement) {
	items := aeSliceToItem(actElements)
	model := list.New(items, NewMgDelegate(), 0, 0)

	// set arabic type on pagination
	model.Paginator.Type = paginator.Arabic

	// disable help
	model.SetShowHelp(false)

	// remove esc from quit keys, since I want to use it only for clearing filter
	model.KeyMap.Quit = key.NewBinding(key.WithKeys("q"))

	// remap left/right to pagedown/pageup to avoid overlapping to navigation
	model.KeyMap.NextPage.SetKeys("pgdown")
	model.KeyMap.PrevPage.SetKeys("pgup")

	// header section
	model.Styles.Title = getListTitleStyle()
	if header := actElements[0].Header(); header != "" {
		model.Title = header
	}

	m.panels[index] = panel{
		list:   model,
		active: true,
		keys:   keys,
	}
}

// UpdatePanelsAfterMoving creates/updates/deletes the panels based on the item
// selected in focused panel
func (m *mgBoard) UpdatePanelsAfterMoving() {
	// update the status of the panels
	m.panels[m.focused].SetStatus(PanelStatusFocused)
	if m.focused-1 >= 0 {
		m.panels[m.focused-1].SetStatus(PanelStatusPrevious)
	}

	// do nothing else if we are at last panel
	lastPanelIndex := len(m.panels) - 1
	if m.focused == lastPanelIndex {
		return
	}

	// delete the next-next panel
	nextNextPanelIndex := m.focused + 2
	if nextNextPanelIndex < len(m.panels) {
		m.deletePanel(nextNextPanelIndex)
	}

	selectedItem := m.panels[m.focused].list.SelectedItem()

	// return if selectedItem is nil, otherwise it crashes when filtered with
	// unmatching strings
	if selectedItem == nil {
		return
	}

	itemsInNextPanel := selectedItem.(ActionableElement).Selected()

	// Add a new panel only if there are elements to show, otherwise delete it.
	if len(itemsInNextPanel) > 0 {
		m.AddNewPanel(m.focused+1, itemsInNextPanel)
		m.panels[m.focused+1].SetStatus(PanelStatusNext)
	} else {
		m.deletePanel(m.focused + 1)
	}
}

func (m *mgBoard) deletePanel(index int) {
	m.panels[index].active = false
}

func (m *mgBoard) setSize(width, height int) {
	m.windowWidth = width
	m.windowHeight = height
}

func (m *mgBoard) GetActivePanels() []*panel {
	var activePanels []*panel
	for i := 0; i < len(m.panels); i++ {
		if m.panels[i].active {
			activePanels = append(activePanels, &m.panels[i])
		}
	}
	return activePanels
}

func (m *mgBoard) IncreaseFocused() {
	// Do nothing if we're on the last panel or if the next panel is empty
	if m.focused == len(m.panels)-1 || !m.panels[m.focused+1].active {
		return
	}
	m.focused++
}

func (m *mgBoard) DecreaseFocused() {
	if m.focused == 0 {
		return
	}

	m.focused--
}

func (m mgBoard) Init() tea.Cmd {
	return nil
}

func (m *mgBoard) updatePanelsSize() {
	m.updatePanelsWidth()
	m.updatePanelsHeight()
}

// updatePanelsWidth resizes each panel based on its status. The
// algorithm prioritizes the panels in this order:
//
// 1. focused panel
// 2. next panel
// 3. focused -1 panel
// 4. focused -2 panel
//
// Here's how the algorithm works: if the space wanted by a panel is
// available, then it gets it. If not, it will receive a `resizeDivisor`
// fraction of the available space. In this way the granted space will be
// less and less as the priority of the panel goes down. Since a panel has a
// minimum space occupied regardless the allowed width, there is the risk
// that the total sum of the panels will exceeds the available width, at the
// expense of the last and prioritized panels, that's why a
// `unresevedWidthPercentage` percentage is not allocated.
func (m *mgBoard) updatePanelsWidth() {
	const (
		resizeDivisor            = 2
		unresevedWidthPercentage = 0.08
	)
	availableWidth := m.windowWidth - int(float64(m.windowWidth)*unresevedWidthPercentage)

	var panelsInPrioOrder []*panel

	// the first one is the focused
	panelsInPrioOrder = append(panelsInPrioOrder, &m.panels[m.focused])

	// add next panels, if it exists
	if m.focused+1 < len(m.panels) && m.panels[m.focused+1].active {
		panelsInPrioOrder = append(panelsInPrioOrder, &m.panels[m.focused+1])
	}

	// add the previous panels in reverse order
	for i := m.focused - 1; i >= 0; i-- {
		panelsInPrioOrder = append(panelsInPrioOrder, &m.panels[i])
	}

	// give each panel the wanted width or a fraction of the available width
	for _, panel := range panelsInPrioOrder {
		wanted := panel.GetWantedWidth()
		var given int

		if wanted < availableWidth {
			given = wanted
		} else {
			given = availableWidth / resizeDivisor
		}

		panel.setWidth(given)
		availableWidth -= given
	}
}

// updatePanelsHeight calculates and sets the height of all panels based on the
// height of clusterInfo panel and mini help
func (m *mgBoard) updatePanelsHeight() {
	clusterInfo := m.clusterInfoPanel.Render(m.windowWidth)
	help := m.help.View(m.keys)

	clusterInfoHeight := lipgloss.Height(clusterInfo)
	helpHeight := lipgloss.Height(help)
	panelsAvailableHeight := m.windowHeight - clusterInfoHeight - helpHeight

	for _, panel := range m.GetActivePanels() {
		panel.setHeight(panelsAvailableHeight - panel.getStyle().GetVerticalFrameSize())
	}
}

func (m mgBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		m.updatePanelsSize()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Right):
			m.IncreaseFocused()
		case key.Matches(msg, m.keys.Left):
			m.DecreaseFocused()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	m.panels[m.focused], cmd = m.panels[m.focused].Update(msg)
	m.UpdatePanelsAfterMoving()
	m.updatePanelsSize()
	return m, cmd
}

func (m mgBoard) View() string {
	activePanels := []string{}
	for _, panel := range m.GetActivePanels() {
		activePanels = append(activePanels, panel.View())
	}

	clusterInfo := lipgloss.PlaceHorizontal(m.windowWidth, lipgloss.Center,
		m.clusterInfoPanel.Render(m.windowWidth))

	panels := lipgloss.JoinHorizontal(lipgloss.Top, activePanels...)

	help := lipgloss.PlaceHorizontal(m.windowWidth, lipgloss.Center,
		m.help.View(m.keys))

	return lipgloss.JoinVertical(lipgloss.Left, clusterInfo, panels, help)
}
