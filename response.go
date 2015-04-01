package util

import (
	"encoding/json"
	"net/http"
)

type ResponseWrapper struct {
	http.ResponseWriter
}

func (r ResponseWrapper) Send(status int, data ...interface{}) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	if len(data) > 0 {
		json.NewEncoder(r).Encode(data[0])
	}
}

func (r ResponseWrapper) Error(message string, args ...interface{}) {
	status := 500
	errObj := map[string]interface{}{
		"message": message,
	}
	argc := len(args)
	if argc > 0 {
		status = args[0].(int)
	}
	if argc > 1 {
		for k, v := range args[1].(map[string]interface{}) {
			errObj[k] = v
		}
	}
	if argc > 2 {
		panic("invalid arguments")
	}

	r.Send(status, map[string]interface{}{"error": errObj})
}

func Response(w http.ResponseWriter) ResponseWrapper {
	return ResponseWrapper{w}
}
