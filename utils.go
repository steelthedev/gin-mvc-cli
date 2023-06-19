package main

// import (
// 	"io/ioutil"
// 	"os"
// )

// func CreateFolderStructure(projectName string) error {
// 	// Create the main project folder
// 	err := os.Mkdir(projectName, 0755)
// 	if err != nil {
// 		return err
// 	}

// 	// Create the app folder
// 	err = os.Mkdir(projectName+"/app", 0755)
// 	if err != nil {
// 		return err
// 	}

// 	// Create the controllers folder
// 	err = os.Mkdir(projectName+"/app/controllers", 0755)
// 	if err != nil {
// 		return err
// 	}

// 	// Create the models folder
// 	err = os.Mkdir(projectName+"/app/models", 0755)
// 	if err != nil {
// 		return err
// 	}

// 	// Create the routes folder
// 	err = os.Mkdir(projectName+"/app/routes", 0755)
// 	if err != nil {
// 		return err
// 	}

// 	//create connections folder
// 	err = os.Mkdir(projectName+"/connections", 0755)
// 	if err != nil {
// 		return err
// 	}

// 	//Create db folder in connections
// 	err = os.Mkdir(projectName+"/connections/db", 0755)

// 	err = createDbFile(projectName)
// 	if err != nil {
// 		return err
// 	}
// 	err = createMainFile(projectName)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func createDbFile(projectName string) error {

// 	dbFile := `
// 	package db

// import (
// 	"log"

// 	"` + projectName + `/models"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func Init() *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open("` + projectName + `.db"), &gorm.Config{
// 		DisableForeignKeyConstraintWhenMigrating: true,
// 	})

// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	//db.AutoMigrate(&models.Model{})

// 	return db

// }
// `

// 	err := ioutil.WriteFile(projectName+"/connections/db/db.go", []byte(dbFile), 0644)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func createMainFile(projectName string) error {
// 	mainFile := `package main

// 	import (
// 		"log"
// 		"net/http"

// 		"github.com/gin-gonic/gin"
// 		"gorm.io/driver/sqlite"
// 		"gorm.io/gorm"
// 		"` + projectName + `/app/controllers"
// 		"` + projectName + `/app/models"
// 		"` + projectName + `/app/routes"
// 		"` + projectName + `/connections/db"
// 	)

// 	func main() {
// 		router := gin.Default()
// 		dbHandler := db.Init(dbURL)

// 		router.GET("/", func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"data": "Welcome to ` + projectName + `"})
// 		})

// 		router.Use(func(c *gin.Context) {
// 			cors.New(cors.Options{
// 				AllowedOrigins: []string{"*"},
// 				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
// 				AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 			}).ServeHTTP(c.Writer, c.Request, func(w http.ResponseWriter, r *http.Request) {
// 			})
// 		})

// 		router.Run(":8000")
// 	}`

// 	err := ioutil.WriteFile(projectName+"/main.go", []byte(mainFile), 0644)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
