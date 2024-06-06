package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Function to handle requests.
func (app *Application) routes() http.Handler {
	middlewareChain := alice.New(app.recoverPanic, app.secureHeaders, app.RequestLogger, app.ResponseLogger)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf)
	mux := http.NewServeMux()
	mux.Handle("/", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Handle("/home", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.home)) // index route (homepage).
	mux.Handle("/addTasks", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.addTasks))  // add tasks
	mux.Handle("/update", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updateTasks)) // update tasks
	mux.Handle("/delete", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.delTasks))    // delete tasks

	// Other routes
	mux.Handle("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	//mux.Handle("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Handle("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	//mux.Handle("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Handle("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir(app.config.StaticDir))   // serve static files
	mux.Handle("/static/", http.StripPrefix("/static", fileServer)) // strip static directory.

	return middlewareChain.Then(mux)
}



