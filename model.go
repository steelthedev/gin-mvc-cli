package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type GinMVCModel struct {
	cursor   int
	choice   TUIChoice
	tuiInput textinput.Model
}

func NewGinMVCModel(tuiInput textinput.Model) *GinMVCModel {
	return &GinMVCModel{
		tuiInput: tuiInput,
	}
}

func (g *GinMVCModel) Init() tea.Cmd {
	return nil
}

func (g *GinMVCModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return g.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		g.tuiInput.Width = msg.Width

	}
	return g, nil
}

func (g *GinMVCModel) View() string {

	if g.choice == TUIChoiceCreateNewGin {

		g.tuiInput.Placeholder = "Input project name"
		g.tuiInput.Focus()
		return g.tuiInput.View()

	}

	if g.choice == TUIChoiceAddControllersAndRoutes {
		g.tuiInput.Placeholder = "Input app name"
		g.tuiInput.Focus()
		return g.tuiInput.View()
	}
	s := strings.Builder{}
	s.WriteString("What do you want to do?\n\n")

	choices := []TUIChoice{TUIChoiceCreateNewGin, TUIChoiceAddControllersAndRoutes}
	for i := 0; i < ChoicesCount; i++ {
		if g.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(choices[i]))
		s.WriteString("\n")
	}
	s.WriteString("\n(press esc or ctrl+c to quit)\n")

	return s.String()
}

func (g *GinMVCModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return g, tea.Quit
	case "enter":
		if g.choice == TUIChoiceCreateNewGin {
			return g.handleTUIChoiceCreateNewGin()
		}
		if g.choice == TUIChoiceAddControllersAndRoutes {
			return g.handleTUIChoiceAddControllersAndRoutes()
		}
		g.setChoice()
	case "down":
		g.cursor++
		if g.cursor > ChoicesCount {
			g.cursor = 0
		}
	case "up":
		g.cursor--
		if g.cursor < 0 {
			g.cursor = ChoicesCount - 1
		}
	default:
		if g.tuiInput.Focused() {
			g.tuiInput, _ = g.tuiInput.Update(msg)
		}
	}
	return g, nil
}

func (g *GinMVCModel) handleTUIChoiceCreateNewGin() (tea.Model, tea.Cmd) {
	g.tuiInput.Blur()
	s, err := validateInput(g.tuiInput.Value())
	if err != nil {
		fmt.Println(err)
		return g, tea.Quit
	}
	s = formatInput(s)
	err = CreateProject(s)
	if err != nil {
		fmt.Println(err)
		return g, tea.Quit
	}
	fmt.Println("Your project " + g.tuiInput.Value() + " is ready!")
	return g, tea.Quit
}

func (g *GinMVCModel) handleTUIChoiceAddControllersAndRoutes() (tea.Model, tea.Cmd) {
	g.tuiInput.Blur()
	s, err := validateInput(g.tuiInput.Value())
	if err != nil {
		fmt.Println(err)
		return g, tea.Quit
	}
	s = formatInput(s)
	err = AddToExistingProject(s)
	if err != nil {
		fmt.Println(err)
		return g, tea.Quit
	}
	fmt.Println("Your app " + g.tuiInput.Value() + " is ready!")
	return g, tea.Quit
}

func (g *GinMVCModel) setChoice() {
	switch g.cursor {
	case 0:
		g.choice = TUIChoiceCreateNewGin
	case 1:
		g.choice = TUIChoiceAddControllersAndRoutes
	}
}
