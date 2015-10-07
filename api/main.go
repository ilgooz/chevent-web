package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codeui/chevent-web/api/conf"
)

func main() {
	conf.Load()
	run()
}

func run() {
	server := &http.Server{
		Addr:    conf.Addr,
		Handler: handler(),
	}

	fmt.Printf("api server started at: %s\n", conf.Addr)
	log.Fatalln(server.ListenAndServe())
}
