package main

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/codeui/chevent-web/api/mware"
	"github.com/codeui/chevent-web/api/route"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	appChain = alice.New(
		mware.Logging,
		mware.Cors,
		gziphandler.GzipHandler,
		context.ClearHandler,
		mware.Mongo,
	)
)

func handler() http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/events").Handler(appChain.ThenFunc(route.CreateEventHandler))
	r.Methods("GET").Path("/events").Handler(appChain.ThenFunc(route.ListEventsHandler))

	r.Methods("GET").Path("/version").Handler(appChain.ThenFunc(VersionHandler))

	// a dummy handler to log all the other requests that directs to non existent endpoints
	r.PathPrefix("/").Handler(appChain.Then(http.DefaultServeMux))

	return r
}
