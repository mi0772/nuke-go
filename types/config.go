package types

import (
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	PartitionFilePath string
	PartitionNumber   uint8
	HttpServerPort    string
	PersistPeriod     int8
	TpcServerPort     string
	HttpServerEnabled bool
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
	if v, ok := os.LookupEnv("HTTP_PORT"); ok == false {
		result.HttpServerPort = "8080"
	} else {
		result.HttpServerPort = v
	}

	if v, ok := os.LookupEnv("TPC_PORT"); ok == false {
		result.TpcServerPort = "18123"
	} else {
		result.TpcServerPort = v
	}

	if v, ok := os.LookupEnv("HTTP_SERVER"); ok == true {
		if num, err := strconv.ParseInt(v, 10, 8); err != nil {
			log.Fatal("variable HTTP_SERVER must be int")
		} else {
			result.HttpServerEnabled = int8(num) > 0
		}
	} else {
		result.HttpServerEnabled = false
	}

	if _pp, ok := os.LookupEnv("PERSIST_PERIOD"); ok == false {
		result.PersistPeriod = -1
	} else {
		persistPeriod, err := strconv.ParseInt(_pp, 10, 8)
		if err != nil {
			log.Fatal("persist period must be integer")
		} else {
			result.PersistPeriod = int8(persistPeriod)
		}
	}
	return result
}
