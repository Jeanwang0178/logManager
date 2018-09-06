package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type TableMapping struct {
	Id           string `orm:"column(id);pk" description:"主键"`
	AliasName    string `orm:"column(aliasName);size(50)" description:"数据库别名"`
	LogTableName string `orm:"column(log_table_name);size(50)" description:"表名称"`
	FieldName    string `orm:"column(field_name);size(50)" description:"字段名称"`
	FieldType    string `orm:"column(field_type);size(16)" description:"字段类型"`
	FieldTitle   string `orm:"column(field_title);size(50)" description:"字段标题"`
	FieldSort    int8   `orm:"column(field_sort)" description:"字段排序"`
	IsKey        int8   `orm:"column(is_key)" description:"是否主键"`
	IsShow       string `orm:"column(is_show);size(1)" description:"是否显示"`
	IsQuery      string `orm:"column(is_query);size(1)" description:"是否查询"`
	Status       int8   `orm:"column(status)" description:"状态，0正常 1禁用"`
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

// 获取字段配置列表
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
