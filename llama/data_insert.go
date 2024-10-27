package llama

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    // "log"
)

func fetchData(db *sql.DB) ([]string, error) {
    rows, err := db.Query("SELECT column_name FROM table_name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var data []string
    for rows.Next() {
        var item string
        if err := rows.Scan(&item); err != nil {
            return nil, err
        }
        data = append(data, item)
    }
    return data, nil
}
