package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		slog.Error("erro ao conectar no banco de dados", "error", err.Error())
		panic(err)
	}
	defer db.Close()

	// Define valor HTTP_PORT
	httpPort := "8080"

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ConsultaCotacaoSiteEconomia(w, r, db)
	})
	slog.Info("Server HTTP UP port " + httpPort)
	http.ListenAndServe(":"+httpPort, nil)

}

// Consulta a cotação do dolar
func ConsultaCotacaoSiteEconomia(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var urlDolar string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

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
	err = GravaCotacao(data.USDBRL, db)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		slog.Error("erro ao gravar cotação no banco de dados", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data.USDBRL.Bid)
	slog.Info("-------------------------------")

}

// Faz gravação da cotação no banco de dados
func GravaCotacao(data Dollar, db *sql.DB) error {

	var timeoutDuration = 10 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	start := time.Now()

	// Verifica se já expirou antes de iniciar
	if ctx.Err() != nil {
		slog.Error("timeout antes de iniciar transação", "error", ctx.Err().Error())
		duration := time.Since(start)
		slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())
		return ctx.Err()
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("erro ao iniciar transação", "error", err.Error())
		duration := time.Since(start)
		slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())
		return err
	}

	// Verifica timeout após iniciar transação
	select {
	case <-ctx.Done():
		tx.Rollback()
		slog.Error("timeout ao iniciar transação", "error", ctx.Err().Error())
		duration := time.Since(start)
		slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())
		return ctx.Err()
	default:
	}

	query := `INSERT INTO dollars (code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`
	_, err = tx.ExecContext(ctx, query, data.Code, data.CodeIn, data.Name, data.High, data.Low, data.VarBid, data.PctChange, data.Bid, data.Ask, data.Timestamp)

	time.Sleep(9 * time.Millisecond) // Simula demora na execução

	// Verifica timeout imediatamente após execução
	if ctx.Err() != nil {
		tx.Rollback()
		slog.Error("tempo de execução excedido", "error", ctx.Err().Error())
		duration := time.Since(start)
		slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())
		return ctx.Err()
	}

	if err != nil {
		tx.Rollback()
		slog.Error("erro ao inserir dados", "error", err.Error())
		duration := time.Since(start)
		slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())
		return err
	}

	if err := tx.Commit(); err != nil {
		slog.Error("erro ao commitar transação", "error", err.Error())
		duration := time.Since(start)
		slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())
		return err
	}

	duration := time.Since(start)
	slog.Info("requisição finalizada", "duration_ms", duration.Milliseconds())

	return nil
}
