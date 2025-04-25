package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gekkowrld/his/api"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/*.sql
var sqldb embed.FS

type ress struct {
	db_name string
	schema  string
}

type rest struct {
	db *sql.DB
}

func main() {
	// can be stored in data/ directory.
	// but db/ works too
	res := ress{db_name: "db/health_sys"}
	rest := rest{}

	db_shema, err := sqldb.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("Error opening schema file for reading: %s", err)
	}
	add_program, err := sqldb.ReadFile("db/add_program.sql")
	if err != nil {
		log.Fatalf("Error opening add_program file for reading: %s", err)
	}
	add_client, err := sqldb.ReadFile("db/add_client.sql")
	if err != nil {
		log.Fatalf("Error opening add_client file for reading: %s", err)
	}
	client_program, err := sqldb.ReadFile("db/client_program.sql")
	if err != nil {
		log.Fatalf("Error opening client_program file for reading: %s", err)
	}
	search_client, err := sqldb.ReadFile("db/client_search.sql")
	if err != nil {
		log.Fatalf("Error opening client_search file for reading: %s", err)
	}

	res.schema = string(db_shema)
	db, err := initial_setup(res)
	if err != nil || db == nil {
		log.Fatalf("Setup error: %v", err)
	}
	rest.db = db
	prog := api.Program{DB: db, InsertSQL: string(add_program)}
	client := api.Client{
		DB:          db,
		InsertSQL:   string(add_client),
		AddPrograms: string(client_program),
		Search:      string(search_client),
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /program", prog.CreateProgram)
	router.HandleFunc("POST /client", client.CreateClient)
	router.HandleFunc("POST /client/{id}", client.ClientProgram)
	router.HandleFunc("GET /search/client", client.SearchClient)

	server := http.Server{Addr: ":2343", Handler: router}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Started server at 2343")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5_000_000)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	teardown(rest)
}

// Setup the db and other resouces before the program runs.
func initial_setup(res ress) (*sql.DB, error) {
	log.Println("Starting Setup...")
	defer log.Println("Finished Setup...")

	db, err := sql.Open("sqlite3", res.db_name)
	if err != nil {
		return nil, err
	}

	db_res, err := db.Exec(res.schema)
	if err != nil {
		return nil, err
	}

	api.PrintDBResult(db_res)
	return db, err
}

// Cleanup the db and other related resouces
func teardown(res rest) {
	log.Println("Started teardown...")
	defer log.Println("Finished teardown...")

	if res.db != nil {
		res.db.Close()
	}
}
