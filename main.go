package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	database, error := engine.InitializeDatabase(PartitionNumber, PartFilePath)
	if error != nil {
		panic("ciao")
	}

	_, e := database.Push("carlo", []byte("Here is a string...."))
	if e != nil {
		fmt.Println(e)
	}
	_, e = database.Push("carlo", []byte("secondo carlo"))
	if e != nil {
		fmt.Println(e)
	}
	_, e = database.Push("antonio", []byte("Here is a string...."))
	if e != nil {
		fmt.Println(e)
	}
	_, e = database.Push("banana", []byte("Here is a string...."))
	if e != nil {
		fmt.Println(e)
	}

	v, e := database.Pop("carlo")
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("il valore trovato è : %s\n", v)
	}

}
