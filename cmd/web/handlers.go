package main

import (
	"adilhaddad.net/agefice-docs/pkg/forms"
	"adilhaddad.net/agefice-docs/pkg/models"
	"adilhaddad.net/agefice-docs/pkg/models/mysql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (app *application) SaveOrCheckLogin(w http.ResponseWriter, r *http.Request) {

	// parse r body as byte and then to player object

	/*b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}
	u := models.User{}
	type mess struct {
		Msg string
	}
	if err := json.Unmarshal(b, &u); err != nil {
		app.serverError(w, err)
		return
	}
	findUser, err := app.dbModel.GetLoginByLogin(u)
	if err == nil && findUser.User != "" {
		if app.comparePasswords(findUser.PassWord, []byte(u.PassWord)) {
			json.NewEncoder(w).Encode(&struct {
				Msg string `json:"msg"`
			}{Msg: "User identified"})
		} else {
			json.NewEncoder(w).Encode(&struct {
				Msg string `json:"msg"`
			}{Msg: "Mot de passe incorrect"})

		}
	} else {
		json.NewEncoder(w).Encode(&struct {
			Msg string `json:"msg"`
		}{Msg: err.Error()})
		hashedPass := app.hashAndSalt([]byte(u.PassWord))
		u.PassWord = hashedPass

		err = app.dbModel.SaveUser(u)
		if err == nil {
			if err := json.NewEncoder(w).Encode("Erreur lors de la creation du user"); err != nil {
				app.serverError(w, err)
				return
			}
		} else {
			app.serverError(w, err)
			return
		}
	}*/
}

func (app *application) getAll(w http.ResponseWriter, r *http.Request) {

	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.
	p, err := app.dbModel.Latest()

	if err != nil {

		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.ToJson(w, &p, 200)

}

func (app *application) getByMfa(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

		app.serverError(w, err)
		return
	}
	b, err := strconv.ParseBool(string(body))
	if err != nil {

		app.serverError(w, err)
		return
	}
	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.
	p, err := app.dbModel.GetByMfa(b)

	if err != nil {

		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", app.addDefaultData(&templateData{Personnes: p}, r))

	// Use the new render helper.
	//app.ToJson(w, &p, 200)

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	var files []string
	if app.env == "DEV" {
		files = []string{
			"./ui/html/login.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}
	} else {
		files = []string{
			"/var/www/go/deploy/agefice/ui/html/login.page.tmpl",
			"/var/www/go/deploy/agefice/ui/html/base.layout.tmpl",
			"/var/www/go/deploy/agefice/ui/html/footer.partial.tmpl",
		}
	}

	ts, err := template.New("login.page.tmpl").Funcs(functions).ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.
	p, err := app.dbModel.Latest()

	if err != nil {

		app.serverError(w, err)
		return
	}

	var files []string
	if app.env == "DEV" {
		files = []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}
	} else {
		files = []string{
			"/var/www/go/deploy/agefice/ui/html/home.page.tmpl",
			"/var/www/go/deploy/agefice/ui/html/base.layout.tmpl",
			"/var/www/go/deploy/agefice/ui/html/footer.partial.tmpl",
		}
	}

	td := app.addDefaultData(&templateData{Personnes: p}, r)

	ts, err := template.New("home.page.tmpl").Funcs(functions).ParseFiles(files...)

	if err != nil {

		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, app.addCurrentYear(td, r))
	if err != nil {

		app.serverError(w, err)

	}
	// Use the new render helper.
	//app.ToJson(w, td, 200)
	//app.render(w, r, "home.page.tmpl", app.addDefaultData(td, r))
}
func (app *application) showPersonne(w http.ResponseWriter, r *http.Request) {

	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {

		app.notFound(w) // Use the notFound() helper.
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	p, err := app.dbModel.Get(id)
	if err == mysql.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	d, err := app.dbModel.GetDocuments()
	if err != nil {
		app.serverError(w, err)
	}
	f, err := app.dbModel.GetFormations()
	if err != nil {
		app.serverError(w, err)
	}
	prs, err := app.dbModel.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	//convert prs array to json
	pj, err := json.Marshal(prs)
	if err != nil {
		app.serverError(w, err)
	}
	mapDocs := map[int]string{}
	for _, v := range p.Document {
		mapDocs[v.Id] = v.Libelle
	}

	td := app.addDefaultData(&templateData{Documents: &d, Formations: f, FormData: nil, Personne: p, MapDocs: mapDocs, FormatDate: app.FormatDate, JsonPersonnes: string(pj)}, r)
	app.templateData = td
	app.render(w, r, "personneCreate.page.tmpl", app.addCurrentYear(td, r))

	// Use the new render helper.
	//app.render(w, r, "show.page.tmpl", td)

}

func (app *application) createPersonneForm(w http.ResponseWriter, r *http.Request) {

	d, err := app.dbModel.GetDocuments()
	if err != nil {
		app.serverError(w, err)
		return
	}
	f, err := app.dbModel.GetFormations()
	if err != nil {
		app.serverError(w, err)
		return
	}
	p, err := app.dbModel.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	//convert prs array to json
	pj, err := json.Marshal(p)
	if err != nil {
		app.serverError(w, err)
	}

	td := app.addDefaultData(&templateData{Documents: &d, Formations: f, FormData: nil, Personne: &models.Personne{}, JsonPersonnes: string(pj)}, r)
	app.templateData = td
	app.render(w, r, "personneCreate.page.tmpl", app.addCurrentYear(td, r))
}

func (app *application) createPersonneObject(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	fmt.Println(r.Body)
	if err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}
	s := models.Personne{}
	if err := json.Unmarshal(b, &s); err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}

	if err = app.dbModel.Insert(s); err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}

	//http.Redirect(w, r, fmt.Sprintf("/snippet/:%d", id), http.StatusSeeOther)
}

