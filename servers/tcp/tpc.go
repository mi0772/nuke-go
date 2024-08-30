package tcp

import (
	"bufio"
	"fmt"
	"github.com/mi0772/nuke-go/types"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mi0772/nuke-go/engine"
)

const (
	maxGoroutines = 100 // Numero massimo di goroutine nel pool
)

var logf = log.New(os.Stdout, "[TCP-Server] ", log.LstdFlags)

func StartTCPServer(database *engine.Database, config *types.Configuration) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.TpcServerPort))
	if err != nil {
		logf.Fatalf("Errore nel creare il server TCP:%s", err)
		return
	}
	defer listener.Close()

	logf.Println("Server TCP in ascolto su :", config.TpcServerPort)

	jobs := make(chan net.Conn, maxGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go worker(&wg, jobs, database)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logf.Println("Errore nell'accettare la connessione:", err)
			continue
		}
		jobs <- conn
	}

	close(jobs)
	wg.Wait()
}

func handleConnection(conn net.Conn, database *engine.Database) {
	defer conn.Close()

	conn.Write([]byte("hello!\n"))

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(20 * time.Second)) // Timeout per la lettura
		message, err := reader.ReadString('\n')
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				logf.Println("Timeout scaduto, chiusura connessione.")
			} else {
				logf.Println("Errore nella lettura:", err)
			}
			return
		}

		message = strings.TrimSpace(message)

		command, err := NewInputCommand(message)
		if err != nil {
			_, err = conn.Write([]byte("command error\n"))
			continue
		}
		if command.commandIdentifier == Quit {
			break
		}

		cmd, err := CommandBuilder(command)
		if err != nil {
			_, err = conn.Write([]byte("invalid command\n"))
			continue
		}
		item, result := cmd.Process(&command, database)
		if result < 0 {
			conn.Write([]byte(fmt.Sprintf("error:%d", err)))
		} else {
			conn.Write(item)
		}
	}
}

func worker(wg *sync.WaitGroup, jobs <-chan net.Conn, database *engine.Database) {
	defer wg.Done()
	for conn := range jobs {
		handleConnection(conn, database)
	}
}
