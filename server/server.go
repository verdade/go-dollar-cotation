package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"time"
)

const URI_COTATION = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetCotationDollar)
	http.ListenAndServe(":8080", mux)
}

func GetCotationDollar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, URI_COTATION, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}

	var c Cotation
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}

	bid := c.USDBRL.Bid
	saveBid(bid)
	w.Write([]byte(string(string(body))))
}

func saveBid(bid string) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	db, err := sql.Open("sqlite3", "cotation.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	smt, err := db.Prepare("INSERT INTO cotations (bid) values (?)")
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
	_, err = smt.ExecContext(ctx, bid)
	if err != nil {
		log.Println(">>> Deu erro %v\n <<<", err)
	}
	log.Println("::.. Inserted the cotation into database! ..::")
}
