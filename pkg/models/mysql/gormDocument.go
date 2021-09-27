package mysql

import (
	"adilhaddad.net/agefice-docs/pkg/models"
	_ "github.com/jinzhu/gorm"
)

func (dm *DataModel) GetDocuments() ([]*models.Document, error) {

	var d []*models.Document

	err := dm.Db.Debug().Find(&d).Error

	return d, err

}

func (dm *DataModel) GetDocumentById(id int) (*models.Document, error) {
	d := models.Document{}
	err := dm.Db.Debug().Where("id = ?", &id).Find(&d).Error
	return &d, err
}

func (dm *DataModel) DeleteDocumentsByPersonne(p models.Personne) error {
	err := dm.Db.Debug().Delete(p.Document).Find(p).Error
	return err
}
