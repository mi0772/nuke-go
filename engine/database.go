package engine

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mi0772/nuke-go/types"
)

type Database struct {
	partitions []Partition
	pathFile   string
}

func (d *Database) new(partition uint8, pathFile string) {
	d.pathFile = pathFile
	d.partitions = make([]Partition, partition)

	for i, p := range d.partitions {
		p.partitionNumber = uint8(i)
		p.entries = make(map[string]Item)
		p.partitionPath = fmt.Sprintf("%s/partition_%d.part", d.pathFile, i)
		p.mutex = &sync.RWMutex{}

		//try to load partitions from file
		if _, err := os.Stat(p.partitionPath); err == nil {
			fmt.Printf("File exists\n")
			file, err := os.ReadFile(p.partitionPath)
			if err != nil {
				log.Fatalf("impossible read file partition %s", p.partitionPath)
			}
			err = json.Unmarshal(file, &p.entries)
			if err != nil {
				log.Fatalf("impossibile unmarshal json content : %s", p.partitionPath)
			}
			log.Printf("successfully resumed partition : %d\n", p.partitionNumber)
		}

		d.partitions[i] = p
	}
}

func (d *Database) Pop(key string) (Item, error) {
	return d.partitions[d.getPartition(key)].pop(key)
}

func (d *Database) Push(key string, value []byte) (Item, error) {
	return d.partitions[d.getPartition(key)].push(key, value)
}

func (d *Database) Read(key string) (*Item, error) {
	return d.partitions[d.getPartition(key)].read(key)
}

func (d *Database) Keys() []types.Key {
	keySize := 0
	for _, p := range d.partitions {
		keySize += len(p.entries)
	}

	keys := make([]types.Key, 0, keySize)

	for _, p := range d.partitions {
		for k := range p.entries {
			key := types.Key{Key: k, Partition: p.partitionNumber}
			keys = append(keys, key)
		}
	}

	return keys
}

func (d *Database) getPartition(key string) uint8 {
	hasher := fnv.New64a()
	hasher.Write([]byte(key))
	hash := hasher.Sum64()
	numPartitions := uint64(len(d.partitions))
	result := uint8(hash % numPartitions)
	fmt.Printf("partizione selezionata : %d\n", result)
	return result
}

func (d *Database) FlushPartitions() {
	wg := sync.WaitGroup{}
	for _, p := range d.partitions {
		wg.Add(1)
		go p.persist(&wg)
	}
	wg.Wait()
	log.Printf("successfully persisted %d partitions", len(d.partitions))
}

func InitializeDatabase(partitionNumber uint8, filePath string) *Database {
	db := &Database{}
	db.new(partitionNumber, filePath)
	return db
}

type Item struct {
	Key    string    `json:"key"`
	Value  []byte    `json:"value"`
	Expire time.Time `json:"expire"`
}
