package handler

import "net/http"

type IHandler interface {
	http.Handler
	Route()
	Close()
}
