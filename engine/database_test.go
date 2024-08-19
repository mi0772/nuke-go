package engine

import (
	"log"
	"sync"
	"testing"

	"golang.org/x/exp/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestWrite(t *testing.T) {

	wg := sync.WaitGroup{}

	database, err := InitializeDatabase(10, ".")
	if err != nil {
		t.Fatal("error initializing database")
	}

	for i := range 10000 {
		log.Printf("memorized %d item\n", i)
		wg.Add(1)
		key := randStringRunes(20)
		value := randStringRunes(100)
		go memorize(database, key, []byte(value), &wg)
	}

	wg.Wait()
}

func TestResume(t *testing.T) {
	_, err := InitializeDatabase(10, ".")
	if err != nil {
		t.Fatal("error initializing database")
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
