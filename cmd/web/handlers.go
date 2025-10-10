package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"snippetbox.ariffil.com/internal/models"
	"snippetbox.ariffil.com/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email		string `form:"email"`
	Password	string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) sendFile(w http.ResponseWriter, r *http.Request) {

	app.infoLog.Println("entered the sendFile handler")

	app.infoLog.Printf("getfile: %s\n", r.URL.Path)

	requestedFile := strings.TrimPrefix(r.URL.Path, "/getfile/")

	app.infoLog.Printf("requested file is: %s\n", requestedFile)

	http.ServeFile(w, r, filepath.Clean((requestedFile)))

}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	latestSnippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Snippets = latestSnippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	snippetId, err := strconv.Atoi(params.ByName("id"))

	if err != nil || snippetId < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(snippetId)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return

	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.BelowMaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(int(form.Expires), 1, 7, 365), "expires", "This field should be 1, 7 or 365")

	if !form.Valid() {

		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)

		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessManager.Put(r.Context(), "flash", "Snippet successfully created!")

	newSnippetURL := fmt.Sprintf("/snippet/view/%d", id)
	http.Redirect(w, r, newSnippetURL, http.StatusSeeOther)

}

func (app *application) serveLoginPage(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.tmpl.html", data)

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var loginForm userLoginForm

	err = app.formDecoder.Decode(&loginForm, r.PostForm)

	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	loginForm.CheckField(validator.NotBlank(loginForm.Email), "email", "Email field can't be empty")
	loginForm.CheckField(validator.IsValidEmail(loginForm.Email), "email", "Please provide a valid email")

	loginForm.CheckField(validator.NotBlank(loginForm.Password), "password", "Password field can't be empty")

	if !loginForm.Valid() {

		data := templateData{}
		data.Form = loginForm

		app.render(w, http.StatusUnauthorized, "login.tmpl.html", &data)
		return
	}

	userID, err := app.users.Authenticate(loginForm.Email, loginForm.Password)

	if err != nil {
		loginForm.AddNonFieldError("Email or password is wrong!")

		data := templateData{}
		data.Form = loginForm

		app.render(w, http.StatusUnauthorized, "login.tmpl.html", &data)
		return

	}

	err = app.sessManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessManager.Put(r.Context(), "authenticatedUserID", userID)

	app.sessManager.Put(r.Context(), "flash", "Successfully logged in!")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (a *application) serveSignupPage(w http.ResponseWriter, r *http.Request) {

	data := a.newTemplateData(r)

	data.Form = userSignupForm{}

	a.render(w, http.StatusOK, "signup.tmpl.html", data)

}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form userSignupForm

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validation for name field
	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be empty")
	form.CheckField(validator.BelowMaxChars(form.Name, 255), "name", "Name cannot be more than 255 characters long")

	// validation for email
	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.BelowMaxChars(form.Email, 255), "email", "Email cannot be more than 255 characters long")
	form.CheckField(validator.IsValidEmail(form.Email), "email", "Please provide a valid email")

	// validation for password
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password should be at least 8 characters long")
	form.CheckField(validator.BelowMaxChars(form.Password, 19), "password", "Password can't be longer than 18 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)

		return

	}

	err = app.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {

			form.AddFieldError("email", "Email address already in use")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)

		} else {
			app.serverError(w, err)
		}

		return
	}


	app.sessManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)


}

func (a *application) logout(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "user will be logged out with the parameters in post form")

}
