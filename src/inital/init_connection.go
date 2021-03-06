package inital

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"logManager/src/common"
	"logManager/src/models"
	"logManager/src/services"
	"logManager/src/utils"
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

	orm.RegisterModel(new(models.User), new(models.BizLog))

	if beego.AppConfig.String("runmode") == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	}

	utils.InitCache()

	data := make([]interface{}, 0)
	data = append(data, "default")
	err := utils.SetCache(common.AliasName, data, 6000000)
	if err != nil {
		common.Logger.Error("conect redis failed ", err.Error())
	}

	//初始化已经存在的数据库链接
	query := make(map[string]string)
	ml, err := services.ConfigDatabaseServiceGetList(query)
	if err != nil {
		common.Logger.Error("services.ConfigDatabaseServiceGetList faile ", err)
	}
	for _, dbcofig := range ml {
		conn, err := services.RegisterDB(&dbcofig)
		if err != nil {
			common.Logger.Error("init one database faield ", conn)
		}
	}

	//初始化KAFKA链接
	models.InitKafka()

}
