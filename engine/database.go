package engine

import (
	"fmt"
	"hash/fnv"
	"time"
)

type Database struct {
	partitions []Partition
}

func (d *Database) new(partition uint8) {
	d.partitions = make([]Partition, partition)
	for i, p := range d.partitions {
		p.partitionNumber = uint8(i)
		p.entries = *(new(Map[string, Item]))
	}
}

func (d *Database) Pop(key string) (*Item, error) {
	return d.partitions[d.getPartition(key)].pop(key)
}

func (d *Database) Push(key string, value []byte) error {
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

type Partition struct {
	entries         Map[string, Item]
	partitionNumber uint8
}

func InitializeDatabase(path string, partitionNumber uint8) (*Database, error) {
	db := &Database{}
	db.new(partitionNumber)
	return db, nil
}

func (p *Partition) push(key string, value []byte) error {
	_, err := p.pop(key)
	if err != nil {
		i := Item{key: key, value: value}
		p.entries.Store(key, i)
		return nil
	} else {
		return fmt.Errorf("item with key %s already present", key)
	}
}

func (p *Partition) pop(key string) (*Item, error) {
	item, ok := p.entries.Load(key)
	if ok {
		return &item, nil
	} else {
		return nil, fmt.Errorf("%s not present", key)
	}
}

type Item struct {
	key    string    `json:"key"`
	value  []byte    `json:"value"`
	expire time.Time `json:"expire"`
}
