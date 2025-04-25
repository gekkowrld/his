package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Client struct {
	DB        *sql.DB
	InsertSQL string
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
	w.Write([]byte(fmt.Sprintf("%s %s", affectedRows(res), client.Id)))
}
