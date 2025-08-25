package main

import (
	"fmt"
	"goexpert-client-server-api/client"
	"goexpert-client-server-api/server"
	"net/http"
)

func main() {

	// Prepara para chamar a consulta
	var w http.ResponseWriter
	var r *http.Request

	// sobe a aplicação de leitura dos dados
	server.InitServer()
	fmt.Println("InitServer - sobre a aplicação de captura das cotações")

	// Consulta cotação do dolar
	server.ConsultaCotacao(w, r)
	fmt.Println("ConsultaCotacao - ")

	//
	client.LeDolar()

}
