package mysql

import (
	"adilhaddad.net/agefice-docs/pkg/models"
	_ "github.com/jinzhu/gorm"
	"strconv"
)

func (dm DataModel) Insert(p models.Personne) error {
	err := dm.Db.Debug().Create(&p).Error
	return err
}

func (dm DataModel) Update(p models.Personne) error {
	err := dm.Db.Debug().Save(&p).Model(p).Association("Document").Replace(p.Document).Error
	return err
}

func (dm DataModel) UpdateFlag(p *models.Personne) error {
	err := dm.Db.Debug().Model(&p).Update("flag_mail", p.FlagMail).Error
	return err
}

func (dm *DataModel) Get(id int) (*models.Personne, error) {
	p := models.Personne{}
	err := dm.Db.Debug().Preload("Entreprise").Preload("Formation").Preload("Document").First(&p, id).Error
	return &p, err
}

func (dm *DataModel) Latest() ([]models.Personne, error) {

	var p []models.Personne

	err := dm.Db.Debug().Order("id asc").Preload("Entreprise").Preload("Document").Preload("Formation").Find(&p).Error

	return p, err

}
func (dm *DataModel) GetByMfa(b bool) ([]models.Personne, error) {

	var p []models.Personne
	err := dm.Db.Debug().Order("id asc").Preload("Entreprise").Preload("Document").Preload("Formation").Where("mfa = ?", b).Find(&p).Error
	//err := dm.Db.Debug().Where("mfa = ?", b).Find(&p).Error
	return p, err

}

func (dm *DataModel) DeletePersonneById(val string) error {
	id, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	p := &models.Personne{Id: id}
	return dm.Db.Debug().Delete(&p).Error
}

func (dm *DataModel) GetSendedMAils() ([]models.Personne, error) {
	var p []models.Personne
	err := dm.Db.Debug().Order("id asc").Preload("Entreprise").Preload("Document").Preload("Formation").Where("flag_mail = ? or flag_mail = ?", 1, 2).Order("nom asc").Find(&p).Error
	return p, err
}
