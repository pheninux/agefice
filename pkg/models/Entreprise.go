package models

import (
	"net/http"
	"strconv"
)

type Entreprise struct {
	Id   int    `gorm:"primary_key;auto_increment:true" json:"id"`
	Nom  string `gorm:"type:varchar(50)" json:"nom"`
	Code string `gorm:"type:varchar(8);unique" json:"code"`
}

func NewEntreprise(r *http.Request) Entreprise {

	if r.PostForm.Get("action") == "Modifier" {
		id, _ := strconv.Atoi(r.PostForm.Get("idEntreprise"))
		return Entreprise{Id: id, Nom: r.PostForm.Get("nomEntreprise"), Code: r.PostForm.Get("code")}
	} else {
		return Entreprise{Nom: r.PostForm.Get("nomEntreprise"), Code: r.PostForm.Get("code")}
	}

}
