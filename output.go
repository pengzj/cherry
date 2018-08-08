package cherry

import (
	"net/http"
	"encoding/json"
)

func OutputJSON(w http.ResponseWriter, v interface{})  {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
