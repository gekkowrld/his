package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	var program ProgramInfo
	err := json.NewDecoder(r.Body).Decode(&program)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	uuid7, err := uuid.NewV7()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	program.Id = uuid7.String()

	res, err := p.DB.Exec(p.InsertSQL,
		program.Id, program.Name, program.Description, program.StartDate, program.EndDate)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	PrintDBResult(res)
	w.Write([]byte(fmt.Sprintf("%s %s", affectedRows(res), program.Id)))
}
