package session

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

// BuildRoutes returns the routes for the application.
func BuildRoutes() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/login-success", LoginSuccessHandler)
	router.HandleFunc("/verify", VerifyHandler)
	router.HandleFunc("/logout", LogoutHandler)

	// profile routes with LoginRequiredMiddleware
	profileRouter := mux.NewRouter()
	profileRouter.HandleFunc("/profile", ProfileHandler)

	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(LoginRequiredMiddleware),
		negroni.Wrap(profileRouter),
	))

	apiRouter := router.PathPrefix("/v1").Subrouter()
	apiRouter.HandleFunc("/user/signup", AddUser).Methods("POST")
	apiRouter.HandleFunc("/user/signin", GetToken).Methods("POST")
	apiRouter.HandleFunc("/event/new", AddEvent).Methods("POST")
	apiRouter.HandleFunc("/event/delete", DeleteEvent).Methods("POST")
	apiRouter.HandleFunc("/event/update", UpdateEvent).Methods("POST")
	apiRouter.HandleFunc("/event/invite", AddAttendance).Methods("POST")
	apiRouter.HandleFunc("/event/disinvite", DeleteAttendance).Methods("POST")

	apiRouter.HandleFunc("/event", GetEvent).Methods("GET")
	apiRouter.HandleFunc("/user", GetUser).Methods("GET")
	apiRouter.HandleFunc("/event/all", GetEvents).Methods("GET")
	apiRouter.HandleFunc("/user/all", GetUsers).Methods("GET")

	// apply the base middleware to the main router
	n := negroni.New(
		negroni.NewRecovery(), negroni.NewLogger(),
		negroni.HandlerFunc(CsrfMiddleware),
		negroni.HandlerFunc(UserMiddleware),
	)

	n.UseHandler(router)

	return n
}
