package main

import (
	"database/sql"
	"fmt"
	models "gpus/my_models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type GpusData struct {
	UserName string `json:"username"`
	Comments string `json:"comments"`
}

func handleAddToGpu(c *gin.Context) {
	data := GpusData{}
	if err := c.BindJSON(&data); err != nil {
         fmt.Println("error in binding!!")
		 return;
	}

	var insertData models.Gpu;
	insertData.Username = data.UserName
	insertData.Comment = data.Comments

	err := insertData.Insert(c, boil.GetContextDB(), boil.Infer())
	
	if err != nil {
		fmt.Println("error in inserting data : ",err)
		return 
	}
	
	c.JSON(http.StatusAccepted,gin.H{"success" : true, "data" : insertData})
	
}

func handleGetGpu(c *gin.Context) {
	Data, err := models.Gpus().All(c, boil.GetContextDB())
	if err != nil {
		fmt.Println("error in getting data : ",err)
		return;
	}
	c.JSON(http.StatusAccepted,gin.H{"success" : true, "data" : Data})
}

func connectDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gpu?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations() error {
	m, err := migrate.New(
		"file://db/migrations", "postgres://postgres:postgres@localhost:5432/gpu?sslmode=disable")

	if err != nil {
		return err
	}
	fmt.Println("Migrations up successfully")
	m.Up()

	return nil
}

func main() {
	DB, err := connectDatabase()
	if err != nil {
		fmt.Println("error in connection database")
		return
	}
	boil.SetDB(DB)
	fmt.Println("database Connected successfully")

	err = RunMigrations()
	if err != nil {
		fmt.Println("Failed to migrate")
		fmt.Println("error in migration is : ", err)
		return
	}

	router := gin.Default()
	router.POST("/Add", handleAddToGpu)
	router.GET("/show", handleGetGpu)

	router.Run(":8080")
}
