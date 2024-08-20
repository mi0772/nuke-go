package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/handlers"
	"github.com/mi0772/nuke-go/middleware"
	"github.com/mi0772/nuke-go/types"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var configuration = types.ParseConfiguration()

	log.Println("===================================================")
	log.Println("= N U K E  CACHE SERVER")
	log.Println("===================================================")
	log.Printf("starting server with %d partition, all files will be stored in %s\n", configuration.PartitionNumber, configuration.PartitionFilePath)

	database := engine.InitializeDatabase(configuration.PartitionNumber, configuration.PartitionFilePath)

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			log.Printf("flushing partition to disk")
			database.FlushPartitions()
		}
	}()

	r := gin.Default()
	r.Use(middleware.DatabaseMiddleware(database))
	r.GET("/admin/keys", handlers.ListKeys)
	r.DELETE("/admin/clear", handlers.Clear)
	r.GET("/admin/partitions/details", handlers.PartitionDetails)
	r.POST("/push_file", handlers.PushFile)
	r.POST("/push_string", handlers.PushString)
	r.GET("/pop/:key", handlers.Pop)
	r.GET("/read/:key", handlers.Read)
	r.Run()

}
