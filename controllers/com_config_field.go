package controllers

import (
	"encoding/json"
	"errors"
	"logManager/models"
	"logManager/services"
	"logManager/utils"
	"strings"
)

// FieldController operations for ConfigField
type FieldController struct {
	BaseController
}

// URLMapping ...
func (c *FieldController) URLMapping() {

	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @router /edit [get,post]
func (ctl *FieldController) Edit() {

	response := make(map[string]interface{})
	var query = make(map[string]string)

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.GetString("tableName")

	query["aliasName"] = aliasName
	query["tableName"] = tableName

	aliasNames := make([]interface{}, 0)
	err := utils.GetCache(utils.AliasName, &aliasNames)
	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		if aliasName != "" && tableName != "" {
			ml, err := services.ConfigFieldServiceGetFieldByDatabase(query)
			if err != nil {
				utils.Logger.Error(err.Error())
				response["code"] = utils.FailedCode
				response["msg"] = err.Error()
			} else {
				response["code"] = utils.SuccessCode
				response["msg"] = utils.SuccessMsg
				response["data"] = ml
			}

		} else {
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
			utils.Logger.Info("query param || ", "aliasName：", aliasName, "tableName", tableName)
		}
	}

	ctl.Data["param"] = query
	ctl.Data["result"] = response
	ctl.Data["aliasNames"] = aliasNames
	ctl.display()
}

// Post ...
// @Title Post
// @Description create ConfigField
// @Param	body		body 	[]models.ConfigField{}	true		"body for ConfigField content"
// @Success 201 {int} models.ConfigField
// @Failure 403 body is empty

// @router /save [post]
func (ctl *FieldController) Save() {

	response := make(map[string]interface{})

	vList := []models.ConfigField{}
	var vbyte = ctl.Ctx.Input.RequestBody
	err := json.Unmarshal(vbyte, &vList)
	if err == nil {
		_, err = services.ConfigFieldServiceAddAll(vList)
		if err != nil {
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			ctl.Ctx.Output.SetStatus(201)
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
		}

	} else {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	}

	ctl.Data["json"] = response
	ctl.ServeJSON()

}

// GetOne ...
// @Title Get One
// @Description get ConfigField by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConfigField
// @Failure 403 :id is empty
// @router /:id [get]
func (c *FieldController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id := idStr
	v, err := models.GetConfigFieldById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get ConfigField
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ConfigField
// @Failure 403
// @router / [get]
func (c *FieldController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}

	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllConfigField(query, fields, sortby, order)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the ConfigField
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ConfigField	true		"body for ConfigField content"
// @Success 200 {object} models.ConfigField
// @Failure 403 :id is not int
// @router /:id [put]
func (c *FieldController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id := idStr
	v := models.ConfigField{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateConfigFieldById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the ConfigField
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *FieldController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id := idStr
	if err := models.DeleteConfigField(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
