package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const AppName = "Go-Fiche"

const slugMap = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

var (
	GitCommit,
	GitBranch,
	GitState,
	BuildDate,
	Version string
)

func init() {
	pflag.BoolP("help", "h", false, "Prints this help message")
	pflag.StringP("output", "o", "./code", "Relative or absolute path to the directory where you want to store user-posted pastes.")
	pflag.StringP("domain", "d", "localhost", "This will be used as a prefix for an output received by the client. Value will be prepended with http[s].")
	pflag.IntP("port", "p", 9999, "Port in which the service should listen on.")
	pflag.BoolP("https", "S", false, fmt.Sprintf("If set, %s returns url with https prefix instead of http.", AppName))
	pflag.IntP("buffer", "B", 32768, "This parameter defines size of the buffer used for getting data from the user. Maximum size (in bytes) of all input files is defined by this value.")
	pflag.StringP("log", "l", "", "Log file. This file has to be user-writable.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("help") {
		fmt.Printf("%s! - Version %s, Built on %s from Git tag [%s:%s-%s)\n", AppName, Version, BuildDate, GitBranch, GitCommit, GitState)
		pflag.Usage()
		os.Exit(2)
	}

	if viper.GetBool("https") {
		viper.Set("uriprefix", "https")
	} else {
		viper.Set("uriprefix", "http")
	}
}

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
	time.Sleep(5)
	os.Exit(0)
}

func generateSlug(seed int64) string {
	stringSeed := fmt.Sprintf("%d", seed)
	evenLength := (len(stringSeed) / 2) * 2
	digitHold := int64(0)
	out := ""
	for i := 0; i < evenLength; i += 2 {
		digitHold, _ = strconv.ParseInt(stringSeed[i:i+2], 10, 8)
		digitHold = digitHold % int64(len(slugMap))
		out += slugMap[digitHold : digitHold+1]
	}
	return out
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
	} else {
		file.Close()
	}

	log.Printf("Wrote %d bytes to %s", stat.Size(), slugFullpath)

	conn.Write([]byte(fmt.Sprintf("%s://%s/%s\r\n", viper.GetString("uriprefix"), viper.GetString("domain"), slug)))
}
