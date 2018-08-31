package services

import (
	"logManager/models"
)

func BizLogServiceGetList(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, count int64, err error) {

	mlist, count, err := models.GetAllBizLog(query, fields, sortby, order, offset, limit)
	if err != nil {
		return
	}

	return mlist, count, nil
}

func BizLogServiceGetById(id string) (bizLog *models.BizLog, err error) {

	bizLog, err = models.GetBizLogById(id)

	if err != nil {
		return
	}

	return bizLog, nil

}

func BizLogServiceUpdate(bizLog *models.BizLog) (err error) {

	err = models.UpdateBizLogById(bizLog)

	if err != nil {
		return
	}

	return

}
