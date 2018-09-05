package services

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"logManager/models"
	"reflect"
	"strings"
)

func ManagerServiceGetList(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, count int64, err error) {

	con := orm.NewOrm()
	qs := con.QueryTable("biz_log")
	logDB, err := orm.GetDB()
	con.Using("default")
	logDB.Exec("CREATE TABLE `com_user_1` ( `id` char(32) NOT NULL, `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',`email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱', `password` char(32) NOT NULL DEFAULT '' COMMENT '密码', `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',`last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间', `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',`status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 -1禁用', PRIMARY KEY (`id`), UNIQUE KEY `idx_user_name` (`user_name`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	logDB.Exec("DROP TABLE IF EXISTS com_user_1")

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

	//qs = qs.OrderBy(sortFields...)

	count, _ = qs.Count()

	var lists []orm.ParamsList

	var logList []models.CommonLog

	colnames := []string{"id2", "userid2", "modulename3"}

	sn, err := con.Raw("select log_id id ,user_id extstr1,module_name extstr2 from com_biz_log limit 3").QueryRows(&commonLog)

	fmt.Println(sn)

	n, err := con.Raw("select log_id id2 ,user_id userid2,module_name modulename2 from com_biz_log ").ValuesList(&lists, colnames...)

	fmt.Println(n)

	if _, err = qs.Limit(limit, offset).All(&logList, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range logList {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range logList {
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

	mlist, count, err := models.GetAllBizLog(query, fields, sortby, order, offset, limit)
	if err != nil {
		return
	}

	return mlist, count, nil
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
