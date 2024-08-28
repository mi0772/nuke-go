package main

import (
	"github.com/mi0772/nuke-go/servers/http"
	"github.com/mi0772/nuke-go/servers/tcp"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/types"
)

var logf = log.New(os.Stdout, "[NUKE-Main] ", log.LstdFlags)

func main() {
	err := godotenv.Load()
	if err != nil {
		logf.Fatal("Error loading .env file")
	}

	var configuration = types.ParseConfiguration()

	logf.Println("===================================================")
	logf.Println("=             N U K E  CACHE SERVER               =  ")
	logf.Println("===================================================")
	logf.Printf("starting server with %d partition, all files will be stored in %s\n", configuration.PartitionNumber, configuration.PartitionFilePath)

	database := engine.InitializeDatabase(configuration.PartitionNumber, configuration.PartitionFilePath)

	logf.Printf("total entries in database is : %d", database.CountEntries())

	if configuration.PersistPeriod != -1 {
		ticker := time.NewTicker(time.Duration(configuration.PersistPeriod) * time.Minute)
		logf.Printf("persist will be executed every %d minutes", configuration.PersistPeriod)
		go func() {
			for range ticker.C {
				logf.Printf("flushing partition to disk")
				database.FlushPartitions()
			}
		}()
	} else {
		logf.Printf("persistence disabled !")
	}

	go http.StartHTTPServer(database)
	go tcp.StartTCPServer(database)

	select {}
}
