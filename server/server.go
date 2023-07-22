package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server := createServer()
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}

func handler(w http.ResponseWriter, r *http.Request) {
	quotation, err := getDolarQuotation()
	if err != nil {
		fmt.Println("Error when getting a dolar quotation: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = saveQuotation(quotation)
	if err != nil {
		fmt.Println("Error when saving quotation on database: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := createResponse(*quotation)

	json.NewEncoder(w).Encode(&response)
}

func createResponse(quotation Quotation) *Response {
	return &Response{Bid: quotation.Bid}
}

func getDolarQuotation() (*Quotation, error) {
	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}
	resp, err := client.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response QuotationReponse

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	return response.USDBRL, nil

}

func createServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handler)
	return mux
}

func saveQuotation(quotation *Quotation) error {
	db, err := sql.Open("sqlite3", "challenge.db")
	if err != nil {
		fmt.Println("Error when opening database connection: ", err.Error())
	}
	defer db.Close()

	db.Exec(`CREATE TABLE IF NOT EXISTS quotations (
		id INTEGER PRIMAY KEY, 
		code varchar(50),
		codein varchar(50),
		name varchar(100),
		high varchar(30),
		low varchar(30),
		varBid varchar(30),
		pctChange varchar(30),
		bid varchar(30),
		ask varchar(30),
		timestamp varchar(30),
		createDate varchar(50)
		 )`,
	)

	stmt, err := db.Prepare(`INSERT INTO quotations (Code, Codein, Name, High, Low, VarBid, PctChange, Bid, Ask, Timestamp, CreateDate)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = stmt.ExecContext(
		ctx,
		quotation.Code,
		quotation.Codein,
		quotation.Name,
		quotation.High,
		quotation.Low,
		quotation.VarBid,
		quotation.PctChange,
		quotation.Bid,
		quotation.Ask,
		quotation.Timestamp,
		quotation.CreateDate,
	)
	if err != nil {
		return err
	}

	fmt.Println("Quotation saved with success")

	return nil
}

type QuotationReponse struct {
	USDBRL *Quotation `json:"USDBRL"`
}

type Quotation struct {
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
}

type Response struct {
	Bid string `json:"bid"`
}
