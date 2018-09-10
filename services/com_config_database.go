package services

import (
	"github.com/satori/go.uuid"
	"logManager/models"
	"strings"
)

func ConfigDatabaseServiceAdd(m *models.ConfigDatabase) (num int64, err error) {

	m.Id = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	num, err = models.AddConfigDatabase(m)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func ConfigDatabaseServiceUpdate(m *models.ConfigDatabase) (err error) {

	err = models.UpdateConfigDatabaseById(m)
	if err != nil {
		return err
	}
	return
}

// DeleteConfigDatabase deletes ConfigDatabase by Id and returns error if
// the record to be deleted doesn't exist
func ConfigDatabaseServiceDelete(id string) (err error) {
	err = models.DeleteConfigDatabase(id)
	if err != nil {
		return err
	}
	return
}

func ConfigDatabaseServiceGetById(id string) (v *models.ConfigDatabase, err error) {

	v, err = models.GetConfigDatabaseById(id)

	if err != nil {
		return nil, err
	}
	return v, nil

}

func ConfigDatabaseServiceGetList(query map[string]string) (ml []interface{}, err error) {
	ml, err = models.GetAllConfigDatabase(query)
	if err != nil {
		return nil, err
	}
	return ml, nil
}
