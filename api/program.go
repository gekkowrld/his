package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Program struct {
	DB        *sql.DB
	InsertSQL string
}

type ProgramInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   int    `json:"start_date"`
	EndDate     int    `json:"end_date"`
}

// Write the program name into db
// If there is an Error, it prints the error to the writer and stdout and exits.
func (p *Program) CreateProgram(w http.ResponseWriter, r *http.Request) {
	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	res, program, err := CreateProgram(json_data, p.DB, p.InsertSQL)
	if err != nil {
		w.Write([]byte(err.Error()))
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
		return
	}

	PrintDBResult(res)
	WriteSuccessProgram(w, program)
}

// Create a program and write it to the database.
func CreateProgram(jsonData []byte, db *sql.DB, insertSQL string) (sql.Result, ProgramInfo, error) {
	var program ProgramInfo
	err := json.Unmarshal(jsonData, &program)
	if err != nil {
		return nil, program, err
	}

	uuid7, err := uuid.NewV7()
	if err != nil {
		return nil, program, err
	}
	program.Id = uuid7.String()

	res, err := db.Exec(insertSQL,
		program.Id, program.Name, program.Description, program.StartDate, program.EndDate)

	if err != nil {
		return nil, program, err
	}

	return res, program, nil
}
