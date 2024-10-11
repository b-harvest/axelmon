package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"bharvest.io/axelmon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var GlobalState *Response

func Run(listenPort int, stateFile string) {
	sf, e := os.OpenFile(stateFile, os.O_RDONLY, 0600)
	if e != nil {
		log.Warn(e.Error())
	}
	b, e := io.ReadAll(sf)
	_ = sf.Close()
	if e != nil {
		log.Warn(e.Error())
	}
	saved := &Response{}
	e = json.Unmarshal(b, saved)

	GlobalState = saved

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GlobalState)
	})

	addr := fmt.Sprintf(":%d", listenPort)
	log.Info(fmt.Sprintf("server listening on %s", addr))

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Error(err)
	}
}
