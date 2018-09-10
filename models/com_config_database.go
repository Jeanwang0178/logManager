package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type ConfigDatabase struct {
	Id         string `orm:"column(id);pk" description:"主键"`
	AliasName  string `orm:"column(aliasName);size(50)" description:"别名"`
	DbHost     string `orm:"column(db_host);size(24)" description:"主机地址"`
	DbUser     string `orm:"column(db_user);size(24)" description:"用户名"`
	DbPassword string `orm:"column(db_password);size(24)" description:"密码"`
	DbPort     string `orm:"column(db_port);size(8)" description:"端口"`
	DbName     string `orm:"column(db_name);size(24)" description:"数据库名称"`
	DbCharset  string `orm:"column(db_charset);size(24)" description:"编码"`
	DbTimezone string `orm:"column(db_timezone);size(24)" description:"数据库时区"`
	DbMaxIdle  int8   `orm:"column(db_maxIdle)" description:"最大空闲链接"`
	DbMaxConn  int8   `orm:"column(db_maxConn)" description:"最大数据库链接"`
	Status     int8   `orm:"column(status)" description:"状态，0正常 1禁用"`
}

func (t *ConfigDatabase) TableName() string {
	return "com_config_database"
}

func init() {
	orm.RegisterModel(new(ConfigDatabase))
}

// AddConfigDatabase insert a new ConfigDatabase into database and returns
// last inserted Id on success.
func AddConfigDatabase(m *ConfigDatabase) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.Insert(m)
	return num, err
}

// GetConfigDatabaseById retrieves ConfigDatabase by Id. Returns error if
// Id doesn't exist
func GetConfigDatabaseById(id string) (v *ConfigDatabase, err error) {
	o := orm.NewOrm()
	v = &ConfigDatabase{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllConfigDatabase retrieves all ConfigDatabase matches certain condition. Returns empty list if
// no records exist
func GetAllConfigDatabase(query map[string]string) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ConfigDatabase))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}

	var l []ConfigDatabase

	for _, v := range l {
		ml = append(ml, v)
	}
	return ml, nil
	return nil, err
}

// UpdateConfigDatabase updates ConfigDatabase by Id and returns error if
// the record to be updated doesn't exist
func UpdateConfigDatabaseById(m *ConfigDatabase) (err error) {
	o := orm.NewOrm()
	v := ConfigDatabase{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteConfigDatabase deletes ConfigDatabase by Id and returns error if
// the record to be deleted doesn't exist
func DeleteConfigDatabase(id string) (err error) {
	o := orm.NewOrm()
	v := ConfigDatabase{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ConfigDatabase{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
