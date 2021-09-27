package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/login", dynamicMiddleware.ThenFunc(app.login))
	mux.Get("/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.home))
	mux.Get("/all", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getAll))
	mux.Post("/getByMfa", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getByMfa))
	mux.Get("/personne/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPersonneForm))
	mux.Post("/personne/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updatePersonneForm))
	mux.Post("/personne/json/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPersonneObject))
	mux.Post("/personne/json/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updatePersonneObject))
	mux.Post("/personne/json/updateFlagMail", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateFlag))
	mux.Post("/delete/personne", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deletePersonneById))
	mux.Post("/personne/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPersonne))
	mux.Post("/formation/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.formationCreate))
	mux.Get("/personne/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showPersonne))
	mux.Get("/mailing/:action", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.mailProcessing))
	mux.Post("/saveOrCheckUser", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.SaveOrCheckLogin))
	mux.Get("/documents", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getDocuments))
	mux.Get("/formations", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getFormations))
	mux.Get("/formation/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getFormationById))
	mux.Get("/sendedMails", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.sendedMails))
	mux.Post("/sendMail", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.foreSendMail))

	// Add the five new routes.
	//mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	//mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	var fileServer http.Handler
	if app.env == "DEV" {
		fileServer = http.FileServer(http.Dir("./ui/static/"))
	} else {
		fileServer = http.FileServer(http.Dir("/var/www/go/deploy/agefice/ui/static/"))
	}

	//mux.Get("/static/", http.StripPrefix("/static", fileServer))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
