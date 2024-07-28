package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type CollectApp struct {
	focusIndex int
	grimoire   Grimoire
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel(g Grimoire) CollectApp {
	app := CollectApp{
		inputs: make([]textinput.Model, 3),
		grimoire: g,
	}

	var t textinput.Model
	for i := range app.inputs {

		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64
		switch i {
		case 0:
			t.Placeholder = "Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Command"
		}

		app.inputs[i] = t

	}

	return app
}

func (app CollectApp) Init() tea.Cmd {
	return textinput.Blink
}

func (app CollectApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return app, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			app.cursorMode++
			if app.cursorMode > cursor.CursorHide {
				app.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(app.inputs))
			for i := range app.inputs {
				cmds[i] = app.inputs[i].Cursor.SetMode(app.cursorMode)
			}
			return app, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && app.focusIndex == len(app.inputs) {

				if len(app.inputs[0].Value()) != 0 && len(app.inputs[2].Value()) != 0 {
					spell, err := NewSpell(app.inputs[1].Value(), app.inputs[2].Value())
					if err != nil {
						log.Fatal("There was an error creating a spell", err)
						return app, tea.Quit

					}
					app.grimoire.AddSpell(app.inputs[0].Value(), spell)
					app.grimoire.FlushToFile()
				}

				return app, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				app.focusIndex--
			} else {
				app.focusIndex++
			}

			if app.focusIndex > len(app.inputs) {
				app.focusIndex = 0
			} else if app.focusIndex < 0 {
				app.focusIndex = len(app.inputs)
			}

			cmds := make([]tea.Cmd, len(app.inputs))
			for i := 0; i <= len(app.inputs)-1; i++ {
				if i == app.focusIndex {
					// Set focused state
					cmds[i] = app.inputs[i].Focus()
					app.inputs[i].PromptStyle = focusedStyle
					app.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				app.inputs[i].Blur()
				app.inputs[i].PromptStyle = noStyle
				app.inputs[i].TextStyle = noStyle
			}

			return app, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := app.updateInputs(msg)

	return app, cmd
}

func (app *CollectApp) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(app.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	///	for i := range app.inputs {
	///		app.inputs[i], cmds[i] = app.inputs[i].Update(msg)
	///	}
	app.inputs[0], cmds[0] = app.inputs[0].Update(msg)
	app.inputs[1], cmds[1] = app.inputs[1].Update(msg)
	app.inputs[2], cmds[2] = app.inputs[2].Update(msg)

	return tea.Batch(cmds...)
}

func (app CollectApp) View() string {
	var b strings.Builder

	b.WriteString(app.inputs[0].View())
	b.WriteRune('\n')
	b.WriteString(app.inputs[1].View())
	b.WriteRune('\n')
	b.WriteString(app.inputs[2].View())
	b.WriteRune('\n')

	button := &blurredButton
	if app.focusIndex == len(app.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(app.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
