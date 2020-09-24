package v1

import (
	"net/http"
)

type API struct {
	Path            string
	Method          string
	Func            http.HandlerFunc
}

func NewAPI(path string, method string, handlerFunc http.HandlerFunc) API {
	return API{path, method, handlerFunc}
}
