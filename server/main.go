package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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

func main() {

	// Define valor HTTP_PORT
	httpPort := "8080"
	//
	http.HandleFunc("/cotacao", ConsultaCotacaoSiteEconomia)
	fmt.Println("Server HTTP UP port " + httpPort)
	http.ListenAndServe(":"+httpPort, nil)

}

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
	// Grava os dados no banco de dados
	GravaCotacao(data.USDBRL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data.USDBRL.Bid)

}

// Faz gravação da cotação no banco de dados
func GravaCotacao(data Dollar) {

	// Inicializa o banco de dados
	db, err := NewSqliteDb()
	if err != nil {
		panic(err)
	}

	var currencys []Dollar

	// corrigir este Find e incluir um where com o camnpo bid, se já estiver cadastrado a cotação com este valor, apenas altera.
	//db.Find(&currencys)
	db.Where("bid = ?", data.Bid).Find(&currencys)

	fmt.Println("existem " + strconv.Itoa(len(currencys)) + " moedas cadastradas no banco")

	if len(currencys) == 0 {
		// Grava a moeda se o banco estiver vazio
		err = db.Create(&data).Error
		if err != nil {
			panic("erro ao inserir dados: " + err.Error())
		}
	} else {
		count := 0
		for _, cur := range currencys {
			if count == 0 {
				// Altera no banco os demais dados da moeda de valor igual
				err = db.Updates(cur).Error
				if err != nil {
					panic("erro ao alterar dados: " + err.Error())
				}
			} else {
				// Grava no banco a moeda que ainda não existe
				err = db.Delete(&currencys[count]).Error
				if err != nil {
					panic("erro ao inserir dados: " + err.Error())
				}
			}
			count++
		}
	}

}

func NewSqliteDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Dollar{})
	return db, nil
}