func (app *application) updatePersonneObject(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {

		app.serverError(w, err)
		return
	}
	s := models.Personne{}
	if err := json.Unmarshal(b, &s); err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}

	if err = app.dbModel.Update(s); err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}

	//http.Redirect(w, r, fmt.Sprintf("/snippet/:%d", id), http.StatusSeeOther)
}

func (app *application) createPersonne(w http.ResponseWriter, r *http.Request) {

	// The check of r.Method != "POST" is now superfluous and can be removed.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// check validation form
	ferr := app.validateForm(r)

	//get all formation

	// If there are any validation errors, re-display the create.page.tmpl
	// template passing in the validation errors and previously submitted
	// r.PostForm data.

	if len(ferr) > 0 {
		app.render(w, r, "personneCreate.page.tmpl", app.addDefaultData(&templateData{FormData: r.PostForm, FormErrors: ferr, Personne: &models.Personne{}}, r))
		return
	}
	// formations
	var frms []models.Formation
	id, _ := strconv.Atoi(r.PostForm.Get("idPersonne"))
	if r.PostForm.Get("new") != "" {
		f, err := app.dbModel.GetFormationByStagiaireInfos(models.Personne{Id: id})
		if err != nil {
			f := *new(models.Formation)
			f.Intitule = r.PostForm.Get("new")
			frms, _ = models.NewFormations(r, f)
		} else {
			f.Intitule = r.PostForm.Get("new")
			frms, _ = models.NewFormations(r, f)
		}

	} else {
		f, _ := app.dbModel.GetFormationByStagiaireInfos(models.Personne{Id: id})
		if err != nil {
			f := *new(models.Formation)
			f.Intitule = r.PostForm["formations"][0]
			frms, _ = models.NewFormations(r, f)
		} else {
			f.Intitule = r.PostForm["formations"][0]
			frms, _ = models.NewFormations(r, f)
		}
	}

	// Entreprise
	e := models.NewEntreprise(r)
	//Document
	var docs []models.Document
	fmt.Println(r.PostForm["documents"])
	for _, id := range r.Form["documents"] {
		id, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		d, err := app.dbModel.GetDocumentById(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		docs = append(docs, *d)
	}
	var personneById *models.Personne
	if id > 0 {
		personneById, err = app.dbModel.Get(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	p, err := models.NewPersonne(r)
	if personneById != nil {
		p.FlagMail = personneById.FlagMail
	}
	p.Entreprise = e
	p.Formation = frms
	p.Document = docs

	if err != nil {

		app.serverError(w, err)
		return
	}

	if  id == 0 {
		err = app.dbModel.Insert(p)
		if err != nil {
			app.serverError(w, err)
		} else {
			// Use the Put() method to add a string value ("Your snippet was saved
			// successfully!") and the corresponding key ("flash") to the session
			// data. Note that if there's no existing session for the current user
			// (or their session has expired) then a new, empty, session for them
			// will automatically be created by the session middleware.

			app.session.Put(r, "flash", "Stagiaire crée avec succés")
		}
	} else {
		p.Id = id
		err = app.dbModel.Update(p)
		if err != nil {
			app.serverError(w, err)
		}
	}

	if err != nil {

		app.serverError(w, err)
		return
	}

	// Change the redirect to use the new semantic URL style of /snippet/:id
	//http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	lp, err := app.dbModel.Latest()
	if err != nil {

		app.serverError(w, err)
	}
	app.render(w, r, "home.page.tmpl", app.addDefaultData(&templateData{Personnes: lp}, r))
}

func (app *application) getFormationById(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {

		app.serverError(w, err)
		return
	}
	f, err := app.dbModel.GetFormationById(id)
	if err != nil {

		app.serverError(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(&f); err != nil {

		app.serverError(w, err)
		return
	}
}

func (app *application) getDocuments(w http.ResponseWriter, r *http.Request) {
	d, err := app.dbModel.GetDocuments()
	if err != nil {

		app.serverError(w, err)
		return
	}
	app.ToJson(w, d, 200)
}

func (app *application) getFormations(w http.ResponseWriter, r *http.Request) {
	d, err := app.dbModel.GetFormations()
	if err != nil {

		app.serverError(w, err)
		return
	}
	app.ToJson(w, d, 200)
	fmt.Println(d)
}

func (app *application) deletePersonneById(w http.ResponseWriter, r *http.Request) {

	//id := r.URL.Query().Get(":id")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {

		app.serverError(w, err)
	}
	type data struct {
		Id int `json:"id"`
	}
	d := data{}
	err = json.Unmarshal(b, &d)
	if err != nil {

		app.serverError(w, err)
	}
	err = app.dbModel.DeletePersonneById(strconv.Itoa(d.Id))
	if err != nil {

		app.serverError(w, err)
	}
	/*r.URL.Path = ""
	/*p, err := app.dbModel.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "home.page.tmpl", &templateData{Personnes: p})*/

}

func (app *application) mailProcessing(w http.ResponseWriter, r *http.Request) {

	action := r.URL.Query().Get(":action")
	if action == "true" {
		app.StartMailProcessing()
	} else {
		app.StopService()
	}

}

func (app *application) updateFlag(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	fmt.Println(r.Body)
	if err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}
	s := models.Personne{}
	if err := json.Unmarshal(b, &s); err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}

	if err = app.dbModel.UpdateFlag(&s); err != nil {

		app.serverError(w, err)
		app.ToJson(w, err, 500)
		return
	}

	//http.Redirect(w, r, fmt.Sprintf("/snippet/:%d", id), http.StatusSeeOther)
}

func (app *application) formationCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {

		app.clientError(w, http.StatusBadRequest)
	}
	/** get formation by id **/
	id, err := strconv.Atoi(r.Form["formations"][0])
	if err != nil {

		app.serverError(w, err)
	}

	var formation *models.Formation
	if id != 0 && r.PostForm.Get("new") == "" {
		formation, err = app.dbModel.GetFormationById(id)
		if err != nil {

			app.serverError(w, err)
			return
		}
	}

	f, err := models.NewFormations(r, *formation)
	if err != nil {

		app.serverError(w, err)
	}
	err = app.dbModel.InsertFormation(f[0])
	if err != nil {

		app.serverError(w, err)
	}

}

func (app *application) updatePersonneForm(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {

		app.serverError(w, err)
	}
	var p models.Personne
	if err = json.Unmarshal(b, &p); err != nil {
		app.serverError(w, err)
	}

	d, err := app.dbModel.GetDocuments()
	if err != nil {
		app.serverError(w, err)
	}
	f, err := app.dbModel.GetFormations()
	if err != nil {
		app.serverError(w, err)
	}
	prs, err := app.dbModel.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	//convert prs array to json
	pj, err := json.Marshal(prs)
	if err != nil {
		app.serverError(w, err)
	}
	mapDocs := map[int]string{}
	for _, v := range p.Document {
		mapDocs[v.Id] = v.Libelle
	}

	td := app.addDefaultData(&templateData{Documents: &d, Formations: f, FormData: nil, Personne: &p, MapDocs: mapDocs, FormatDate: app.FormatDate, JsonPersonnes: string(pj)}, r)
	app.templateData = td
	app.render(w, r, "personneCreate.page.tmpl", app.addCurrentYear(td, r))
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{
		Form: forms.New(nil)}, r))
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
		return
	}

	err = app.dbModel.InsertUser(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", app.addDefaultData(&templateData{
		Form: forms.New(nil)}, r))
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error // message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.dbModel.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
			return
		} else {
			form.Errors.Add("generic", "no matching record found")
			app.render(w, r, "login.page.tmpl", app.addDefaultData(&templateData{Form: form}, r))
			return
		}
	}
	// Add the ID of the current user to the session, so that they are now 'logged // in'.
	app.session.Put(r, "authenticatedUserID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been // logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) sendedMails(w http.ResponseWriter, r *http.Request) {

	p, err := app.dbModel.GetSendedMAils()
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "mail.page.tmpl", app.addDefaultData(&templateData{Personnes: p}, r))
}
func (app *application) foreSendMail(w http.ResponseWriter, r *http.Request) {
	var p models.Personne
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
	}
	err = json.Unmarshal(b, &p)
	if err != nil {
		app.serverError(w, err)
	}
	app.sendMail(p)

}
