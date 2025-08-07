package client

import (
	"io"
	"net/http"
)

func LeDolar() {
	requisicao, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	response, err := io.ReadAll(requisicao.Body)
	if err != nil {
		panic(err)
	}
	println(response)

}
