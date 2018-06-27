package main

import (
	"config"
	"db"
	"fmt"
	"github.com/gin-gonic/gin"
	"handlers"
	"io"
	"log"
	"os"
	"path/filepath"
	"utils"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	err = config.InitConfig(dir + "/config.json")
	//err = config.InitConfig("./config.json")
	if err != nil {
		log.Panic(err.Error())
	}

	db.InitDB()
	defer db.DB.Close()

	port := config.AppConfig.Server.Port
	if len(port) == 0 {
		port = ":8080"
	}

	gin.DisableConsoleColor()

	// Logging to a file.
	f, err := os.Create(dir + "/logs/gin.log")
	if err != nil {
		log.Panic(err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(f)

	gin.SetMode(config.AppConfig.Server.ModeStart)

	router := gin.Default()

	api := router.Group("/api")

	v1 := api.Group("/v1")

	usersGroup := v1.Group("/user")

	{
		usersGroup.POST("/token", handlers.CheckEmailHandler, handlers.AddUsersHandler)
		usersGroup.POST("/update", handlers.CheckAuthenticationMiddleware, handlers.UpdateUsersHandler)
	}

	phonesGroup := v1.Group("/phone", handlers.CheckAuthenticationMiddleware)

	{
		phonesGroup.GET("/all", handlers.AllPhonesHandler)
		phonesGroup.GET("/lastupdate", handlers.GetPhonesLastUpdateHandler)
	}

	organizationGroup := v1.Group("/organization", handlers.CheckAuthenticationMiddleware)

	{
		organizationGroup.GET("/all", handlers.AllOrganizationHandler)
	}

	go utils.CheckerFile()

	//log.Fatal(autotls.Run(router, "35.234.94.146"))
	router.Run(port)
}
