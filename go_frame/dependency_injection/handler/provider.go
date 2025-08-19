package handler

import "github.com/google/wire"

var (
	HdlSet = wire.NewSet(NewGinHandler)
)
