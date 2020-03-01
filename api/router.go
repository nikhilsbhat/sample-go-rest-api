package api

import (
	"io"

	"github.com/gorilla/mux"
)

// MuxIn implements the method NewRouter to create a newrouter and holds data for the same.
type MuxIn struct {
	Apilog io.Writer
}

// NewRouter gives back the router which it created to the function/method who called it.
func (log *MuxIn) NewRouter() *mux.Router {

	rout := mux.NewRouter().StrictSlash(true)
	rout.Use(TimeoutHandler)

	//initializing logger with log path
	test := new(LogInit)
	test.Logpath = log.Apilog
	rout.Use(test.Logger)

	for _, route := range Routes {
		rout.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		//rout.Use(mid.JsonHandler)
	}
	return rout
}
