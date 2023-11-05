package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

func GetCotationDollar() (*Cotation, error) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, URI_COTATION, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c Cotation
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	bid := c.USDBRL.Bid

	saveBid(bid)
	return &c, nil
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

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS cotations (id INTEGER PRIMARY KEY, bid VARCHAR(30))")
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table cotation!")
	}
	statement.Exec()

	statement, _ = db.Prepare("INSERT INTO cotations (bid) values (?)")
	statement.ExecContext(ctx)
	log.Println("Inserted the cotation into database!")

}
