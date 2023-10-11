package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func loggerHVTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
			var requestData map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&requestData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			loggedData, err := json.Marshal(requestData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			buffer := bytes.NewBuffer(loggedData)

			readCloser := io.NopCloser(buffer)

			r, err := http.NewRequest("POST", r.URL.String(), readCloser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			amount, ok := requestData["amount"].(float64)
			txType, okType := requestData["transaction_type"].(string)
			if !ok || amount <= 10000 && txType == "deposit" || !okType || txType != "deposit" {
				next.ServeHTTP(w, r)
				return
			}

			fmt.Println("request amount > 10000:", string(loggedData))
			next.ServeHTTP(w, r)
			return
		}
	})
}
