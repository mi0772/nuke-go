package engine

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"golang.org/x/exp/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestWrite(t *testing.T) {

	wg := sync.WaitGroup{}

	database := InitializeDatabase(10, ".")

	for i := range 10000 {
		log.Printf("memorized %d item\n", i)
		wg.Add(1)
		key := randStringRunes(20)
		value := randStringRunes(100)
		go memorize(database, key, []byte(value), &wg)
	}

	wg.Wait()
}

func TestDatabaseNew(t *testing.T) {
	// Creare una directory temporanea per i test
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir) // Pulire la directory temporanea dopo il test

	// Inizializzare una nuova istanza di Database
	db := &Database{}

	// Chiamare il metodo new con 3 partizioni
	db.new(3, tempDir)

	// Verificare che il numero di partizioni sia corretto
	if len(db.partitions) != 3 {
		t.Errorf("Numero di partizioni errato, ottenuto: %d, atteso: %d", len(db.partitions), 3)
	}

	// Verificare che il percorso del file sia stato impostato correttamente
	if db.pathFile != tempDir {
		t.Errorf("Percorso del file errato, ottenuto: %s, atteso: %s", db.pathFile, tempDir)
	}

	// Verificare che ogni partizione sia stata inizializzata correttamente
	for i, p := range db.partitions {
		expectedPath := fmt.Sprintf("%s/partition_%d.part", tempDir, i)
		if p.partitionPath != expectedPath {
			t.Errorf("Percorso della partizione errato, ottenuto: %s, atteso: %s", p.partitionPath, expectedPath)
		}
		if p.partitionNumber != uint8(i) {
			t.Errorf("Numero della partizione errato, ottenuto: %d, atteso: %d", p.partitionNumber, i)
		}
		if p.entries == nil {
			t.Errorf("Mappa delle entries non inizializzata per la partizione %d", i)
		}
		if p.mutex == nil {
			t.Errorf("Mutex non inizializzato per la partizione %d", i)
		}
	}
}

func TestResume(t *testing.T) {
	_ = InitializeDatabase(10, ".")
}
func TestCountEntries(t *testing.T) {
	database := InitializeDatabase(10, ".")

	// Add some entries to the database
	database.Push("key1", []byte("value1"))
	database.Push("key2", []byte("value2"))
	database.Push("key3", []byte("value3"))

	// Get the count of entries
	count := database.CountEntries()

	// Verify the count is correct
	expectedCount := uint(3)
	if count != expectedCount {
		t.Errorf("CountEntries returned %d, expected %d", count, expectedCount)
	}
}
func memorize(database *Database, key string, value []byte, wg *sync.WaitGroup) {
	database.Push(key, []byte(value))
	wg.Done()
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
