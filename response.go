package owl

import (
	"encoding/json"
	"net/http"
)

// JSON sends a JSON response with the given status code.
func JSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

// Text sends a plain text response with the given status code.
func Text(w http.ResponseWriter, code int, text string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	_, err := w.Write([]byte(text))
	return err
}
