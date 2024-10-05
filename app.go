package pkg

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type App struct {
	// Intented for complex operations instead use `Query` method
	Db     *sql.DB
	Server *Server
}

func NewApp(db *sql.DB, server *Server) *App {
	return &App{Db: db, Server: server}
}

func (r *App) Query(query string, params *map[string]interface{}) ([]map[string]interface{}, error) {
	return ExecuteQuery(r.Db, query, params)
}

func (r *App) Json(w http.ResponseWriter, data any) error {
	// try to marshal data and send it as response
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
	return nil
}
