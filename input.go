package cherry

import (
	"encoding/json"
	"net/http"
)

func ParseJson(request *http.Request, data interface{}) error  {
	decoder := json.NewDecoder(request.Body)
	return decoder.Decode(data)
}
