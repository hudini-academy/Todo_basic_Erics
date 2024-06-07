package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/justinas/nosurf"
)

// Handles server error responses.
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // stores the error message and stacktrace.
	app.errorlog.Output(2, trace) // To get the exact file and line number.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Handles client error responses.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Handles client error responses
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Renders the template for the application.
func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		log.Println(ok)
		return
	}

	buffer := new(bytes.Buffer)

	err := ts.Execute(buffer, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buffer.WriteTo(w)
}

// Authenticate the user by checking whether the userId is in the session.
func (app *Application) authenticatedUser(r *http.Request) int {
	log.Println(app.session.GetInt(r, "userID"))
	return app.session.GetInt(r, "userID")
}
// Loads the default data required when the template is rendered.
func (app *Application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.Flash = app.session.PopString(r, "flash")
	return td
}
