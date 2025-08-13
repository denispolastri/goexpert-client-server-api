package server

import (
	"goexpert-client-server-api/infrastructure"
	"net/http"
)

func main() {
	infrastructure.NewSqliteDb()
	http.HandleFunc("/cotacao", ConsultaDolar)
	http.ListenAndServe(":8080", nil)
}

func ConsultaDolar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("teste"))
}
