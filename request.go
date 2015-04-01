package util

import (
	"encoding/json"
	"net/http"
)

type RequestWrapper struct {
	*http.Request
}

func (r RequestWrapper) DecodeBody(to interface{}) error {
	return json.NewDecoder(r.Body).Decode(to)
}

func Request(r *http.Request) RequestWrapper {
	return RequestWrapper{r}
}
