package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/verdade/go-dollar-cotation/server"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotation, err := server.GetCotationDollar()
	if err != nil {
		panic(err)
	}

	res, err := json.Marshal(cotation)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprint(os.Stderr, "Deu erro %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotation.USDBRL.Bid))
	if err != nil {
		fmt.Fprint(os.Stderr, "Deu erro %v\n", err)
	}
	w.Write([]byte(string(res)))
}
