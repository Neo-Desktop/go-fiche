package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/viper"
)

// AppName is the canonical name of this program
const AppName = "Go-Fiche"

func main() {
	log.Printf("Starting %s on %s...", AppName, time.Now().Format(time.RFC822Z))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.Get("port")))
	if err != nil {
		// handle error
		log.Fatalf("Could not bind to port: %d!", viper.Get("port"))
		os.Exit(-1)
	}
	log.Printf("Server started listening on port: %d.", viper.Get("port"))
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error on accepting connection!")
			continue
		}
		go fiche(conn)
	}
}
