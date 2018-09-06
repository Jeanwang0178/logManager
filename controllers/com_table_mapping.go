package controllers

import (
	"encoding/json"
	"errors"
	"logManager/models"
	"logManager/services"
	"logManager/utils"
	"strings"
)

// MappingController operations for TableMapping
type MappingController struct {
	BaseController
}

// URLMapping ...
func (c *MappingController) URLMapping() {

	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @router /add [get,post]
func (ctl *MappingController) MappingEdit() {

	response := make(map[string]interface{})
	var query = make(map[string]string)

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.GetString("tableName")

	query["aliasName"] = aliasName
	query["tableName"] = tableName

	if aliasName != "" && tableName != "" {
		ml, err := services.MappingServiceGetFieldByDatabase(query)
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
		utils.Logger.Info("query param || ", "aliasName：", aliasName, "tableName", tableName)
	}

	ctl.Data["param"] = query
	ctl.Data["result"] = response
	ctl.display()
}

// Post ...
// @Title Post
// @Description create TableMapping
// @Param	body		body 	[]models.TableMapping{}	true		"body for TableMapping content"
// @Success 201 {int} models.TableMapping
// @Failure 403 body is empty

// @router /mappingSave [post]
func (ctl *MappingController) MappingSaveAll() {

	response := make(map[string]interface{})

	vList := []models.TableMapping{}
	var vbyte = ctl.Ctx.Input.RequestBody
	err := json.Unmarshal(vbyte, &vList)
	if err == nil {
		_, err = services.AddAllTableMapping(vList)
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
// @Description get TableMapping by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.TableMapping
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MappingController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id := idStr
	v, err := models.GetTableMappingById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get TableMapping
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.TableMapping
// @Failure 403
// @router / [get]
func (c *MappingController) GetAll() {
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

	l, err := models.GetAllTableMapping(query, fields, sortby, order)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the TableMapping
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.TableMapping	true		"body for TableMapping content"
// @Success 200 {object} models.TableMapping
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MappingController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id := idStr
	v := models.TableMapping{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTableMappingById(&v); err == nil {
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
// @Description delete the TableMapping
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *MappingController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id := idStr
	if err := models.DeleteTableMapping(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
