package server

import (
	"encoding/json"
	"fmt"
	"goexpert-client-server-api/infrastructure"
	"goexpert-client-server-api/types"
	"io"
	"net/http"
	"os"
)

const urlDolar string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func InitServer() {
	http.HandleFunc("/cotacao", ConsultaCotacao)
	http.ListenAndServe(":8080", nil)
}

// Consulta a cotação do dolar
func ConsultaCotacao(w http.ResponseWriter, r *http.Request) {

	request, err := http.Get(urlDolar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao fazer requisição: %v\n", err)
	}
	defer request.Body.Close()
	response, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao ler resposta: %v\n", err)
	}
	var data types.DollarBR
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao fazer o parse da resposta: %v\n", err)
	}
	fmt.Println(data.USDBRL.Bid)

	// Grava os dados no banco de dados
	GravaCotacao(data.USDBRL)

}

// Faz gravação da cotação no banco de dados
func GravaCotacao(data types.Dollar) {

	// Inicializa o banco de dados
	db, err := infrastructure.NewSqliteDb()
	if err != nil {
		panic(err)
	}

	// Grava os dados no banco
	err = db.Create(&data).Error
	if err != nil {
		panic("erro ao inserir dados: " + err.Error())
	}

	var currency []types.Dollar

	db.Find(&currency)
	for _, dollar := range currency {
		fmt.Println(dollar.Bid)
	}

	fmt.Println("fim do processamento")
}
