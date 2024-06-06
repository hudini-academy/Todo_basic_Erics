package main

import (
	"log"
	"net/http"
	"strconv"
	"todo/pkg/forms"
	"todo/pkg/models"
)

var errorData *forms.Form

// Home page for tasks application.
func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	task, errGetting := app.todo.GetAll()
	if errGetting != nil {
		app.errorlog.Println(errGetting.Error())
		app.serverError(w, errGetting)
		log.Println(errGetting)
		return
	}

	// Check for the flash message.
	message := app.session.PopString(r, "flash")
	if errorData != nil {
		app.render(w, r, "home.page.tmpl", &templateData{
			Form:     errorData,
			Snippets: task,
		})
		errorData = nil
		return
	}
	// Render the remplate.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: task,
		Flash:    message,
	})
}

// Add tasks for tasks application.
func (app *Application) addTasks(w http.ResponseWriter, r *http.Request) {

	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Data validation.
	form := forms.New(r.PostForm)
	form.Required("task")

	if !form.Valid() {
		errorData = form
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	// Insert into the database.
	_, errInsert := app.todo.Insert(r.FormValue("task"))
	if errInsert != nil {
		app.errorlog.Println(errInsert)
		return
	}

	// Set a success message in the session.
	app.session.Put(r, "flash", "Task successfully created!")
	// Redirect to home page.
	http.Redirect(w, r, "/home", http.StatusFound)
}

// Delete tasks for tasks application.
func (app *Application) delTasks(w http.ResponseWriter, r *http.Request) {
	idToDelete, _ := strconv.Atoi(r.FormValue("id")) // convert task id to integer.

	_, errDel := app.todo.DelTaskDB(idToDelete)
	if errDel != nil {
		app.errorlog.Println(errDel)
		log.Println(errDel)
		return
	}
	// Redirect to home page.
	http.Redirect(w, r, "/home", http.StatusFound)
}

// Update tasks for tasks application.
func (app *Application) updateTasks(w http.ResponseWriter, r *http.Request) {
	idToUpdate, _ := strconv.Atoi(r.FormValue("id")) // Convert the id to integer.

	_, errUpdate := app.todo.UpdateTaskDB(idToUpdate, r.FormValue("updateTitle"))
	if errUpdate != nil {
		app.errorlog.Println(errUpdate)
		log.Println(errUpdate)
		return
	}
	// Redirect to home page.
	http.Redirect(w, r, "/home", http.StatusFound)
}

// signup form for users
func (app *Application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		app.signupUser(w, r)
		return
	}
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// signup users
func (app *Application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 5)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// If there are any errors, re-display the signup form with the errors.
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming tha
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Login user Form
func (app *Application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		app.loginUser(w, r)
		return
	}
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// Login user
func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic message.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	log.Println(err)
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		log.Println(err)
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Logout user
func (app *Application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've be
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/user/login", http.StatusBadRequest)
}
