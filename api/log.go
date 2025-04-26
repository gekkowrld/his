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
	writeHttp(w, write)
}

func WriteSucessClient(w http.ResponseWriter, client ClientInfo) {
	data, _ := json.Marshal(client)
	log.Printf("\n%v\n", data)
	writeHttp(w, data)
}

func writeHttp(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func WriteClientSearchResults(w http.ResponseWriter, search []search_value) {
	data, err := json.Marshal(search)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
	}
	writeHttp(w, data)
}

func WriteClientProfile(w http.ResponseWriter, profile clientProfile) {
	data, err := json.Marshal(profile)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
	}
	log.Printf("%s\n", profile)
	writeHttp(w, data)
}

func WriteSucessClientProgram(w http.ResponseWriter, clientId string, programs []string ) {
	data, err := json.Marshal(struct{
		ClientId string `json:"client_id"`
		Programs []string `json:"programs"`
	}{
		ClientId: clientId,
		Programs: programs,
	})
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
	}
	writeHttp(w,data)
}
