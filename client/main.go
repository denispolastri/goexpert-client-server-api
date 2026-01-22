package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

func main() {

	// Logger default
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	var conteudo string
	var lErroCotacao bool = false

	start := time.Now()

	request, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		slog.Error("erro ao fazer a requisição", "error", err)
	}
	defer request.Body.Close()
	response, err := io.ReadAll(request.Body)
	if err != nil {
		slog.Error("erro ao ler o body da resposta", "error", err)
	}

	duration := time.Since(start)
	slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())

	var bid string
	err = json.Unmarshal(response, &bid)
	if err != nil {
		bid = "erro ao ler cotação"
		slog.Error("erro ao fazer o parse da resposta", "error", err)
	} else {
		lErroCotacao = true
		slog.Info("cotação do dólar lida com sucesso", "bid", bid)
	}

	// Cria o arquivo cotacao.txt
	file, err := os.Create("cotacao.txt")
	if err != nil {
		slog.Error("erro ao criar o arquivo", "error", err)
	}
	defer file.Close()

	if lErroCotacao {
		conteudo = "Dólar:{" + bid + "}"
	} else {
		conteudo = "Dólar:{erro ao ler cotação}"
	}
	_, err = file.WriteString(conteudo)
	if err != nil {
		slog.Error("erro ao escrever no arquivo", "error", err)
	}
}
