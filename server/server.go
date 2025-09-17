package server

import (
	"encoding/json"
	"fmt"
	"goexpert-client-server-api/client"
	"io"
	"net/http"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Dollar struct {
	gorm.Model
	Code      string `json:"code"`
	CodeIn    string `json:"codeIn"`
	Name      string `json:"name"`
	High      string `json:"high"`
	Low       string `json:"low"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Timestamp string `json:"timestamp"`
}

type DollarBR struct {
	USDBRL Dollar `json:"USDBRL"`
}

const urlDolar string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

// Consulta a cotação do dolar
func ConsultaCotacaoSiteEconomia(w http.ResponseWriter, r *http.Request) {

	request, err := http.Get(urlDolar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao fazer requisição: %v\n", err)
	}
	defer request.Body.Close()
	response, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao ler resposta: %v\n", err)
	}
	var data DollarBR
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao fazer o parse da resposta: %v\n", err)
	}
	fmt.Println(data.USDBRL.Bid)

	// Grava os dados no banco de dados
	GravaCotacao(data.USDBRL)

}

// Faz gravação da cotação no banco de dados
func GravaCotacao(data Dollar) {

	// Inicializa o banco de dados
	db, err := NewSqliteDb()
	if err != nil {
		panic(err)
	}

	// Grava os dados no banco
	err = db.Create(&data).Error
	if err != nil {
		panic("erro ao inserir dados: " + err.Error())
	}

	var currency []Dollar

	db.Find(&currency)
	for _, dollar := range currency {
		fmt.Println(dollar.Bid)
	}

	client.LeDolarBancoDeDados()

	fmt.Println("fim do processamento")
}

func NewSqliteDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("server/sqlite.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Dollar{})
	return db, nil
}
