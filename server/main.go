package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

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

func main() {

	// Logger default
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Define valor HTTP_PORT
	httpPort := "8080"

	http.HandleFunc("/cotacao", ConsultaCotacaoSiteEconomia)
	slog.Info("Server HTTP UP port " + httpPort)
	http.ListenAndServe(":"+httpPort, nil)

}

// Consulta a cotação do dolar
func ConsultaCotacaoSiteEconomia(w http.ResponseWriter, r *http.Request) {

	var urlDolar string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	ctx := r.Context()

	request, err := http.Get(urlDolar)
	if err != nil {
		slog.Error("erro ao fazer requisição", "error", err)
	}
	defer request.Body.Close()
	response, err := io.ReadAll(request.Body)
	if err != nil {
		slog.Error("erro ao ler resposta", "error", err)
	}
	var data DollarBR
	err = json.Unmarshal(response, &data)
	if err != nil {
		slog.Error("erro ao fazer o parse da resposta", "error", err)
	}
	// Grava os dados no banco de dados
	err = GravaCotacao(ctx, data.USDBRL)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(data.USDBRL.Bid)
	slog.Info("-------------------------------")

}

// Faz gravação da cotação no banco de dados
func GravaCotacao(ctx context.Context, data Dollar) error {

	// Inicializa o banco de dados
	db, err := NewSqliteDb()
	if err != nil {
		slog.Error("erro ao conectar no banco de dados", "error", err)
		return err
	}
	// Define um timeout para a operação de gravação
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	start := time.Now()

	// Grava a moeda se o banco estiver vazio
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		slog.Error("erro ao inserir dados", "error", err.Error())
	}

	switch {
	case err == nil:
		// Gravação concluída com sucesso
		duration := time.Since(start)
		slog.Info("Dados inseridos com sucesso", "duration_ms", duration.Milliseconds())
		return nil

	case errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled):
		// Operação de gravação cancelada ou expirou o tempo limite
		duration := time.Since(start)
		slog.Info("Tempo limite excedeu", "duration_ms", duration.Milliseconds())
		return err
	}
	return nil
}

func NewSqliteDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		slog.Error("erro ao conectar no banco de dados", "error", err.Error())
		return nil, err
	}
	db.AutoMigrate(&Dollar{})
	return db, nil
}
