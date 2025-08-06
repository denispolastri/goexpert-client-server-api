package main

import (
	"goexpert-client-server-api/server"
	"net/http"
)

func main() {
	var w http.ResponseWriter
	var r *http.Request
	server.ConsultaDolar(w, r)
}
