package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"sync"
	"time"
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

func (d *Database) Pop(key string) (*Item, error) {
	return d.partitions[d.getPartition(key)].pop(key)
}

func (d *Database) Push(key string, value []byte) (Item, error) {
	return d.partitions[d.getPartition(key)].push(key, value)
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

type Partition struct {
	entries         map[string]Item
	partitionNumber uint8
	partitionPath   string
	mutex           *sync.RWMutex
}

func InitializeDatabase(partitionNumber uint8, filePath string) *Database {
	db := &Database{}
	db.new(partitionNumber, filePath)
	return db
}

func (p *Partition) push(key string, value []byte) (Item, error) {
	_, err := p.pop(key)
	if err != nil {
		i := Item{Key: key, Value: value}
		p.mutex.Lock()
		p.entries[key] = i
		p.mutex.Unlock()
		return i, nil
	} else {
		return Item{}, fmt.Errorf("item with key %s already present", key)
	}
}

func (p *Partition) pop(key string) (*Item, error) {
	p.mutex.RLock()
	item, ok := p.entries[key]
	p.mutex.RUnlock()
	if ok {
		return &item, nil
	} else {
		return nil, fmt.Errorf("%s not present", key)
	}
}

func (p *Partition) persist(wg *sync.WaitGroup) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	defer wg.Done()
	json, err := json.Marshal(p.entries)

	if err != nil {
		log.Fatalf("cannot marshal partition data : %d\n", err)
	}

	file, err := os.Create(p.partitionPath)
	if err != nil {
		log.Fatalf("Errore nella creazione del file: %s", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(json)
	if err != nil {
		log.Fatalf("Errore nella scrittura con il buffer: %s", err)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatalf("Errore nel flush dei dati al file: %s", err)
	}
}

type Item struct {
	Key    string    `json:"key"`
	Value  []byte    `json:"value"`
	Expire time.Time `json:"expire"`
}
