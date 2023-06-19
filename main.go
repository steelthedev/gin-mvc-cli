package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Create Gin", "Add to project"}

type Model struct {
	cursor int
	choice string
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
			m.choice = choices[m.cursor]
			switch m.choice {
			case choices[0]:
				CreateFolderStructure("test")
			}
			return m, tea.Quit
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

		}

	}
	return m, nil
}

func (m *Model) View() string {
	s := strings.Builder{}
	s.WriteString("Choose a command \n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {

			s.WriteString("(•)")
		} else {
			s.WriteString("( )")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")

	}
	s.WriteString("\n (press q or esc to quit) \n")
	return s.String()
}

func main() {
	p := tea.NewProgram(&Model{})

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(*Model); ok && m.choice != "" {
		fmt.Printf("\n---\nYou chose %s!\n", m.choice)
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

	"` + projectName + `/app/models"
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

	db.AutoMigrate(&models.Welcome{})

	return db

}
`

	err := ioutil.WriteFile(projectName+"/connections/db/db.go", []byte(dbFile), 0644)
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
		WelcomeController := controllers.NewWelcomeController(dbHandler)

		routes.RegisterUserRoutes(router, WelcomeController)

	
		router.Run(":8000")
	}`

	err := ioutil.WriteFile(projectName+"/main.go", []byte(mainFile), 0644)
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
		"` + projectName + `/app/models"
	)
	
	type WelcomeController struct {
		db *gorm.DB
	}
	
	func NewWelcomeController(db *gorm.DB) *WelcomeController {
		return &WelcomeController{
			db: db,
		}
	}
	
	func (wc *WelcomeController) GetWelcome(c *gin.Context) {
		var welcome []models.Welcome
		err := wc.db.Find(&welcome).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, welcome)
	}`

	err := ioutil.WriteFile(projectName+"/app/controllers/welcome_controllers.go", []byte(defaultController), 0644)
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
	
	func RegisterUserRoutes(r *gin.Engine, wc *controllers.WelcomeController) {
		routes := r.Group("welcome")
		routes.GET("/", wc.GetWelcome)
	}`

	err := ioutil.WriteFile(projectName+"/app/routes/routes.go", []byte(routeFile), 0644)
	if err != nil {
		return err
	}
	return nil
}

func createModelFile(projectName string) error {
	modelFile := `package models
	import "gorm.io/gorm"

	
	type Welcome struct {
		gorm.Model
		ID        uint  
		Message   string 
	}
	`
	err := ioutil.WriteFile(projectName+"/app/models/welcome_models.go", []byte(modelFile), 0644)
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
