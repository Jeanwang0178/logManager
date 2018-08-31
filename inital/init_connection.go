package inital

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"logManager/models"
	"net/url"
)

func Init() {

	//默认数据库
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)

	//日志数据库
	ecdbhost := beego.AppConfig.String("ec.db.host")
	ecdbport := beego.AppConfig.String("ec.db.port")
	ecdbuser := beego.AppConfig.String("ec.db.user")
	ecdbpassword := beego.AppConfig.String("ec.db.password")
	ecdbname := beego.AppConfig.String("ec.db.name")
	ectimezone := beego.AppConfig.String("ec.db.timezone")

	ec := ecdbuser + ":" + ecdbpassword + "@tcp(" + ecdbhost + ":" + ecdbport + ")/" + ecdbname + "?charset=utf8"

	if ectimezone != "" {
		ec = ec + "&loc=" + url.QueryEscape(ectimezone)
	}

	orm.RegisterDataBase("common", "mysql", ec)

	orm.RegisterModel(new(models.User), new(models.BizLog))

	if beego.AppConfig.String("runmode") == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	}

}
