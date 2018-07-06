package main

import (
	"config"
	"context"
	"db"
	"github.com/gin-gonic/gin"
	"handlers"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"utils"
)

func main() {

	utils.ReadConfigFile(false)

	db.InitDB()
	defer db.DB.Close()

	port := config.AppConfig.Server.Port
	if len(port) == 0 {
		port = ":7070"
	}

	gin.DisableConsoleColor()

	// Logging to a file.
	f, err := os.Create(config.AppConfig.RootDirPath + "logs" + string(os.PathSeparator) + "gin.log")
	if err != nil {
		log.Panic(err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	gin.SetMode(config.AppConfig.Server.ModeStart)

	router := gin.Default()

	api := router.Group("/api")

	v1 := api.Group("/v1")

	usersGroup := v1.Group("/user")

	{
		usersGroup.POST("/token", handlers.CheckEmailHandler, handlers.AddUsersHandler)
		usersGroup.POST("/tokenisexist", handlers.CheckAuthenticationMiddleware, handlers.TokenIsExistHandler)
		usersGroup.POST("/update", handlers.CheckAuthenticationMiddleware, handlers.UpdateUsersHandler)
	}

	phonesGroup := v1.Group("/phone", handlers.CheckAuthenticationMiddleware)

	{
		phonesGroup.GET("/all", handlers.AllPhonesHandler)
		phonesGroup.GET("/lastupdate", handlers.GetPhonesLastUpdateHandler)
		phonesGroup.GET("/idorganization/:id", handlers.PhonesByOrganizationIDHandler)
	}

	organizationGroup := v1.Group("/organization", handlers.CheckAuthenticationMiddleware)

	{
		organizationGroup.GET("/all", handlers.AllOrganizationHandler)
	}

	apkGroup := v1.Group("/apk", handlers.CheckAuthenticationMiddleware)

	{
		apkGroup.GET("/lastupdate", handlers.GetLastUpdateAPKHandler)
		apkGroup.GET("/download", handlers.DownloadLastUpdateAPKHandler)
	}

	go utils.CheckerFile()
	go utils.ReadConfigFile(true)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	//log.Fatal(autotls.Run(router, "35.234.94.146"))
	//router.Run(port)
}
