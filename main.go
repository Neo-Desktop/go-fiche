package main

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// AppName is the canonical name of this program
const AppName = "Go-Fiche"

func main() {
	log.Printf("Starting %s on %s...", AppName, time.Now().Format(time.RFC822Z))
	go ficheInit()

	if viper.GetBool("http") {
		go serveHttp()
	}

	for {
		time.Sleep(time.Minute)
	}
}
