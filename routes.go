package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
        session "github.com/oceandiver/letsgo/lib/session"
        api "github.com/oceandiver/letsgo/lib/api"
)

// BuildRoutes returns the routes for the application.
func BuildRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", session.HomeHandler)
	router.HandleFunc("/login-success", session.LoginSuccessHandler)
	router.HandleFunc("/verify", session.VerifyHandler)
	router.HandleFunc("/logout", session.LogoutHandler)

        router.HandleFunc("/adduser", api.AddUser)
        router.HandleFunc("/addevent", api.AddEvent)
        router.HandleFunc("/event", api.GetEvent)

	// profile routes with LoginRequiredMiddleware
	profileRouter := mux.NewRouter()
	profileRouter.HandleFunc("/profile", session.ProfileHandler)

	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(session.LoginRequiredMiddleware),
		negroni.Wrap(profileRouter),
	))

	// apply the base middleware to the main router
	n := negroni.New(
		negroni.HandlerFunc(session.CsrfMiddleware),
		negroni.HandlerFunc(session.UserMiddleware),
	)
	n.UseHandler(router)

	return n
}
