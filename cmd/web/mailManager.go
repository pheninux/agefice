package main

import (
	"adilhaddad.net/agefice-docs/pkg/models"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/matryer/runner"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

const (
	FROM     = "jchassagnac@utas.fr"
	PASSWORD = "Jcre*951"
)

type mail struct {
	host    string
	from    string
	to      string
	object  string
	content string
}

const (
	mailTypePart1 = `Bonjour , <br><br>

Votre formation est à présent terminée.<br> 
Pour la demande de remboursement vous voudrez bien me faire parvenir les éléments suivants :`

	mailTypePart2 = `
L’ensemble des documents doit être signé et tamponné par l’organisme de formation.<br><br>

Je reste à votre disposition pour tout complément d’information.<br><br>

Cordialement,<br>
Jessica CHASSAGNAC<br>
MEDEF Douaisis <br>
03.27.08.10.70`
)

func (app *application) StartMailProcessing() error {
	c := make(chan []models.Personne)
	e := make(chan error, 1)
	app.serviceMail.task = runner.Go(func(s runner.S) error {
		fmt.Println("Starting service mail")
		app.serviceMail.IsStarting = true
		for {
			p, err := app.dbModel.Latest()
			if err != nil {
				e <- err
			}
			c <- p
			time.Sleep(time.Hour * 24)
			if s() {
				break
			}
		}
		return nil
	})

	if app.serviceMail.task.Err() != nil {
		log.Fatalln("task failed:", app.serviceMail.task.Err())
		app.serviceMail.IsStarting = false
		return app.serviceMail.task.Err()
	}
	for {
		select {
		case err, ok := <-e:
			if ok {
				fmt.Println(err)
				return nil
			}
		case p, ok := <-c:
			if ok {
				fmt.Println("starting manage personne")
				go app.managePersonnes(p)
			}
		default:

		}
	}

}
func (app *application) StopService() {
	fmt.Println("stoping service mail")
	app.serviceMail.task.Stop()
	app.serviceMail.IsStarting = false
}

func (app *application) sendMail(p models.Personne) {

	body := app.composeMailBody(p)
	m := gomail.NewMessage()
	m.SetHeader("From", FROM)
	m.SetHeader("To", p.Mail)
	m.SetHeader("Subject", "Dossier AGEFICE")
	m.SetBody("text/html", body)
	mailer := gomail.NewDialer("smtp.office365.com", 587, FROM, PASSWORD)
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email
	if err := mailer.DialAndSend(m); err != nil {
		p.FlagMail = 2
		app.updateFlagMail(p)
		return
	}
	p.FlagMail = 1
	app.updateFlagMail(p)
}

func (app *application) managePersonnes(p []models.Personne) {

	for _, v := range p {
		dNow, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		if err != nil {
			fmt.Println("Error :", err)
			return
		}
		if !v.StopMail {
			if dNow.Sub(v.Formation[0].DateFin).Hours()/24 == 1 && v.FlagMail == 0 {
				go app.sendMail(v)
			}
		}
	}

}

func (app *application) updateFlagMail(p models.Personne) {

	/*_ , err := goreq.Request{
		Method: "POST",
		Uri:    "http://localhost:4000/personne/json/updateFlagMail",
		Body:   &models.Personne{Id:p.Id,FlagMail:true},
	}.Do()*/
	fmt.Println("updating flag mail ...")
	err := app.dbModel.UpdateFlag(&models.Personne{Id: p.Id, FlagMail: p.FlagMail})
	if err != nil {
		fmt.Println("Error [update flag mail] :", err)
	}

}

func (app *application) composeMailBody(p models.Personne) string {
	docs, err := app.dbModel.GetDocuments()
	if err != nil {
		fmt.Println("Error when geting documents [mailManager]")
	}
	//labels all docs
	var l []string
	//labels doc personne
	var lp []string
	for _, v := range docs {
		l = append(l, v.Libelle)
	}
	for _, x := range p.Document {
		lp = append(lp, x.Libelle)
	}
	return composteTemplaceMissingDocs(missingDocs(lp, l))

}

func missingDocs(a, b []string) []string {
	ma := make(map[string]bool)
	for _, ka := range a {
		ma[ka] = true
	}
	var r []string
	for _, kb := range b {
		if !ma[kb] {
			r = append(r, kb)
		}
	}
	return r
}

func composteTemplaceMissingDocs(docs []string) string {
	var buf bytes.Buffer
	buf.WriteString(mailTypePart1)
	buf.WriteString("<ul>")
	for _, v := range docs {
		buf.WriteString(fmt.Sprintf("<li>%s</li>", v))
	}
	buf.WriteString("</ul>")
	buf.WriteString(mailTypePart2)
	return buf.String()
}
