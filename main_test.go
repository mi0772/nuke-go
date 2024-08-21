package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/types"
	"golang.org/x/exp/rand"
)

func TestInsert(t *testing.T) {
	wg := sync.WaitGroup{}
	var configuration = types.ParseConfiguration()

	database := engine.InitializeDatabase(configuration.PartitionNumber, configuration.PartitionFilePath)

	for i := 0; i < 1_000_000; i++ {
		wg.Add(1)
		go func() {
			database.Push(fmt.Sprintf("key_%d", i), []byte(RandStringRunes(100)))
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("done")
}

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
