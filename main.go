package main

import (
	"goexpert-client-server-api/infrastructure"
	"goexpert-client-server-api/server"
	"net/http"
)

func main() {

	// Inicializa o banco de dados
	db, err := infrastructure.NewSqliteDb()
	if err != nil {
		panic(err)
	}

	// Inicia o servidor
	server.InitServer()

	// Prepara para chamar a consulta
	var w http.ResponseWriter
	var r *http.Request

	// Consulta cotação do dolar
	server.ConsultaCotacao(w, r)

	server.GravaCotacao(db)

}
