package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type ResponseWrapper struct {
	http.ResponseWriter
}

func (r ResponseWrapper) Send(status int, data ...interface{}) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	if len(data) > 0 {
		err := json.NewEncoder(r).Encode(data[0])
		if err != nil {
			panic(err)
		}
	}
}

func (r ResponseWrapper) Error(message interface{}, args ...interface{}) {
	status := 500
	errObj := map[string]interface{}{
		"message": fmt.Sprintf("%s", message),
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

	if status == 500 {
		log.Println(message)
		debug.PrintStack()
	}
	r.Send(status, map[string]interface{}{"error": errObj})
}

func (r ResponseWrapper) RawError(message interface{}, args ...interface{}) {
	status := 500
	argc := len(args)
	if argc > 0 {
		status = args[0].(int)
	}
	if status == 500 {
		debug.PrintStack()
	}
	r.ResponseWriter.WriteHeader(status)
	r.ResponseWriter.Write([]byte(fmt.Sprintf("%s", message)))
}

func Response(w http.ResponseWriter) ResponseWrapper {
	return ResponseWrapper{w}
}
