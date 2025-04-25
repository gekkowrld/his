package api

import (
	"database/sql"
	"fmt"
	"log"
)

func PrintDBResult(result sql.Result) {
	log.Println(affectedRows(result))
}

func affectedRows(result sql.Result) string {
	affected, _ := result.RowsAffected()
	return fmt.Sprintf("%d rows affected", affected)
}
