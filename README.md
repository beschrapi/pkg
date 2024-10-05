# Beschrapi

Beschrapi (bəˈʃraːpi), derived from the German word "beschreiben" (meaning "to describe") and the English acronym "API" (Application Programming Interface), is a tool designed for building APIs using YAML files.
It offers a simple and efficient way to create APIs with minimal code.

For more complex API requirements, Beschrapi allows you to include a code file to manage the request and response, giving you full flexibility to customize and implement your API logic as needed.

## Usage

This package is used to provide the APIs from Beschrapi. Its used create code files to include in Beschrapi.

Here is an example of how to use this package:

```go
package main

import (
	"net/http"
	"github.com/beschrapi/pkg" // Here we import the package
)

func Run(w http.ResponseWriter, r *http.Request, app *pkg.App) {
	// Simulate a SQL query with UNION to fake the data
	query := `SELECT 'Hello' AS message`
	result, err := app.Query(query, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	app.Json(w, result)
}
```

This will create an Endpoint that returns a JSON response with the message "Hello".

