package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bharvest.io/axelmon/log"
)

var GlobalState Response

func Run(listenPort int) {
	GlobalState = Response{}

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(GlobalState)
    })

	addr := fmt.Sprintf(":%d", listenPort)
	log.Info(fmt.Sprintf("server listening on %s", addr))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Error(err)
	}
}
