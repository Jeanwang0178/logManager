package services

import "logManager/models"

func BizLogServiceGetList(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {

	mlist, err := models.GetAllBizLog(query, fields, sortby, order, offset, limit)
	if err != nil {
		return
	}

	return mlist, nil
}
