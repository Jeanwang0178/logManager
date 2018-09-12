package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
)

type ConfigRemote struct {
	Id         string    `orm:"column(id);pk" description:"主键"`
	RemoteAddr string    `orm:"column(remoteAddr);size(128)" description:"请求地址"`
	Header     string    `orm:"column(header);size(256)" description:"header"`
	Param      string    `orm:"column(param);size(256)" description:"param"`
	Body       string    `orm:"column(body);size(1024)" description:"body"`
	Method     string    `orm:"column(method);size(4)" description:"方法"`
	Pattern    string    `orm:"column(pattern);size(50)" description:"格式"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
}

func (t *ConfigRemote) TableName() string {
	return "com_config_remote"
}

func init() {
	orm.RegisterModel(new(ConfigRemote))
}

// AddConfigRemote insert a new ConfigRemote into database and returns
// last inserted Id on success.
func AddConfigRemote(m *ConfigRemote) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetConfigRemoteById retrieves ConfigRemote by Id. Returns error if
// Id doesn't exist
func GetConfigRemoteById(id string) (v *ConfigRemote, err error) {
	o := orm.NewOrm()
	v = &ConfigRemote{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllConfigRemote retrieves all ConfigRemote matches certain condition. Returns empty list if
// no records exist
func GetAllConfigRemote(query map[string]string, offset int64, limit int64) (ml []ConfigRemote, num int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ConfigRemote))
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
	// order by:
	var sortFields = []string{"-create_time"}

	num, err = qs.Count()

	if err != nil {
		return nil, num, err
	}

	var l []ConfigRemote
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l); err == nil {

		for _, v := range l {
			ml = append(ml, v)
		}

		return ml, num, nil
	}
	return nil, num, err
}

// UpdateConfigRemote updates ConfigRemote by Id and returns error if
// the record to be updated doesn't exist
func UpdateConfigRemoteById(m *ConfigRemote) (err error) {
	o := orm.NewOrm()
	v := ConfigRemote{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteConfigRemote deletes ConfigRemote by Id and returns error if
// the record to be deleted doesn't exist
func DeleteConfigRemote(id string) (err error) {
	o := orm.NewOrm()
	v := ConfigRemote{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ConfigRemote{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
