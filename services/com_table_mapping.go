package services

import (
	"github.com/satori/go.uuid"
	"logManager/models"
	"strings"
)

func AddAllTableMapping(mds []models.TableMapping) (num int64, err error) {

	for index, _ := range mds {
		if mds[index].Id == "" {
			mds[index].Id = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
		}
	}

	num, err = models.AddAllTableMapping(mds)
	if err != nil {
		return 0, err
	}
	return num, nil

}

func MappingServiceGetList(query map[string]string, fields []string, sortby []string, order []string) (ml []models.TableMapping, err error) {

	mlist, err := models.GetAllTableMapping(query, fields, sortby, order)
	if err != nil {
		return nil, nil
	}

	return mlist, nil
}

func MappingServiceGetFieldByDatabase(query map[string]string) (ml []models.TableMapping, err error) {

	ml, err = models.GetFieldByDatabase(query)

	return ml, err

}

func MappingServiceGetById(id string) (TableMapping *models.TableMapping, err error) {

	TableMapping, err = models.GetTableMappingById(id)

	if err != nil {
		return
	}

	return TableMapping, nil

}

func MappingServiceUpdate(TableMapping *models.TableMapping) (err error) {

	err = models.UpdateTableMappingById(TableMapping)

	if err != nil {
		return
	}

	return

}
