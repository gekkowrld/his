package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func PrintDBResult(result sql.Result) {
	log.Println(affectedRows(result.RowsAffected))
}

func affectedRows(result func() (int64, error)) string {
	affected, _ := result()
	var row = "rows"
	if affected == 1 {
		row = "row"
	}
	return fmt.Sprintf("%d %s affected", affected, row)
}

type EJ struct {
	Message error
	Code    int
}

// Return JSON error instead of text
func ErrorJSON(err EJ) []byte {
	erc := err.Code
	if erc <= 0 {
		erc = 404
	}
	data, _ := json.Marshal(struct {
		Message string `json:"name"`
		Code    int    `json:"code"`
	}{
		Message: err.Message.Error(),
		Code:    erc,
	})
	return data
}

func WriteError(w http.ResponseWriter, err EJ) {
	write := ErrorJSON(err)
	log.Printf("\nCode: %d\nMessage: %s\n", err.Code, err.Message)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)
	w.Write(write)
}
