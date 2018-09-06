package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
	"time"
)

type BizLog struct {
	Id         string    `orm:"column(log_id);pk" description:"日志表id，uuid"`
	UserId     string    `orm:"column(user_id);size(32);null" description:"用户id,记录操作用户"`
	ModuleName string    `orm:"column(module_name);size(225);null" description:"模块名称"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"操作时间"`
	ClassName  string    `orm:"column(class_name);size(225);null" description:"类名称"`
	MethodName string    `orm:"column(method_name);size(225);null" description:"方法名称"`
	Params     string    `orm:"column(params);null" description:"传入参数"`
	Ip         string    `orm:"column(ip);size(225);null" description:"操作ip"`
	Commemts   string    `orm:"column(commemts);null" description:"备注"`
	Status     int       `orm:"column(status);null"`
}

func (u *BizLog) TableName() string {
	return "com_biz_log"
}

// BizLogAdd insert a new BizLog into database and returns
// last inserted Id on success.
func BizLogAdd(bizLog *BizLog) (int64, error) {

	return orm.NewOrm().Insert(bizLog)
}

// GetBizLogById retrieves BizLog by Id. Returns error if
// Id doesn't exist
func GetBizLogById(id string) (v *BizLog, err error) {
	o := orm.NewOrm()
	v = &BizLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBizLog retrieves all BizLog matches certain condition. Returns empty list if
// no records exist
func GetAllBizLog(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, count int64, err error) {

	con := orm.NewOrm()
	con.Using("default")
	qs := con.QueryTable(new(BizLog))
	//dbt,err :=  orm.GetDB()
	//con.Using("default")
	//dbt.Exec( "CREATE TABLE `com_user_1` ( `id` char(32) NOT NULL, `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',`email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱', `password` char(32) NOT NULL DEFAULT '' COMMENT '密码', `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',`last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间', `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',`status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 -1禁用', PRIMARY KEY (`id`), UNIQUE KEY `idx_user_name` (`user_name`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	//dbt.Exec("DROP TABLE IF EXISTS com_user_1");

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
					return nil, count, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, count, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, count, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, count, errors.New("Error: unused 'order' fields")
		}
	}

	var l []BizLog
	qs = qs.OrderBy(sortFields...)

	count, _ = qs.Count()

	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, count, nil
	}
	return nil, count, err
}

// UpdateBizLog updates BizLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateBizLogById(m *BizLog) (err error) {
	o := orm.NewOrm()
	v := BizLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBizLog deletes BizLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBizLog(id string) (err error) {
	o := orm.NewOrm()
	v := BizLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&BizLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
