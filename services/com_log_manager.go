package services

import (
	"bytes"
	"github.com/astaxie/beego/orm"
	"github.com/beego/bee/generate"
	"logManager/models"
	"logManager/utils"
	"reflect"
	"strings"
)

func ManagerServiceGetLogList(query map[string]string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, titleMap map[string]string, count int64, err error) {

	con := orm.NewOrm()
	con.Using("default")
	qs := con.QueryTable("com_table_mapping")
	titleMap = make(map[string]string)
	commonLogs := []models.CommonLog{}

	//查询默认数据库获取字段配置
	cofigList, err := MappingServiceGetList(query, []string{}, sortby, order)

	//获取查询语句SQL select old_column1 new_column1,old_column2 new_column2 from table_name
	fields := []string{}
	sql, err := getAliasColSql(cofigList, &fields, titleMap)

	if err != nil {
		utils.Logger.Error("MappingServiceGetList failed ", err.Error())
		return nil, titleMap, 0, err
	}

	aliasName := query["aliasName"]
	tableName := query["tableName"]
	sql.WriteString(" from " + tableName + " limit 3 ")

	querySql := sql.String()
	utils.Logger.Info("query log sql :【" + querySql + "】")

	con.Using(aliasName)
	qs = con.QueryTable(tableName)

	//查询数据
	count, _ = qs.Count()
	sn, err := con.Raw(querySql).QueryRows(&commonLogs)

	utils.Logger.Info("query data num  ", sn)

	//过滤未使用的字段 trim unused fields
	if len(fields) == 0 {
		for _, v := range commonLogs {
			ml = append(ml, v)
		}
	} else {
		for _, v := range commonLogs {
			m := make(map[string]interface{})
			val := reflect.ValueOf(v)
			for _, fname := range fields {
				fnameArr := strings.Split(fname, ":")
				fname = fnameArr[0]
				ftype := fnameArr[1]
				if strings.Index(ftype, "int") >= 0 || strings.Index(ftype, "bool") >= 0 {
					m[fname] = val.FieldByName(fname).Int()
				} else if strings.Index(ftype, "string") >= 0 {
					m[fname] = val.FieldByName(fname).String()
				} else if strings.Index(ftype, "float32") >= 0 {
					m[fname] = val.FieldByName(fname).Float()
				} else if strings.Index(ftype, "float64") >= 0 {
					m[fname] = val.FieldByName(fname).Float()
				} else if strings.Index(ftype, "time") >= 0 {
					m[fname] = val.FieldByName(fname).String()
				}
			}
			ml = append(ml, m)
		}
	}

	return ml, titleMap, count, nil
}

func getAliasColSql(cofigList []models.TableMapping, fields *[]string, titleMap map[string]string) (sql bytes.Buffer, err error) {

	sql.WriteString("select ")

	var comonLog models.CommonLog
	mysqlDB := generate.MysqlDB{}

	for index, mapping := range cofigList {

		fieldName := mapping.FieldName
		fieldType := mapping.FieldType
		fieldTitle := mapping.FieldTitle

		col := new(generate.Column)
		col.Type, err = mysqlDB.GetGoDataType(fieldType)
		if err != nil {
			utils.Logger.Error(err.Error())
			return sql, err
		}
		col.Name = fieldName

		aliasMap := utils.ReflectField2Map(&comonLog)

		if strings.Index(col.Type, "int") >= 0 || strings.Index(col.Type, "bool") >= 0 {
			col.Type = "int"
		} else if strings.Index(col.Type, "string") >= 0 {
			col.Type = "string"
		} else if strings.Index(col.Type, "float32") >= 0 {
			col.Type = "float32"
		} else if strings.Index(col.Type, "float64") >= 0 {
			col.Type = "double"
		} else if strings.Index(col.Type, "time") >= 0 {
			col.Type = "time.Time"
		}

		for fname, ftype := range aliasMap {
			if ftype == col.Type {
				if mapping.IsKey == 1 {
					col.Name = col.Name + " id "
					*fields = append(*fields, "Id"+":"+ftype)
					titleMap["Id"] = fieldTitle
				} else {
					col.Name = col.Name + " " + strings.ToLower(fname)
					*fields = append(*fields, fname+":"+ftype)
					titleMap[fname] = fieldTitle
				}

				delete(aliasMap, fname)
				break
			}
		}

		sql.WriteString(col.Name)
		if index < len(cofigList)-1 {
			sql.WriteString(",")
		}
	}

	return sql, nil

}

func ManagerServiceGetById(id string) (bizLog *models.BizLog, err error) {

	bizLog, err = models.GetBizLogById(id)

	if err != nil {
		return
	}

	return bizLog, nil

}

func ManagerServiceUpdate(bizLog *models.BizLog) (err error) {

	err = models.UpdateBizLogById(bizLog)

	if err != nil {
		return
	}

	return

}
