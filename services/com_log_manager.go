package services

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/bee/generate"
	"logManager/models"
	"logManager/utils"
	"reflect"
	"strings"
	"time"
)

func ManagerServiceGetDataList(query map[string]string, offset int64, limit int64) (retArray []interface{},
	titleMap map[string]string, fieldsSort []string, count int64, err error) {

	con := orm.NewOrm()
	con.Using("default")

	commonLogs := []models.CommonLog{}

	//查询默认数据库获取字段配置
	cofigList, err := ConfigFieldServiceGetList(query, []string{}, []string{"field_sort"}, []string{"asc"})

	if err != nil {
		utils.Logger.Error("MappingServiceGetList failed ", err.Error())
		return nil, titleMap, fieldsSort, 0, err
	}

	if len(cofigList) == 0 {
		return nil, titleMap, fieldsSort, 0, errors.New("请先配置表【" + query["tableName"] + "】字段")
	}

	//获取查询语句SQL select old_column1 new_column1,old_column2 new_column2 from table_name
	//fieldsSort = []string{} 排序  收集使用的字段/类型
	titleMap, fieldsSort, sql, orderBy, err := getAliasColSql(cofigList, true)

	if err != nil {
		utils.Logger.Error("getAliasColSql failed ", err.Error())
		return nil, titleMap, fieldsSort, 0, err
	}

	aliasName := query["aliasName"]
	tableName := query["tableName"]

	querySql := sql.String()
	querySql = beego.Substr(querySql, 0, len(querySql)-1)
	querySql += " from " + tableName + orderBy + " limit ? offset ?  "

	utils.Logger.Info("query log sql :【" + querySql + "】")

	con.Using(aliasName)
	qs := con.QueryTable(tableName)

	//查询数据
	count, _ = qs.Count()
	sn, err := con.Raw(querySql, limit, offset).QueryRows(&commonLogs)

	utils.Logger.Info("query data num  ", sn)

	//过滤未使用的字段 trim unused fields
	for _, data := range commonLogs {

		dataMap := filterMapFields(data, fieldsSort)
		retArray = append(retArray, dataMap)
	}

	return retArray, titleMap, fieldsSort, count, nil
}

func getAliasColSql(configList []models.ConfigField, filterShow bool) (titleMap map[string]string, fieldsSort []string, sql bytes.Buffer, orderBy string, err error) {

	sql.WriteString("select ")

	var comonLog models.CommonLog
	titleMap = make(map[string]string) //标题map
	aliasMap := utils.ReflectField2Map(&comonLog)

	for _, config := range configList {

		if config.OrderBy == "DESC" {
			orderBy += config.FieldName + " " + config.OrderBy + ","
		}
		if len(orderBy) > 0 {
			orderBy = beego.Substr(orderBy, 0, len(orderBy)-1)
			orderBy = " order by " + orderBy
		}

		if filterShow && config.IsShow == "0" && config.IsPrimary != 1 {
			continue
		}
		fieldName := config.FieldName
		fieldType := config.FieldType
		fieldTitle := config.FieldTitle

		col := new(generate.Column)

		// mysqlDB := generate.MysqlDB{} mysqlDB.GetGoDataType(fieldType) 获取mysql对应go字段类型
		col.Type = fieldType
		col.Name = fieldName

		if err != nil {
			utils.Logger.Error(err.Error())
			return nil, nil, sql, orderBy, err
		}

		for fname, ftype := range aliasMap {
			if ftype == col.Type {
				if config.IsPrimary == 1 {
					fname = "Id"
				}

				col.Name = col.Name + " " + strings.ToLower(fname)
				fieldsSort = append(fieldsSort, fname)
				titleMap[fname] = fieldTitle
				delete(aliasMap, fname)

				break
			}
		}

		sql.WriteString(col.Name + ",")
	}

	return titleMap, fieldsSort, sql, orderBy, nil

}

/**
  根据ID查询数据详细信息
*/
func ManagerServiceGetDataById(query map[string]string) (dataMap map[string]interface{}, titleMap map[string]string, fieldsSort []string, err error) {

	con := orm.NewOrm()
	con.Using("default")
	commonLogs := models.CommonLog{}

	//查询默认数据库获取字段配置
	cofigList, err := ConfigFieldServiceGetList(query, []string{}, []string{"field_sort"}, []string{"asc"})

	if err != nil {
		utils.Logger.Error("ConfigFieldServiceGetList failed ", err.Error())
		return nil, titleMap, fieldsSort, err
	}

	if len(cofigList) == 0 {
		return nil, titleMap, fieldsSort, errors.New("请先配置表【" + query["tableName"] + "】字段")
	}

	fieldId := ""
	for _, config := range cofigList {
		if config.IsPrimary == 1 {
			fieldId = config.FieldName
			break
		}
	}
	if fieldId == "" {
		fieldId = cofigList[0].FieldName
	}

	//获取查询语句SQL select old_column1 new_column1,old_column2 new_column2 from table_name where id = ?
	//fieldsSort = []string{} //排序  收集使用的字段/类型
	titleMap, fieldsSort, sql, _, err := getAliasColSql(cofigList, false)

	if err != nil {
		utils.Logger.Error("getAliasColSql failed ", err.Error())
		return nil, titleMap, fieldsSort, err
	}

	aliasName := query["aliasName"]
	tableName := query["tableName"]
	id := query["id"]

	querySql := sql.String()
	querySql = beego.Substr(querySql, 0, len(querySql)-1)
	querySql += " from " + tableName + " where " + fieldId + " = ?  "

	utils.Logger.Info("query log sql :【" + querySql + "】")

	con.Using(aliasName)

	//查询数据
	err = con.Raw(querySql, id).QueryRow(&commonLogs)
	if err != nil {
		return nil, titleMap, fieldsSort, err
	}

	utils.Logger.Info("query data   ", commonLogs)

	//过滤未使用的字段 trim unused fields
	dataMap = filterMapFields(commonLogs, fieldsSort)

	return dataMap, titleMap, fieldsSort, nil
}

/**
过滤未使用的字段
*/
func filterMapFields(commonLogs models.CommonLog, fields []string) (dataMap map[string]interface{}) {
	if dataMap == nil {
		dataMap = make(map[string]interface{})
	}
	val := reflect.ValueOf(commonLogs)
	for _, fname := range fields {
		vtype := val.FieldByName(fname).Type().String()
		if strings.Index(vtype, "time") < 0 {
			dataMap[fname] = val.FieldByName(fname).Interface()
		} else {
			vlue := val.FieldByName(fname).Interface()
			val := vlue.(time.Time).Format("2006-01-02 15:04:05")
			dataMap[fname] = val

		}

		/*if strings.Index(ftype, "int") >= 0 || strings.Index(ftype, "bool") >= 0 {
			dataMap[fname] = val.FieldByName(fname).Int()
		} else if strings.Index(ftype, "string") >= 0 {
			dataMap[fname] = val.FieldByName(fname).String()
		} else if strings.Index(ftype, "float") >= 0 {
			dataMap[fname] = val.FieldByName(fname).Float()
		} else if strings.Index(ftype, "time") >= 0 {
			vlue :=  val.FieldByName(fname).Interface()
			val := vlue.(time.Time).Format("2006-01-02 15:04:05")
			dataMap[fname] = val
		}*/
	}
	return dataMap
}
