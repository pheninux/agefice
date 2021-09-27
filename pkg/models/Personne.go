package models

import (
	"net/http"
	"strconv"
	"time"
)

type Personne struct {
	Id             int        `gorm:"primary_key,auto_increment" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	Mfa            bool       `json:"mfa"`
	Nom            string     `gorm:"type:varchar(25);not null" json:"nom"`
	Prenom         string     `gorm:"type:varchar(25)" json:"prenom"`
	Age            int        `json:"age"`
	DateNaissance  time.Time  `json:"date_naissance"`
	Tel            string     `json:"tel"`
	Mail           string     `json:"mail"`
	Adresse        string     `json:"adresse"`
	Entreprise     Entreprise `gorm:"foreignkey:EntrepriseId" json:"entreprise"`
	EntrepriseId   int
	Document       []Document  `gorm:"many2many:personnes_documents" json:"document"`
	Formation      []Formation `gorm:"many2many:personnes_formations" json:"formation"`
	FlagMail       int         `json:"flag_mail"`
	Nsocial        string      `json:"nsocial"`
	Status         string      `json:"status"`
	Commentaire    string      `json:"commentaire"`
	StopMail       bool        `json:"stop_mail"`
	Prospection    bool        `json:"prospection"`
	ComProspection string      `json:"com_prospection"`
}

func NewPersonne(r *http.Request) (Personne, error) {
	const (
		layoutISO = "2006-01-02"
		layoutUS  = "January 2, 2006"
	)
	var err error
	mfa := false

	prospection := false
	if r.PostForm.Get("prospection") != "" {
		prospection = true
	}
	if r.PostForm.Get("mfa") != "" {
		mfa = true
	}
	stopM := false
	if r.PostForm.Get("stopMail") != "" {
		stopM = true
	}
	//Date naissance de la personne
	dn, err := time.Parse(layoutISO, r.PostForm.Get("dateN"))
	// date début de la formation
	//dd, _ := time.Parse(layoutISO, r.PostForm.Get("dateDeb"))
	////date fin de la formation
	//df, _ := time.Parse(layoutISO, r.PostForm.Get("dateFin"))
	age, err := strconv.Atoi(r.PostForm.Get("age"))
	//// nombre d'heure de la formation
	//nbrHeurs, _ := strconv.Atoi(r.PostForm.Get("nbrHeures"))
	//// Cout de la formation
	//cout, _ := strconv.ParseFloat(r.PostForm.Get("cout"), 32)
	return Personne{
		Mfa:            mfa,
		CreatedAt:      time.Now(),
		Nom:            r.PostForm.Get("nom"),
		Prenom:         r.PostForm.Get("prenom"),
		Age:            age,
		DateNaissance:  dn,
		Tel:            r.PostForm.Get("tel"),
		Mail:           r.PostForm.Get("mail"),
		Adresse:        r.PostForm.Get("adresse"),
		Nsocial:        r.PostForm.Get("nsocial"),
		Status:         r.PostForm.Get("status"),
		Commentaire:    r.PostForm.Get("commentaire"),
		StopMail:       stopM,
		Prospection:    prospection,
		ComProspection: r.PostForm.Get("comPros"),
	}, err
}
