package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id        string `orm:"column(id);pk"`
	UserName  string `orm:"column(user_name);size(20)" description:"用户名"`
	Email     string `orm:"column(email);size(50)" description:"邮箱"`
	Password  string `orm:"column(password);size(32)" description:"密码"`
	Salt      string `orm:"column(salt);size(10)" description:"密码盐"`
	LastLogin int    `orm:"column(last_login)" description:"最后登录时间"`
	LastIp    string `orm:"column(last_ip);size(15)" description:"最后登录IP"`
	Status    int8   `orm:"column(status)" description:"状态，0正常 -1禁用"`
}

func (u *User) TableName() string {
	return "com_user"
}

func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

// UserAdd insert a new User into database and returns
// last inserted Id on success.
func UserAdd(user *User) (int64, error) {
	return orm.NewOrm().Insert(user)
}

// UserGetById retrieves User by Id. Returns error if
// Id doesn't exist
func UserGetById(id string) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable("com_user").Filter("id", id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UserGetByName retrieves User by userName. Returns error if
// no records exist
func UserGetByName(userName string) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable("com_user").Filter("user_name", userName).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UserUpdate updates User by Id and returns error if
// the record to be updated doesn't exist
func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}
