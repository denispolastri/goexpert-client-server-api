package main

import (
	"fmt"
	"goexpert-client-server-api/server"
	"net/http"
)

func main() {

	// Define valor HTTP_PORT
	httpPort := "8080"
	//
	http.HandleFunc("/cotacao", server.ConsultaCotacaoSiteEconomia)
	fmt.Println("Server HTTP UP port " + httpPort)
	http.ListenAndServe(":"+httpPort, nil)

}
