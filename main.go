package main

import (
	"embed"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"text/template"
)

//go:embed templates/*
var f embed.FS
var templateEngine *template.Template

func init() {
	templateEngine = template.New("gin-mvc-cli-template-renderer")
}

type TUIChoice string

const (
	TUIChoiceCreateNewGin            TUIChoice = "Create new Gin"
	TUIChoiceAddControllersAndRoutes TUIChoice = "Add controllers and routes"
	ChoicesCount                               = 2
)

func main() {
	p := tea.NewProgram(NewGinMVCModel(textinput.New()))

	_, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
