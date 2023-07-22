package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	quotation, err := getQuotation()
	if err != nil {
		fmt.Println("Error when try to get a quotation: ", err.Error())
	}
	err = saveOnFile(quotation)
	if err != nil {
		fmt.Println("Error when saving quotation on file: ", err.Error())
	}
}

func saveOnFile(quotation *Quotation) error {
	file, err := os.Create("quotation.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %v", quotation.Bid))
	if err != nil {
		return err
	}
	return nil
}

func getQuotation() (*Quotation, error) {
	client := http.Client{
		Timeout: 300 * time.Millisecond,
	}
	resp, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var quotation *Quotation

	err = json.Unmarshal(res, &quotation)

	if err != nil {
		return nil, err
	}

	return quotation, nil

}

type Quotation struct {
	Bid string `json:"bid"`
}
