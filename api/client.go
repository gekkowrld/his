package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Client struct {
	DB          *sql.DB
	InsertSQL   string
	AddPrograms string
	Search      string
	Profile     string
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
	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
		return
	}
	defer r.Body.Close()

	res, client, err := CreateClient(json_data, c.DB, c.InsertSQL)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
		return
	}
	PrintDBResult(res)
	w.Write([]byte(fmt.Sprintf("%s %s", affectedRows(res.RowsAffected), client.Id)))
}

func CreateClient(jsonData []byte, db *sql.DB, insertSQL string) (sql.Result, ClientInfo, error) {
	var client ClientInfo
	err := json.Unmarshal(jsonData, &client)
	if err != nil {
		return nil, client, err
	}

	uuid7, err := uuid.NewV7()
	if err != nil {
		return nil, client, err
	}
	client.Id = uuid7.String()

	res, err := db.Exec(insertSQL,
		client.Id, client.FirstName, client.MiddleNames, client.LastName, client.Day, client.Month, client.Year)

	if err != nil {
		return nil, client, err
	}

	return res, client, err
}

type clps struct {
	Programs []string `json:"programs"`
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
	json_data, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
		return
	}
	defer r.Body.Close()

	afr, clp, err := ClientProgram(client_id, json_data, c.DB, c.AddPrograms)
	rows_affected := fmt.Sprintf("%s client=%s programs=%v", affectedRows(afr), client_id, clp.Programs)
	log.Println(rows_affected)
	w.Write([]byte(rows_affected))
}

// Associate a client and multiple programs
// The id is the user id and the JSON data expects an array of program ids
// An example:
//
//	{"programs": [01966c3f-306f-75bb-93c5-d602d78f5c49, 01966c39-5f7a-727e-b62f-f557af6453ea]}
//
// This will associate with two programs.
func ClientProgram(id string, jsonData []byte, db *sql.DB, addPrograms string) (func() (int64, error), clps, error) {
	var clp clps

	err := json.Unmarshal(jsonData, &clp)
	if err != nil {
		return nil, clp, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, clp, err
	}

	stmt, err := tx.Prepare(addPrograms)
	if err != nil {
		tx.Rollback()
		return nil, clp, err
	}
	defer stmt.Close()

	var affected int64
	for _, pid := range clp.Programs {
		res, err := stmt.Exec(id, pid)
		if err != nil {
			tx.Rollback()
			return nil, clp, err
		}
		af, _ := res.RowsAffected()
		affected += af
	}

	if err := tx.Commit(); err != nil {
		return nil, clp, err
	}

	return func() (int64, error) { return affected, nil }, clp, nil
}

type search_value struct {
	Id        string
	FirstName string
	LastName  string
}

// search the client table for matches.
// user searches are passed directly to the db AS IS
// NOTE: This poses a security risk as it is vulnerable to SQL injeection
func (c *Client) SearchClient(w http.ResponseWriter, r *http.Request) {
	search_term := r.URL.Query().Get("q")
	search_values, err := SearchClient(search_term, c.DB, c.Search)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
		return
	}

	str := fmt.Sprintf("Results for: %s\n%v\n", search_term, search_values)
	log.Println(str)
	w.Write([]byte(str))
}

func SearchClient(query string, db *sql.DB, search string) ([]search_value, error) {
	rows, err := db.Query(search, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var search_values []search_value
	for rows.Next() {
		s := &search_value{}
		if err := rows.Scan(&s.Id, &s.FirstName, &s.LastName); err != nil {
			return search_values, err
		}
		search_values = append(search_values, *s)
	}
	if err := rows.Err(); err != nil {
		return search_values, err
	}

	return search_values, nil
}

type clientProfile struct {
	Id        string
	FirstName string
	LastName  string
}

func (c *Client) ClientProfile(w http.ResponseWriter, r *http.Request) {
	client_id := r.PathValue("id")
	info, err := ClientProfile(client_id, c.DB, c.Profile)
	if err != nil {
		WriteError(w, EJ{Message: err, Code: http.StatusInternalServerError})
		return
	}
	w.Write([]byte(fmt.Sprintf("%#v\n", info)))
}

// Get user profile and return it.
func ClientProfile(id string, db *sql.DB, search string) (clientProfile, error) {
	var profile clientProfile
	row := db.QueryRow(search, id)
	row.Scan(&profile.Id, &profile.FirstName, &profile.LastName)

	return profile, nil
}
