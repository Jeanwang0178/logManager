package services

import (
	"github.com/satori/go.uuid"
	"logManager/src/models"
	"strings"
)

func ConfigFieldServiceAddAll(mds []models.ConfigField) (num int64, err error) {

	for index, _ := range mds {
		if mds[index].Id == "" {
			mds[index].Id = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
		}
	}

	num, err = models.AddAllConfigField(mds)
	if err != nil {
		return 0, err
	}
	return num, nil

}

func ConfigFieldServiceGetList(query map[string]string, fields []string, sortby []string, order []string) (ml []models.ConfigField, err error) {

	mlist, err := models.GetAllConfigField(query, fields, sortby, order)
	if err != nil {
		return nil, nil
	}

	return mlist, nil
}

func ConfigFieldServiceGetFieldByDatabase(query map[string]string) (ml []models.ConfigField, err error) {

	ml, err = models.GetFieldByDatabase(query)

	return ml, err

}

func ConfigFieldServiceGetById(id string) (ConfigField *models.ConfigField, err error) {

	ConfigField, err = models.GetConfigFieldById(id)

	if err != nil {
		return
	}

	return ConfigField, nil

}

func ConfigFieldServiceUpdate(ConfigField *models.ConfigField) (err error) {

	err = models.UpdateConfigFieldById(ConfigField)

	if err != nil {
		return
	}

	return

}
