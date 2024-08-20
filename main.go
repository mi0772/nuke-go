package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/handlers"
	"github.com/mi0772/nuke-go/middleware"
)

var PartFilePath string
var PartitionNumber uint8
var ok bool

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PartFilePath, ok = os.LookupEnv("PART_FILE")
	if !ok {
		log.Fatal("variabile PART_FILE not specified")
	}

	_partitionNumber, ok := os.LookupEnv("PARTITIONS")
	if !ok {
		log.Fatal("variabile PART_FILE not specified")
	}
	num, err := strconv.ParseUint(_partitionNumber, 10, 8)

	if err != nil {
		fmt.Println("Errore:", err)
		return
	}

	PartitionNumber = uint8(num)

	fmt.Println("nuke cache server")
	fmt.Printf("starting server with %d partition, all files will be stored in %s\n", PartitionNumber, PartFilePath)

	database := engine.InitializeDatabase(PartitionNumber, PartFilePath)

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			log.Printf("flushing partition to disk")
			database.FlushPartitions()
		}
	}()

	r := gin.Default()
	r.Use(middleware.DatabaseMiddleware(database))
	r.GET("/keys", handlers.ListKeys)
	r.POST("/push_file", handlers.PushFile)
	r.GET("/pop/:key", handlers.Pop)
	r.GET("/read/:key", handlers.Read)
	r.Run()

}
