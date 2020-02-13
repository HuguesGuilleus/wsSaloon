// 2020, GUILLEUS Hugues <ghugues@netc.fr>
// BSD 3-Clause License

package wsSaloon

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// A group of client, when a message is sending by a client, it's sending to
// Received (if non nil) that returns a response for all clients.
type Saloon struct {
	// A callBack called when a message is received.
	Received CBReceived

	// Message to all the clients
	cast chan []byte
	// Message from client
	ms chan []byte

	// List the client conn
	client map[*websocket.Conn]bool
	// Nouvelle connexion
	cr chan *websocket.Conn
	// Remove the connexion to the list of connexion
	rm chan *websocket.Conn

	// Stop the connexion
	stop   chan bool
	stoped bool
}

// Creat a new sallon.
func SaloonNew(r CBReceived) *Saloon {
	s := &Saloon{
		Received: r,
		client:   make(map[*websocket.Conn]bool),
		ms:       make(chan []byte),
		cast:     make(chan []byte),
		cr:       make(chan *websocket.Conn),
		rm:       make(chan *websocket.Conn),
	}
	go s.run()
	return s
}

// run the saloon in background
func (s *Saloon) run() {
	for {
		select {
		case ms := <-s.cast:
			for c := range s.client {
				c.WriteMessage(websocket.TextMessage, ms)
			}
		case c := <-s.cr:
			s.client[c] = true
		case c := <-s.rm:
			if ok := s.client[c]; ok {
				delete(s.client, c)
			}
		}
	}
}

// Create a new client from a HTTP request.
func (s *Saloon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WS ERROR]", err)
		return
	}
	s.cr <- conn

	go func() {
		for {
			t, m, _ := conn.ReadMessage()
			if t == -1 {
				s.rm <- conn
				return
			} else {
				s.Received.run(s, m)
			}
		}
	}()
}

// Write a text message to all the clients.
func (s *Saloon) Write(data []byte) (int, error) {
	s.cast <- data
	return len(data), nil
}

// Write a text message to all the clients.
func (s *Saloon) WriteString(data string) (int, error) {
	s.cast <- []byte(data)
	return len(data), nil
}
