package pkg

import (
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	// "github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for sql.DB
)

type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
	// TBD
)

func NewPool(
	databaseType DatabaseType,
	user string,
	password string,
	host string,
	port int,
	database string,
) (*sql.DB, error) {
	connectionUrl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, database)
	db, err := sql.Open("pgx", connectionUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to verify connection %v", err)
	}

	return db, nil
}

func prepareQuery(query string, params *map[string]interface{}) (string, *[]any) {
	if params == nil {
		return query, nil
	}

	// Prepare the query with the params
	// The params must be returned in correct order
	// SELECT * FROM table WHERE column1 = {param1} AND column2 = {param2}
	// The params must be in the order of [param1(value), param2(value)]
	// Parameters get replaced with the desired placeholder value
	var queryParams []interface{}
	re := regexp.MustCompile(`\{([^}]+)\}`)
	query = re.ReplaceAllStringFunc(query, func(match string) string {
		key := match[1 : len(match)-1] // Extract key without braces {name} -> name
		queryParams = append(queryParams, (*params)[key])
		return fmt.Sprintf(" $%d ", len(queryParams))
	})

	query = strings.TrimSpace(query)
	return query, &queryParams
}

func ExecuteQuery(db *sql.DB, query string, params *map[string]interface{}) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	query, queryParams := prepareQuery(query, params)
	slog.Debug("Prepared query", slog.Any("query", query))

	if queryParams != nil {
		rows, err = db.Query(query, (*queryParams)...)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		columnValues := make([]interface{}, len(columns))
		for i := range columnValues {
			var col interface{}
			columnValues[i] = &col
		}

		if err := rows.Scan(columnValues...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := *(columnValues[i].(*interface{}))
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}

		results = append(results, rowMap)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %v", rows.Err())
	}

	return results, nil
}
