package api

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"version": "1.0",
	}
	err := indexTemplate.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
