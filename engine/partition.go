package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Partition struct {
	entries         map[string]Item
	partitionNumber uint8
	partitionPath   string
	mutex           *sync.RWMutex
}

func (p *Partition) push(key string, value []byte) (Item, error) {
	_, err := p.read(key)
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

func (p *Partition) pop(key string) (Item, error) {
	p.mutex.Lock()
	item, ok := p.entries[key]
	delete(p.entries, key)
	p.mutex.Unlock()
	if ok {
		return item, nil
	} else {
		return Item{}, fmt.Errorf("%s not present", key)
	}
}

func (p *Partition) read(key string) (*Item, error) {
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
