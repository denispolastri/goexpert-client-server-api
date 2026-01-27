package main

import (
	"context"
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

	start := time.Now()

	// Cria um contexto com timeout de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Cria a requisição com o contexto
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		slog.Error("erro ao criar a requisição", "error", err)
		return
	}

	// Executa a requisição
	client := &http.Client{}
	request, err := client.Do(req)
	if err != nil {
		slog.Error("erro ao fazer a requisição", "error", err)
		// Verifica se o erro foi por timeout
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("timeout: requisição excedeu 300ms")
		}
		return
	}
	defer request.Body.Close()

	response, err := io.ReadAll(request.Body)
	if err != nil {
		slog.Error("erro ao ler o body da resposta", "error", err)
		return
	}

	duration := time.Since(start)
	slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())

	var bid string
	err = json.Unmarshal(response, &bid)
	if err != nil {
		slog.Error("erro ao fazer o parse da resposta", "error", err)
	} else {
		slog.Info("cotação do dólar lida com sucesso", "bid", bid)
	}

	// Cria o arquivo cotacao.txt
	file, err := os.Create("cotacao.txt")
	if err != nil {
		slog.Error("erro ao criar o arquivo", "error", err)
	}
	defer file.Close()

	if err == nil {
		_, err = file.WriteString("Dólar:{" + bid + "}")
		if err != nil {
			slog.Error("erro ao escrever no arquivo", "error", err)
		}
	}

}
