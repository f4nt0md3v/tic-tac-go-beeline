package httpx

import (
	"encoding/json"
	"net/http"
)

const (
	ContentTypeHeader   = "Content-Type"
	ContentTypeBinary   = "application/octet-stream"
	ContentTypeForm     = "application/x-www-form-urlencoded"
	ContentTypeJSON     = "application/json"
	ContentTypeJsonUtf8 = "application/json;charset=utf-8"
	ContentTypeHTML     = "text/html;charset=UTF-8"
	ContentTypeText     = "text/plain;charset=UTF-8"
)

func ParseJSON(r *http.Request, item interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	err := decoder.Decode(item)
	if err != nil {
		return err
	}

	return nil
}

func WriteJson(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", ContentTypeJsonUtf8)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}
