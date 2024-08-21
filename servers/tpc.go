package servers

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/mi0772/nuke-go/engine"
)

const (
	maxGoroutines = 100 // Numero massimo di goroutine nel pool
)

var logf = log.New(os.Stdout, "[TCP-Server] ", log.LstdFlags)

func StartTCPServer(database *engine.Database) {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		logf.Fatalf("Errore nel creare il server TCP:%s", err)
		return
	}
	defer listener.Close()

	logf.Println("Server TCP in ascolto su :9090")

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

	reader := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Timeout per la lettura
	message, err := reader.ReadString('\n')
	if err != nil {
		logf.Println("Errore nella lettura:", err)
		return
	}

	logf.Print("Messaggio ricevuto: ", message)
	conn.Write([]byte("Messaggio ricevuto\n"))
}

func worker(wg *sync.WaitGroup, jobs <-chan net.Conn, database *engine.Database) {
	defer wg.Done()
	for conn := range jobs {
		handleConnection(conn, database)
	}
}
