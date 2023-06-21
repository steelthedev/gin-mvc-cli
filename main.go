package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Create new Gin", "Add controllers and routes"}

type Model struct {
	cursor      int
	choice      string
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
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.choice == choices[0] {
				m.projectName.Blur()
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
				fmt.Println("Your app " + m.projectName.Value() + "is ready!")
				return m, tea.Quit
			}
			m.choice = choices[m.cursor]
		case "down", "j":
			m.cursor++
			if m.cursor > len(choices) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		default:
			m.projectName, _ = m.projectName.Update(msg)

		}

	case tea.WindowSizeMsg:
		m.projectName.Width = msg.Width

	}

	return m, nil

}

func (m *Model) View() string {

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
	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
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

	dbFile := `
	package db

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("` + projectName + `.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// db.AutoMigrate(&models.` + projectName + `{})

	return db

}
`

	err := os.WriteFile(projectName+"/connections/db/db.go", []byte(dbFile), 0644)
	if err != nil {
		return err
	}

	return nil
}

func createMainFile(projectName string) error {
	mainFile := `package main

	import (
		"net/http"	
		"github.com/gin-gonic/gin"
		"` + projectName + `/app/controllers"
		"` + projectName + `/app/routes"
		"` + projectName + `/connections/db"
		"github.com/rs/cors"
	)
	
	func main() {
		router := gin.Default()
	
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "Welcome to ` + projectName + `"})
		})
	
		router.Use(func(c *gin.Context) {
			cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
			}).ServeHTTP(c.Writer, c.Request, func(w http.ResponseWriter, r *http.Request) {
			})
		})
		

		dbHandler := db.Init()

		/* Routes and controllers */

		` + projectName + `Controller := controllers.New` + capitalizeFirstLetter(projectName) + `Controller(dbHandler)

		routes.Register` + capitalizeFirstLetter(projectName) + `Routes(router, ` + projectName + `Controller)

	
		router.Run(":8000")
	}`

	err := os.WriteFile(projectName+"/main.go", []byte(mainFile), 0644)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultController(projectName string) error {

	defaultController := `package controllers

	import (
		"net/http"

		"github.com/gin-gonic/gin"
		"gorm.io/gorm"
	
	)
	
	type ` + capitalizeFirstLetter(projectName) + `Controller struct {
		db *gorm.DB
	}
	
	func New` + capitalizeFirstLetter(projectName) + `Controller(db *gorm.DB) *` + capitalizeFirstLetter(projectName) + `Controller {
		return &` + capitalizeFirstLetter(projectName) + `Controller{
			db: db,
		}
	}
	
	func (nc *` + capitalizeFirstLetter(projectName) + `Controller) ` + capitalizeFirstLetter(projectName) + `(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":"A new controller for ` + projectName + `",
		})
	}`

	err := os.WriteFile(projectName+"/app/controllers/controllers.go", []byte(defaultController), 0644)
	if err != nil {
		return err
	}
	return nil
}

func createRouteFile(projectName string) error {

	routeFile := `package routes

	import (
		"github.com/gin-gonic/gin"
		"` + projectName + `/app/controllers"
	)
	
	func Register` + capitalizeFirstLetter(projectName) + `Routes(r *gin.Engine, nc *controllers.` + capitalizeFirstLetter(projectName) + `Controller) {
		routes := r.Group("` + projectName + `")
		routes.GET("/", nc.` + capitalizeFirstLetter(projectName) + `)
	}`

	err := os.WriteFile(projectName+"/app/routes/routes.go", []byte(routeFile), 0644)
	if err != nil {
		return err
	}
	return nil
}

func createModelFile(projectName string) error {
	modelFile := `package models
	
	`
	err := os.WriteFile(projectName+"/app/models/models.go", []byte(modelFile), 0644)
	if err != nil {
		return err
	}
	return nil
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

	// write to controllers
	filePath := "app/controllers/controllers.go"

	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		return err
	}
	defer file.Close()

	enterSpace := "\n\n"

	controllerString := `

	` + enterSpace + `

	// ` + appName + ` controllers  
	type ` + capitalizeFirstLetter(appName) + `Controller struct {
		db *gorm.DB
	}
	
	func New` + capitalizeFirstLetter(appName) + `Controller(db *gorm.DB) *` + capitalizeFirstLetter(appName) + `Controller {
		return &` + capitalizeFirstLetter(appName) + `Controller{
			db: db,
		}
	}
	
	func (nc *` + capitalizeFirstLetter(appName) + `Controller) ` + capitalizeFirstLetter(appName) + `(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":"A new controller for ` + appName + `",
		})
	}`
	_, err = io.WriteString(file, controllerString)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err.Error())
		return err
	}

	return nil
}

func AddToRoutes(appName string) error {
	// write to controllers
	filePath := "app/routes/routes.go"

	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		return err
	}
	defer file.Close()

	enterSpace := "\n\n"

	routesString := `

	` + enterSpace + `
		
	func Register` + capitalizeFirstLetter(appName) + `Routes(r *gin.Engine, nc *controllers.` + capitalizeFirstLetter(appName) + `Controller) {
		routes := r.Group("` + appName + `")
		routes.GET("/", nc.` + capitalizeFirstLetter(appName) + `)
	}`

	_, err = io.WriteString(file, routesString)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err.Error())
		return err
	}

	return nil
}
