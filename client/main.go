package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

func main() {

	LeDolarBancoDeDados()

}

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

func LeDolarBancoDeDados() {

	request, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()
	response, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var data DollarBR
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao fazer o parse da resposta: %v\n", err)
	}
	fmt.Println(data.USDBRL)

}
