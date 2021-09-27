package models

import (
	"net/http"
	"strconv"
	"time"
)

type Formation struct {
	Id        int `gorm:"primary key;auto_increment" json:"id"`
	CreatedAt time.Time
	Intitule  string    `json:"intitule"`
	DateDebut time.Time `json:"date_debut"`
	DateFin   time.Time `json:"date_fin"`
	NbrHeures int       `json:"nbr_heures"`
	Cout      float64   `json:"cout"`
}

func NewFormations(r *http.Request, formation Formation) ([]Formation, error) {

	const (
		layoutISO = "2006-01-02"
		layoutUS  = "January 2, 2006"
	)
	var err error
	// Id formation
	//id , err :=  strconv.Atoi(r.PostForm.Get("formations"))

	//id, err := strconv.Atoi(r.Form["formations"][0])
	if err != nil {
		return nil, err
	}

	// date début de la formation
	dd, err := time.Parse(layoutISO, r.PostForm.Get("dateDeb"))
	//date fin de la formation
	df, err := time.Parse(layoutISO, r.PostForm.Get("dateFin"))
	// nombre d'heure de la formation
	nbrHeurs, err := strconv.Atoi(r.PostForm.Get("nbrHeures"))
	// Cout de la formation
	cout, err := strconv.ParseFloat(r.PostForm.Get("cout"), 32)

	formation.CreatedAt = time.Now()
	formation.DateDebut = dd
	formation.DateFin = df
	formation.Cout = cout
	formation.NbrHeures = nbrHeurs

	var frms []Formation
	frms = append(frms, formation)
	return frms, err

}
