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
)

func ManagerServiceGetLogList(query map[string]string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, titleMap map[string]string, sortFields []string, count int64, err error) {

	con := orm.NewOrm()
	con.Using("default")
	qs := con.QueryTable("com_table_mapping")
	titleMap = make(map[string]string) //标题map
	commonLogs := []models.CommonLog{}

	//查询默认数据库获取字段配置
	cofigList, err := MappingServiceGetList(query, []string{}, []string{"field_sort"}, []string{"asc"})

	if err != nil {
		utils.Logger.Error("MappingServiceGetList failed ", err.Error())
		return nil, titleMap, sortFields, 0, err
	}

	if len(cofigList) == 0 {
		return nil, titleMap, sortFields, 0, errors.New("请先配置表【" + query["tableName"] + "】字段")
	}

	//获取查询语句SQL select old_column1 new_column1,old_column2 new_column2 from table_name
	fields := []string{}    // 收集使用的字段/类型
	sortFields = []string{} //排序
	sql, err := getAliasColSql(cofigList, &fields, titleMap, &sortFields)

	if err != nil {
		utils.Logger.Error("getAliasColSql failed ", err.Error())
		return nil, titleMap, sortFields, 0, err
	}

	aliasName := query["aliasName"]
	tableName := query["tableName"]

	querySql := sql.String()
	querySql = beego.Substr(querySql, 0, len(querySql)-1)
	querySql += " from " + tableName + " limit 3 "

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
				} else if strings.Index(ftype, "float") >= 0 {
					m[fname] = val.FieldByName(fname).Float()
				} else if strings.Index(ftype, "time") >= 0 {
					m[fname] = val.FieldByName(fname).String()
				}
			}
			ml = append(ml, m)
		}
	}

	return ml, titleMap, sortFields, count, nil
}

func getAliasColSql(cofigList []models.TableMapping, fields *[]string, titleMap map[string]string, sortFields *[]string) (sql bytes.Buffer, err error) {

	sql.WriteString("select ")

	var comonLog models.CommonLog

	aliasMap := utils.ReflectField2Map(&comonLog)

	for _, mapping := range cofigList {

		if mapping.IsShow == "0" {
			continue
		}
		fieldName := mapping.FieldName
		fieldType := mapping.FieldType
		fieldTitle := mapping.FieldTitle

		col := new(generate.Column)

		// mysqlDB := generate.MysqlDB{} mysqlDB.GetGoDataType(fieldType) 获取mysql对应go字段类型
		col.Type = fieldType
		col.Name = fieldName

		if err != nil {
			utils.Logger.Error(err.Error())
			return sql, err
		}

		for fname, ftype := range aliasMap {
			if ftype == col.Type {
				if mapping.IsPrimary == 1 {
					fname = "Id"
				}

				col.Name = col.Name + " " + strings.ToLower(fname)
				*fields = append(*fields, fname+":"+ftype)
				*sortFields = append(*sortFields, fname)
				titleMap[fname] = fieldTitle
				delete(aliasMap, fname)

				break
			}
		}

		sql.WriteString(col.Name + ",")
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
