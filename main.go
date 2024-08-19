package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mi0772/nuke-go/engine"
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

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			log.Printf("flushing partition to disk")
			database.FlushPartitions()
		}
	}()
	time.Sleep(60 * time.Second)
	ticker.Stop()

}
