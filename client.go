package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotation struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

const URI_COTATION = "http://localhost:8080/cotacao"

func main() {
	log.Println("::.. Inicializando busca cotação Dollar ..::")
	defer log.Println("::.. Finalizando Requisição ..::")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", URI_COTATION, nil)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}

	var cotation Cotation
	err = json.Unmarshal(body, &cotation)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
	createFile(&cotation)

}

func createFile(cotation *Cotation) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotation.USDBRL.Bid))
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
}
