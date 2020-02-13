// 2020, GUILLEUS Hugues <ghugues@netc.fr>
// BSD 3-Clause License

package main

import (
	wsSaloon ".."
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("[LAUNCH]")

	// Launch web socket saloon
	saloon := wsSaloon.SaloonNew(func(data []byte) []byte {
		log.Println("[MESSAGE]", string(data))
		return data
	})
	http.Handle("/ws/chatt", saloon)
	// Write a message to all the clients
	go func() {
		t := time.NewTicker(5*time.Second)
		for {
			<-t.C
			saloon.WriteString("Hello world, I'm the server!")
		}
	}()

	// Reste du server web
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
