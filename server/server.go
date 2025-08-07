package server

import (
	"net/http"

	"gorm.io/gorm"
)

func InitServer() {
	http.HandleFunc("/cotacao", ConsultaCotacao)
	http.ListenAndServe(":8080", nil)
}

// Consulta a cotação do dolar
func ConsultaCotacao(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("teste"))

}

// Faz gravação da cotação no banco de dados
func GravaCotacao(db *gorm.DB) {

}
