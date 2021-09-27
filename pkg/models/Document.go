package models

type Document struct {
	Id      int    `gorm:"primary key;auto_increment" json:"id"`
	Libelle string `gorm:"type:varchar(50);not null" json:"libelle"`
}
