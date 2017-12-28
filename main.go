package main

import (
	"bufio"
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteHost, err := net.LookupCNAME(conn.RemoteAddr().String())
	if err != nil {
		remoteHost = conn.RemoteAddr().String()
	}
	defer func() {
		log.Printf("Closing connection to: %s (%s).\n", remoteHost, conn.RemoteAddr().String())
		conn.Close()
	}()
	log.Printf("Incoming connection from: %s (%s).\n", remoteHost, conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	slug := generateSlug(time.Now().UnixNano())
	slugFullpath := fmt.Sprintf("%s/%s", viper.GetString("output"), slug)
	log.Printf("Writing file to %s", slugFullpath)

	file, err := os.Create(slugFullpath)
	if err != nil {
		log.Fatalf("Unable to create slug file: %s", slugFullpath)
		conn.Write([]byte(fmt.Sprintf("%s", "Internal error encountered - Please try again later...")))
		return
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		file.WriteString(scanner.Text())
		file.Write([]byte("\n"))
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	}

	stat, _ := file.Stat()
	if stat.Size() <= 0 {
		file.Close()
		log.Printf("0 bytes received, aborting...")
		os.Remove(slugFullpath)
		return
	}

	file.Close()

	log.Printf("Wrote %d bytes to %s", stat.Size(), slugFullpath)

	conn.Write([]byte(fmt.Sprintf("%s://%s/%s\r\n", viper.GetString("uriprefix"), viper.GetString("domain"), slug)))
}
