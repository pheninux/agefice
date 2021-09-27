package mysql

import (
	"adilhaddad.net/agefice-docs/pkg/models"
	_ "github.com/jinzhu/gorm"
)

func (dm DataModel) InsertFormation(p models.Formation) error {
	err := dm.Db.Debug().Create(&p).Error
	return err
}

func (dm *DataModel) GetFormationById(id int) (*models.Formation, error) {
	f := models.Formation{}
	err := dm.Db.Debug().Where("id = ?", &id).Find(&f).Error

	return &f, err

}

func (dm *DataModel) GetFormations() ([]models.Formation, error) {

	var f []models.Formation

	//err := dm.Db.Debug().Find(&f).Error
	err := dm.Db.Debug().Table("formations").Select("distinct(intitule)").Find(&f).Error

	return f, err

}

func (dm *DataModel) GetFormationByStagiaireInfos(p models.Personne) (models.Formation, error) {

	err := dm.Db.Debug().Preload("Formation").Where("id = ?", p.Id).Find(&p).Error
	if err != nil {
		return models.Formation{}, err
	}
	return p.Formation[0], err

}
