package types

import (
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	PartitionFilePath string
	PartitionNumber   uint8
	ServerPort        string
}

func ParseConfiguration() Configuration {
	var result = Configuration{}
	var ok bool

	if result.PartitionFilePath, ok = os.LookupEnv("PART_FILE"); ok == false {
		result.PartitionFilePath = "./part_file"
	}

	if _partitionNumber, ok := os.LookupEnv("PARTITIONS"); ok == false {
		result.PartitionNumber = 10
	} else {
		if num, err := strconv.ParseUint(_partitionNumber, 10, 8); err != nil {
			log.Fatal("variable PARTITIONS must be unsigned int")
		} else {
			result.PartitionNumber = uint8(num)
		}
	}
	if v, ok := os.LookupEnv("PORT"); ok == false {
		result.ServerPort = "8080"
	} else {
		result.ServerPort = v
	}
	return result
}
