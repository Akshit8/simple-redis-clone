package server

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/Akshit8/simple-redis-clone/data"
)

type Hub struct {
	// tcp listener
	listener net.Listener

	// channel to listen for quit signal
	quit chan struct{}

	// channel to communicate exit signal
	exit chan struct{}

	// data store
	store *data.Store

	// tcp connections map
	connections map[int]net.Conn

	connTimeout time.Duration
}

func NewServer(store *data.Store) *Hub {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("failed to listen on port 6379: %v", err)
	}

	hub := &Hub{
		listener:    l,
		quit:        make(chan struct{}),
		exit:        make(chan struct{}),
		store:       store,
		connections: make(map[int]net.Conn),
		connTimeout: time.Second * 10,
	}

	go hub.run()

	return hub
}

func (h *Hub) run() {
	log.Println("starting cache server")

	var id int

	for {
		select {
		case <-h.quit:
			log.Println("releasing tcp connection")
			err := h.listener.Close()
			if err != nil {
				log.Printf("failed to close tcp listener: %v", err)
			}

			if len(h.connections) > 0 {
				h.warnAllConnections()
				<-time.After(h.connTimeout)
				h.closeAllConnections()
			}

			close(h.exit)
			return
		default:
			tcpListener, ok := h.listener.(*net.TCPListener)
			if !ok {
				return
			}

			err := tcpListener.SetDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				log.Fatalf("failed to set tcp listener deadline: %v", err)
			}

			conn, err := tcpListener.Accept()
			if err != nil {
				log.Fatalf("failed to accept tcp connection: %v", err)
			}

			h.connections[id] = conn

			go func(id int) {
				log.Printf("client with id [%d] connected\n", id)
				// handle connection
				h.handleConnection(conn)
				delete(h.connections, id)
				log.Printf("client with id [%d] disconnected\n", id)
			}(id)

			id++
		}
	}
}

func (h *Hub) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		pr := parseRequest(scanner.Text())

		switch {
		case len(pr) == 3 && pr[0] == "set":
			h.store.Set(pr[1], pr[2])
			writeToClient(conn, []byte("OK"))
		case len(pr) == 2 && pr[0] == "get":
			val, err := h.store.Get(pr[1])
			if err != nil {
				writeToClient(conn, []byte("ERR:"+err.Error()))
				return
			}

			writeToClient(conn, []byte(val))
		case len(pr) == 2 && pr[0] == "del":
			h.store.Delete(pr[1])
			writeToClient(conn, []byte("OK"))
		case len(pr) == 1 && pr[0] == "quit":
			err := conn.Close()
			if err != nil {
				log.Printf("failed to close tcp connection: %v", err)
			}

			writeToClient(conn, []byte("OK"))
		default:
			writeToClient(conn, []byte("ERR: invalid command"))
		}
	}
}

func (h *Hub) warnAllConnections() {
	for _, conn := range h.connections {
		writeToClient(conn, []byte("server is cloaing soon"))
	}
}

func (h *Hub) closeAllConnections() {
	log.Println("closing all connections")
	for id, conn := range h.connections {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection [%d]: %v\n", id, err)
		}
	}
}

func (h *Hub) Stop() {
	log.Println("stopping cache server")
	close(h.quit)
	<-h.exit
	log.Println("taking data snapshot")
	log.Println("cache server stopped successfully")
}
