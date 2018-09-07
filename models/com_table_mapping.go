package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/beego/bee/generate"
)

type TableMapping struct {
	Id           string `orm:"column(id);pk" description:"主键"`
	AliasName    string `orm:"column(aliasName);size(50)" description:"数据库别名"`
	LogTableName string `orm:"column(log_table_name);size(50)" description:"表名称"`
	FieldName    string `orm:"column(field_name);size(50)" description:"字段名称"`
	FieldType    string `orm:"column(field_type);size(16)" description:"字段类型"`
	FieldTitle   string `orm:"column(field_title);size(50)" description:"字段标题"`
	FieldSort    int    `orm:"column(field_sort)" description:"字段排序"`
	OrderBy      string `orm:"column(order_by)" description:"排序（ASC:升序,DESC:降序)"`

	IsPrimary int    `orm:"column(is_primary)" description:"是否主键"`
	IsShow    string `orm:"column(is_show);size(1)" description:"是否显示"`
	Status    int    `orm:"column(status)" description:"状态，0正常 1禁用"`
}

func (t *TableMapping) TableName() string {
	return "com_table_mapping"
}

func init() {
	orm.RegisterModel(new(TableMapping))
}

// AddTableMapping insert a new TableMapping into database and returns
// last inserted Id on success.
func AddTableMapping(m *TableMapping) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return 0, err
	}
	return id, err
}

func AddAllTableMapping(mds []TableMapping) (int64, error) {
	o := orm.NewOrm()
	aliasName := mds[0].AliasName
	logTableName := mds[1].LogTableName
	if aliasName == "" || logTableName == "" {
		return 0, errors.New("aliasName :" + aliasName + " and logTableName is " + logTableName)
	}
	db, err := orm.GetDB("default")
	db.Exec("delete from com_table_mapping where aliasName = ? and log_table_name = ? ", aliasName, logTableName)
	num, err := o.InsertMulti(len(mds), mds)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// GetTableMappingById retrieves TableMapping by Id. Returns error if
// Id doesn't exist
func GetTableMappingById(id string) (v *TableMapping, err error) {
	o := orm.NewOrm()
	v = &TableMapping{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//根据数据库\表名称 获取数据库字段

func GetFieldByDatabase(query map[string]string) (ml []TableMapping, err error) {

	tableName := query["tableName"]
	aliasName := query["aliasName"]

	ml, err = GetAllTableMapping(query, []string{}, []string{"field_sort"}, []string{"asc"})

	if err == nil && len(ml) > 0 {
		return ml, err
	}

	db, err := orm.GetDB(aliasName)
	trans := &generate.MysqlDB{}
	if err != nil {
		return nil, err
	}

	blackList := make(map[string]bool)

	tb := new(generate.Table)
	tb.Name = tableName
	tb.Fk = make(map[string]*generate.ForeignKey)
	trans.GetColumns(db, tb, blackList)

	columns := tb.Columns

	for index, col := range columns {
		mapping := TableMapping{}
		tag := col.Tag
		mapping.AliasName = aliasName
		mapping.LogTableName = tableName
		mapping.FieldName = tag.Column
		mapping.FieldType = col.Type
		mapping.FieldTitle = tag.Comment
		if mapping.FieldTitle == "" {
			mapping.FieldTitle = tag.Column
		}
		mapping.FieldSort = index
		mapping.IsShow = "1"
		mapping.Status = 0
		if index == 0 {
			mapping.IsPrimary = 1
		} else {
			mapping.IsPrimary = 0
		}
		ml = append(ml, mapping)
	}

	return ml, nil

}

// 获取数据库字段配置列表
func GetAllTableMapping(query map[string]string, fields []string, sortby []string, order []string) (ml []TableMapping, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(TableMapping))

	qs = qs.Filter("aliasName", query["aliasName"])
	qs = qs.Filter("logTableName", query["tableName"])

	//排序
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var list []TableMapping
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.All(&list); err == nil {
		for _, v := range list {
			ml = append(ml, v)
		}

		return ml, nil
	}
	return nil, err
}

// UpdateTableMapping updates TableMapping by Id and returns error if
// the record to be updated doesn't exist
func UpdateTableMappingById(m *TableMapping) (err error) {
	o := orm.NewOrm()
	v := TableMapping{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTableMapping deletes TableMapping by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTableMapping(id string) (err error) {
	o := orm.NewOrm()
	v := TableMapping{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TableMapping{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
