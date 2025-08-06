package server

import (
	"net/http"
)

func init() {
	http.HandleFunc("/cotacao", ConsultaDolar)
	http.ListenAndServe(":8080", nil)
}

func ConsultaDolar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("teste"))
}
