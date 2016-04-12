package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func main() {
	wsContainer := doInit()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.BindAddress, config.BindPort),
		wsContainer))
}
