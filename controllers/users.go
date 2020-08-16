package controllers

import (
	"MyGoApi/models"
	"MyGoApi/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"strconv"
	"strings"
	"time"
)

// UsersController operations for Users
type UsersController struct {
	beego.Controller
}

// URLMapping ...
func (c *UsersController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}



// Post ...
// @Title Post
// @Description create Users
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 403 body is empty
// @router / [post]
func (c *UsersController) Post() {
	var v models.Users
	//if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
	//	if _, err := models.AddUsers(&v); err == nil {
	//		c.Ctx.Output.SetStatus(201)
	//		c.Data["json"] = v
	//	} else {
	//		c.Data["json"] = err.Error()
	//	}
	//} else {
	//	c.Data["json"] = err.Error()
	//}
	//c.ServeJSON()
	valid := validation.Validation{}
	valid.Required(v.Name, "name")
	valid.Required(v.Email, "Email")
	valid.Required(v.Password, "Password")
	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}
}

// GetOne ...
// @Title Get One
// @Description get Users by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Users
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsersController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUsersById(id)
	if err != nil {
		c.Data["json"] = Error(err.Error())
	} else {
		c.Data["json"] = Success(v)
	}
	c.ServeJSON()
}



// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Users
// @Failure 403
// @router / [get]
func (c *UsersController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
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

	l, err := models.GetAllUsers(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Users
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 200 {object} models.Users
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UsersController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Users{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUsersById(&v); err == nil {
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
// @Description delete the Users
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UsersController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteUsers(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @router /login [post]
func (c *UsersController) Login() {
	var data LoginVerify
	user,_ := models.GetUsersById(1)
	var v utils.User
	v.Id = user.Id
	v.Name = user.Name
	var token Token
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &data); err == nil {
		beego.Info(data)
		if err := utils.CheckLogin(data.Name,data.Password); err != "ok"{
			c.Data["json"] = Error(err)
		}else{
			token.Token = utils.GenerateToken(0,v)
			c.Data["json"] = Success(token)
		}
	} else {
		c.Data["json"] = Error(err.Error())
	}

	c.ServeJSON()
}

// @router /register [post]
func (c *UsersController) Register() {

	var data RegisterVerify
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &data); err == nil {
		beego.Info(data)
		if err := utils.CheckRegister(data.Name,data.Email,data.Password,data.Repassword); err != "ok"{
			c.Data["json"] = Error(err)
		}else{
			c.Data["json"] = Success("登录成功")
		}
	} else {
		c.Data["json"] = Error(err.Error())
	}

	c.ServeJSON()
}

// @router /info [get]
func (c *UsersController) UserInfo() {
	token := c.Ctx.Input.Header("Authorization")

	v,_ := utils.ValidateToken(token)

	var maps []orm.Params //[map, map, map]
	type User struct {
		Id        int64     `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		AvatarUrl string    `json:"avatar_url"`
		Captcha   int       `json:"captcha" `
		Intro     string    `json:"intro"`
		Admin     bool      `json:"admin"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	orm := orm.NewOrm()
	_, err := orm.QueryTable("users").Filter("id",v.User.Id).Values(&maps, )
	if err != nil {
		c.Data["json"] = Error("查询出错")
		return
	}

	var user User
	for _, m := range maps {
		a,ok := m["Phone"].(string)
		if ok {
			user.Phone     = a

		}else{
			user.Phone = ""
		}
		user.Id        = m["Id"].(int64)
		user.Name      = m["Name"].(string)
		//user.Phone     = m["Phone"].(string)
		user.Email     = m["Email"].(string)

		//user.AvatarUrl = m["AvatarUrl"].(string)
		//user.Captcha   = If(m["Captcha"] != nil,m["Captcha"].(int),nil).(int)
		//user.Intro     = m["Intro"].(string)
		//user.Admin     = m["Admin"].(bool)
		//user.CreatedAt = m["CreatedAt"].(time.Time)
		//user.UpdatedAt = m["UpdatedAt"].(time.Time)
	}
	fmt.Println(user)
	c.Data["json"] = Success(user)
	c.ServeJSON()
}