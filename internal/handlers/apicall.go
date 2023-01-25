package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type TablesRequest struct {
	TablesName []string
	Uuid       string
}

type TablesResponse struct {
	Success   bool           `json:"success"`
	DateSince string         `json:"date_since"`
	TablesIds map[string]int `json:"tables_ids"`
}

type TableRequest struct {
	Uuid string `json:"uuid"`
	Data []byte `json:"data"`
}

type TableResponse struct {
	Data   map[string]int `json:"data"`
	Status int         `json:"status"`
}

func LastIds(host string, data TablesRequest) (*TablesResponse, error) {
	res, err := http.Get(host + "api/get_last_ids/")
	if err != nil {
		return nil, err
	}

	var tablesData *TablesResponse
	err = json.NewDecoder(res.Body).Decode(&tablesData)
	if err != nil {
		return nil, err
	} else {
		log.Println("Got last ids")
	}

	return tablesData, nil
}

func SendTableData(host string, data TableRequest) (*TableResponse, error) {
	postBody, _ := json.Marshal(data)
	responseBody := bytes.NewBuffer(postBody)

	res, err := http.Post(host+"api/accept_data/", "application/json", responseBody)
	if err != nil {
		return nil, err
	}

	var tableData *TableResponse
	err = json.NewDecoder(res.Body).Decode(&tableData)
	if err != nil {
		return nil, err
	}

	return tableData, nil
}
