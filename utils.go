package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CreateProject(projectName string) error {

	if err := createProjectDirectories(projectName); err != nil {
		return err
	}

	if err := createProjectFiles(projectName); err != nil {
		return err
	}

	if err := initializeProject(projectName); err != nil {
		fmt.Printf("Error initializing project: %s\n", err.Error())
		return err
	}
	return nil
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

func createProjectDirectories(projectName string) error {
	var directories = []struct {
		name string
		perm os.FileMode
	}{
		// main project folder
		{name: projectName, perm: 0755},
		// app folder
		{name: projectName + "/app", perm: 0755},
		// controllers folder
		{name: projectName + "/app/controllers", perm: 0755},
		// models folder
		{name: projectName + "/app/models", perm: 0755},
		// routes folder
		{name: projectName + "/app/routes", perm: 0755},
		// connections folder
		{name: projectName + "/connections", perm: 0755},
		// db folder
		{name: projectName + "/connections/db", perm: 0755},
	}

	for _, dir := range directories {
		err := os.Mkdir(dir.name, dir.perm)
		if err != nil {
			return err
		}
	}
	return nil
}

func createProjectFiles(projectName string) error {
	if err := createDbFile(projectName); err != nil {
		return err
	}
	if err := createDefaultController(projectName); err != nil {
		return err
	}
	if err := createRouteFile(projectName); err != nil {
		return err
	}
	if err := createModelFile(projectName); err != nil {
		return err
	}
	if err := createMainFile(projectName); err != nil {
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
	fmt.Println("Initializing project", projectName)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	fmt.Println("Downloading dependencies for project", projectName, "please wait...")
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

func validateInput(input string) (string, error) {
	// Check if the string contains '/' or '-'
	if strings.ContainsAny(input, "/-") {
		return "", errors.New("string contains invalid characters")
	}
	return input, nil
}

func formatInput(input string) string {
	// Remove spaces between words
	formattedString := strings.Join(strings.Fields(input), "")
	return formattedString
}
