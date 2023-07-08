package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"embed"
	"text/template"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed templates/*
var f embed.FS
var templateEngine *template.Template
var choices = []string{"Create new Gin", "Add controllers and routes"}

func init() {
	templateEngine = template.New("gin-mvc-cli-template-renderer")
}

func main() {
	p := tea.NewProgram(&Model{
		projectName: textinput.New(),
	})

	// Run returns the model as a tea.Model.
	_, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

}

type Model struct {
	cursor      int
	choice      string
	loaded      bool
	projectName textinput.Model
}

func validateAndFormatString(input string) (string, error) {
	// Check if the string contains '/' or '-'
	if strings.ContainsAny(input, "/-") {
		return "", errors.New("string contains invalid characters")
	}

	// Remove spaces between words
	formattedString := strings.Join(strings.Fields(input), "")

	return formattedString, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.choice == choices[0] {
				m.projectName.Blur()
				m.loaded = true
				s, err := validateAndFormatString(m.projectName.Value())
				if err != nil {
					fmt.Println(err)
					return m, tea.Quit
				}

				err = CreateFolderStructure(s)
				if err != nil {
					fmt.Println(err)
					return m, tea.Quit
				}

				m.loaded = false
				fmt.Println("Your project " + m.projectName.Value() + " is ready!")
				return m, tea.Quit
			}

			if m.choice == choices[1] {
				m.projectName.Blur()
				s, err := validateAndFormatString(m.projectName.Value())
				if err != nil {
					fmt.Println(err)
					return m, tea.Quit
				}
				err = AddToExistingProject(s)
				if err != nil {
					fmt.Println(err)
					return m, tea.Quit
				}
				m.loaded = true
				fmt.Println("Your app " + m.projectName.Value() + " is ready!")
				return m, tea.Quit
			}
			m.choice = choices[m.cursor]
		case "down":
			m.cursor++
			if m.cursor > len(choices) {
				m.cursor = 0
			}
		case "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		default:
			if m.projectName.Focused() {
				m.projectName, _ = m.projectName.Update(msg)
			}

		}

	case tea.WindowSizeMsg:
		m.projectName.Width = msg.Width

	}

	return m, nil

}

func (m *Model) View() string {

	if m.loaded {
		fmt.Println("Loading...........")
	}

	if m.choice == choices[0] {

		m.projectName.Placeholder = "Input project name"
		m.projectName.Focus()
		return m.projectName.View()

	}

	if m.choice == choices[1] {
		m.projectName.Placeholder = "Input app name"
		m.projectName.Focus()
		return m.projectName.View()
	}
	s := strings.Builder{}
	s.WriteString("What do you want to do?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press esc to quit)\n")

	return s.String()
}

func CreateFolderStructure(projectName string) error {
	// Create the main project folder
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		return err
	}

	// Create the app folder
	err = os.Mkdir(projectName+"/app", 0755)
	if err != nil {
		return err
	}

	// Create the controllers folder
	err = os.Mkdir(projectName+"/app/controllers", 0755)
	if err != nil {
		return err
	}

	// Create the models folder
	err = os.Mkdir(projectName+"/app/models", 0755)
	if err != nil {
		return err
	}

	// Create the routes folder
	err = os.Mkdir(projectName+"/app/routes", 0755)
	if err != nil {
		return err
	}

	//create connections folder
	err = os.Mkdir(projectName+"/connections", 0755)
	if err != nil {
		return err
	}

	//Create db folder in connections
	err = os.Mkdir(projectName+"/connections/db", 0755)

	err = createDbFile(projectName)
	if err != nil {
		return err
	}

	err = createDefaultController(projectName)
	if err != nil {
		return err
	}

	err = createRouteFile(projectName)
	if err != nil {
		return err
	}

	err = createModelFile(projectName)
	if err != nil {
		return err
	}

	err = createMainFile(projectName)
	if err != nil {
		return err
	}

	err = initializeProject(projectName)
	if err != nil {
		fmt.Printf("Error initializing project: %s\n", err.Error())
		return err
	}
	return nil

}

func createDbFile(projectName string) error {
	return createFile(projectName+"/connections/db/db.go", os.O_RDWR|os.O_CREATE, "templates/db.txt", struct{ ProjectName string }{ProjectName: projectName})
}

func createMainFile(projectName string) error {
	data := struct {
		ProjectName            string
		ProjectNameCapitalized string
	}{
		ProjectName:            projectName,
		ProjectNameCapitalized: capitalizeFirstLetter(projectName),
	}
	return createFile(projectName+"/main.go", os.O_RDWR|os.O_CREATE, "templates/main.txt", data)
}

func createDefaultController(projectName string) error {
	data := struct {
		ProjectName            string
		ProjectNameCapitalized string
	}{
		ProjectName:            projectName,
		ProjectNameCapitalized: capitalizeFirstLetter(projectName),
	}
	return createFile(projectName+"/app/controllers/controllers.go", os.O_RDWR|os.O_CREATE, "templates/controller.txt", data)
}

func createRouteFile(projectName string) error {
	data := struct {
		ProjectName            string
		ProjectNameCapitalized string
	}{
		ProjectName:            projectName,
		ProjectNameCapitalized: capitalizeFirstLetter(projectName),
	}
	return createFile(projectName+"/app/routes/routes.go", os.O_RDWR|os.O_CREATE, "templates/routes.txt", data)
}

func createModelFile(projectName string) error {
	return createFile(projectName+"/app/models/models.go", os.O_RDWR|os.O_CREATE, "templates/models.txt", nil)
}

func initializeProject(projectName string) error {

	err := os.Chdir(projectName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "init", projectName)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func capitalizeFirstLetter(word string) string {
	if len(word) == 0 {
		return word
	}

	firstLetter := strings.ToUpper(string(word[0]))
	restOfWord := word[1:]
	return firstLetter + restOfWord
}

func AddToExistingProject(appName string) error {
	err := AddToController(appName)
	if err != nil {
		return err
	}

	err = AddToRoutes(appName)
	if err != nil {
		return err
	}

	return nil
}

func AddToController(appName string) error {
	data := struct {
		AppName            string
		AppNameCapitalized string
	}{
		AppName:            appName,
		AppNameCapitalized: capitalizeFirstLetter(appName),
	}
	return createFile("app/controllers/controllers.go", os.O_WRONLY|os.O_APPEND, "templates/add_controller.txt", data)
}

func AddToRoutes(appName string) error {
	data := struct {
		AppName            string
		AppNameCapitalized string
	}{
		AppName:            appName,
		AppNameCapitalized: capitalizeFirstLetter(appName),
	}
	return createFile("app/routes/routes.go", os.O_WRONLY|os.O_APPEND, "templates/add_router.txt", data)
}

func createFile(fileName string, flag int, template string, templateData interface{}) error {
	file, err := os.OpenFile(fileName, flag, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	content, err := f.ReadFile(template)
	if err != nil {
		return err
	}
	fileTemplate, err := templateEngine.Parse(string(content))
	if err != nil {
		return err
	}
	err = fileTemplate.Execute(file, templateData)
	return err
}
