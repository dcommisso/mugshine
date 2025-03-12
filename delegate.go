package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"io"
)

const ellipsis = "…"

type MgDelegate struct {
	Styles  MgItemStyles
	height  int
	spacing int
}

func NewMgDelegate() MgDelegate {
	const defaultHeight = 1
	const defaultSpacing = 1
	return MgDelegate{
		Styles:  NewMgItemStyles(),
		height:  defaultHeight,
		spacing: defaultSpacing,
	}
}

type MgDelegateItem interface {
	list.Item
	Title() string
	IsFailed() bool
}

type MgItemStyles struct {
	list.DefaultItemStyles
	FailedTitle         lipgloss.Style
	FailedSelectedTitle lipgloss.Style
}

func NewMgItemStyles() (s MgItemStyles) {
	s.DefaultItemStyles = list.NewDefaultItemStyles()

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(colorSelectedItem).
		Foreground(colorSelectedItem).
		Padding(0, 0, 0, 1)

	s.FailedTitle = lipgloss.NewStyle().
		Foreground(colorFailed).
		Padding(0, 0, 0, 2)

	s.FailedSelectedTitle = lipgloss.NewStyle().
		Foreground(colorFailed).
		Padding(0, 0, 0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(colorSelectedItem))

	return s
}

// Adapted from: https://github.com/charmbracelet/bubbles/blob/cdc743f1f4881343ebc57b3ae820d2dc03331acb/list/defaultitem.go#L140
func (d MgDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title        string
		isFailed     bool
		matchedRunes []int
		s            = &d.Styles
	)

	if i, ok := item.(MgDelegateItem); ok {
		title = i.Title()
		isFailed = i.IsFailed()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered {
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	switch {
	case isFailed && isSelected:
		title = s.FailedSelectedTitle.Render(title)
	case isFailed:
		title = s.FailedTitle.Render(title)
	case emptyFilter:
		title = s.DimmedTitle.Render(title)
	case isSelected && m.FilterState() != list.Filtering:
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
	default:
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
	}

	fmt.Fprintf(w, "%s", title) //nolint: errcheck
}

func (d MgDelegate) Height() int                               { return d.height }
func (d MgDelegate) Spacing() int                              { return d.spacing }
func (d MgDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
