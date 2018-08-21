package cherry

import (
	"net/http"
	"encoding/json"
)

func OutputJSON(w http.ResponseWriter, v interface{})  {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(v)
}

func OutputError(w http.ResponseWriter, code int, msg string)  {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)

	v := struct {
		Code int `json:"code"`
		Message string `json:"message,omitempty"`
	}{
		Code:code,
		Message:msg,
	}

	json.NewEncoder(w).Encode(v)
}
