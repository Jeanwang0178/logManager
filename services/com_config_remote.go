package services

import (
	"github.com/satori/go.uuid"
	"logManager/models"
	"strings"
	"time"
)

func ConfigRemoteServiceAdd(m *models.ConfigRemote) (num int64, err error) {

	m.Id = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	m.Pattern = "json"
	m.CreateTime = time.Now()
	num, err = models.AddConfigRemote(m)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func ConfigRemoteServiceUpdate(m *models.ConfigRemote) (err error) {

	m.CreateTime = time.Now()
	err = models.UpdateConfigRemoteById(m)
	if err != nil {
		return err
	}
	return
}

// DeleteConfigRemote deletes ConfigRemote by Id and returns error if
// the record to be deleted doesn't exist
func ConfigRemoteServiceDelete(id string) (err error) {
	err = models.DeleteConfigRemote(id)
	if err != nil {
		return err
	}
	return
}

func ConfigRemoteServiceGetById(id string) (v *models.ConfigRemote, err error) {

	v, err = models.GetConfigRemoteById(id)

	if err != nil {
		return nil, err
	}
	return v, nil

}

func ConfigRemoteServiceGetList(query map[string]string, offset int64, limit int64) (ml []models.ConfigRemote, num int64, err error) {
	ml, num, err = models.GetAllConfigRemote(query, offset, limit)
	return ml, num, nil
}
