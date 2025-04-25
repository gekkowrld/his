package api

import (
	"database/sql"
	"fmt"
	"log"
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
