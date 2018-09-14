package services

import (
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"logManager/models"
	"logManager/utils"
	"net/url"
	"strings"
)

func ConfigDatabaseServiceAdd(m *models.ConfigDatabase) (num int64, err error) {

	ec, err := RegisterDB(m)

	if err != nil {
		utils.Logger.Error("orm.RegisterDataBase failed %s", ec)
		return 0, err
	}

	m.Id = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	num, err = models.AddConfigDatabase(m)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func ConfigDatabaseServiceUpdate(m *models.ConfigDatabase) (err error) {

	ec, err := RegisterDB(m)

	if err != nil {
		utils.Logger.Error("orm.RegisterDataBase failed %s", ec)
		return err
	}

	err = models.UpdateConfigDatabaseById(m)
	if err != nil {
		return err
	}
	return
}

func RegisterDB(m *models.ConfigDatabase) (string, error) {
	ecdbhost := m.DbHost
	ecdbport := m.DbPort
	ecdbuser := m.DbUser
	ecdbpassword := m.DbPassword
	ecdbname := m.DbName
	ectimezone := m.DbTimezone
	ecMaxIdle := int(m.DbMaxIdle)
	ecMaxConn := int(m.DbMaxConn)
	if ecMaxIdle == 0 {
		ecMaxIdle = 2
	}
	if ecMaxConn == 0 {
		ecMaxConn = 3
	}
	ec := ecdbuser + ":" + ecdbpassword + "@tcp(" + ecdbhost + ":" + ecdbport + ")/" + ecdbname + "?charset=utf8"
	if ectimezone != "" {
		ec = ec + "&loc=" + url.QueryEscape(ectimezone)
	}
	err := orm.RegisterDataBase(m.AliasName, "mysql", ec, ecMaxIdle, ecMaxConn)

	if err == nil {
		aliasNames := make([]interface{}, 0)
		err := utils.GetCache(utils.AliasName, &aliasNames)
		if err != nil {
			utils.Logger.Error("utils.GetCache(utils.AliasName)", err)
		}

		aliasNames = append(aliasNames, m.AliasName)
		utils.SetCache(utils.AliasName, aliasNames, 6000000)
	}

	return ec, err
}

// DeleteConfigDatabase deletes ConfigDatabase by Id and returns error if
// the record to be deleted doesn't exist
func ConfigDatabaseServiceDelete(id string) (err error) {

	m, err := models.GetConfigDatabaseById(id)

	if err != nil {
		return err
	} else {
		newData := make([]interface{}, 0)
		aliasNames := make([]interface{}, 0)
		err := utils.GetCache(utils.AliasName, &aliasNames)
		for _, alias := range aliasNames {
			if alias != m.AliasName {
				newData = append(newData, alias)
			}
		}

		utils.SetCache(utils.AliasName, newData, 6000000)
		err = models.DeleteConfigDatabase(id)

		if err != nil {
			return err
		}
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

func ConfigDatabaseServiceGetList(query map[string]string) (ml []models.ConfigDatabase, err error) {
	ml, err = models.GetAllConfigDatabase(query)
	if err != nil {
		return nil, err
	}
	return ml, nil
}
