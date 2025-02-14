package scenes

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/ficsit-cli/tea/components"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

var _ tea.Model = (*exitMenu)(nil)

type exitMenu struct {
	root components.RootModel
	list list.Model
}

func NewExitMenu(root components.RootModel) tea.Model {
	model := mainMenu{
		root: root,
	}

	items := []list.Item{
		utils.SimpleItem{
			ItemTitle: "Exit Saving Changes",
			Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
				if err := root.GetGlobal().Save(); err != nil {
					message := "failed to save"
					menu := currentModel.(exitMenu)
					log.Error().Err(err).Msg(message)
					cmd := menu.list.NewStatusMessage(message)
					return currentModel, cmd
				}
				return currentModel, tea.Quit
			},
		},
		utils.SimpleItem{
			ItemTitle: "Exit Discarding Changes",
			Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
				return currentModel, tea.Quit
			},
		},
	}

	model.list = list.NewModel(items, utils.NewItemDelegate(), root.Size().Width, root.Size().Height-root.Height())
	model.list.SetShowStatusBar(false)
	model.list.SetFilteringEnabled(false)
	model.list.Title = "Save Changes?"
	model.list.Styles = utils.ListStyles
	model.list.DisableQuitKeybindings()
	model.list.SetSize(model.list.Width(), model.list.Height())
	model.list.StatusMessageLifetime = time.Second * 3

	return model
}

func (m exitMenu) Init() tea.Cmd {
	return nil
}

func (m exitMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case KeyControlC:
			return m, tea.Quit
		case KeyEnter:
			i, ok := m.list.SelectedItem().(utils.SimpleItem)
			if ok {
				if i.Activate != nil {
					i.Activate(msg, m)
					return m, tea.Quit
				}
			}
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := lipgloss.NewStyle().Margin(2, 2).GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.root.SetSize(msg)
	}

	return m, nil
}

func (m exitMenu) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}
