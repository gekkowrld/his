package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Client struct {
	DB          *sql.DB
	InsertSQL   string
	AddPrograms string
}

type ClientInfo struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	MiddleNames string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Day         int    `json:"day"`
	Month       int    `json:"month"`
	Year        int    `json:"year"`
}

// Add client infomation to the database
func (c *Client) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client ClientInfo
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	uuid7, err := uuid.NewV7()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	client.Id = uuid7.String()

	res, err := c.DB.Exec(c.InsertSQL,
		client.Id, client.FirstName, client.MiddleNames, client.LastName, client.Day, client.Month, client.Year)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	PrintDBResult(res)
	w.Write([]byte(fmt.Sprintf("%s %s", affectedRows(res.RowsAffected), client.Id)))
}

// Associate a client and multiple programs
// The {id} is the user id and the POST request expects an array of program ids
// An example:
//
//	{"programs": [01966c3f-306f-75bb-93c5-d602d78f5c49, 01966c39-5f7a-727e-b62f-f557af6453ea]}
//
// This will associate with two programs.
func (c *Client) ClientProgram(w http.ResponseWriter, r *http.Request) {
	client_id := r.PathValue("id")
	var clp struct {
		Programs []string `json:"programs"`
	}
	err := json.NewDecoder(r.Body).Decode(&clp)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tx, err := c.DB.Begin()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	stmt, err := tx.Prepare(c.AddPrograms)
	if err != nil {
		tx.Rollback()
		w.Write([]byte(err.Error()))
		return
	}
	defer stmt.Close()

	var affected int64
	for _, pid := range clp.Programs {
		res, err := stmt.Exec(client_id, pid)
		if err != nil {
			tx.Rollback()
			w.Write([]byte(err.Error()))
			return
		}
		af, _ := res.RowsAffected()
		affected += af
	}

	if err := tx.Commit(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	rows_affected := fmt.Sprintf("%s client=%s programs=%v", affectedRows(func() (int64, error) { return affected, nil }), client_id, clp.Programs)
	log.Println(rows_affected)
	w.Write([]byte(rows_affected))
}
