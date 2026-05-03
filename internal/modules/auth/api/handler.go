package api

import "github.com/go-playground/validator/v10"

type Serviceauth interface {
	authService
}

type Handlerauth struct {
	svcAuth   Serviceauth
	validator *validator.Validate
	debug     bool
}

func NewHandlerauth(svcAuth Serviceauth, debug bool) *Handlerauth {
	return &Handlerauth{
		svcAuth:   svcAuth,
		validator: validator.New(),
		debug:     debug,
	}
}
