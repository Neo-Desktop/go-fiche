package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/pilu/go-base62"
	"github.com/spf13/viper"
)

func fiche(conn net.Conn) {
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

	slug := base62.Encode(int(time.Now().UnixNano()))
	slugFullpath := fmt.Sprintf("%s/%s", viper.GetString("output"), slug)
	log.Printf("Writing file to %s", slugFullpath)

	file, err := os.Create(slugFullpath)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("%s", "Internal error encountered - Please try again later...\n\n")))
		log.Fatalf("Unable to create slug file: %s", slugFullpath)
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
