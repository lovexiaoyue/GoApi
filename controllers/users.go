package controllers

import (
	"MyGoApi/models"
	"MyGoApi/utils"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
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
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddUsers(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()

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
		c.Data["json"] = Error(400,err.Error())
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
	var v utils.User
	var token Token

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &data); err == nil {
		if err := utils.CheckLogin(data.Name,data.Password); err != "ok"{
			c.Data["json"] = Error(400,err)
		}else{
			user,_ := models.GetUsersByName(data.Name)
			if user == nil {
				c.Data["json"] = Error(400,"用户名或密码错误")
				c.ServeJSON()
			}
			secret,_ := GeneratePassWd(data.Password)
			beego.Info(string(secret))
			beego.Info("===============")
			beego.Info(user.Password)
			ok,_:=ValidatePassWd(data.Password,user.Password)
			beego.Info(ok)
			if !ok{
				c.Data["json"] = Error(400,"用户名或密码错误")
				c.ServeJSON()
			}

			//v = new(utils.User)
			v.Id = user.Id
			v.Name = user.Name
			token.Token = utils.GenerateToken(0,v)
			c.SetSession("name",user.Name)
			c.SetSession("admin",user.IsAdmin)
			beego.Info(token)
			c.Data["json"] = Success(token)

		}
	} else {
		c.Data["json"] = Error(400,err.Error())
	}

	c.ServeJSON()
}

// @router /register [post]
func (c *UsersController) Register() {

	var data RegisterVerify
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &data); err == nil {
		if err := utils.CheckRegister(data.Name,data.Email,data.Password,data.Repassword); err != "ok"{
			c.Data["json"] = Error(400,err)
		}else{
			user := new(models.Users)
			user.Name = data.Name
			user.Email = data.Email
			user.Phone = ""
			secret,_ := GeneratePassWd(data.Password)
			user.Password = string(secret)
			o := orm.NewOrm()
			_,err := o.Insert(user)
			if err != nil {
				beego.Info(err)
				c.Data["json"] = Error(400,"用户名已存在")
			}else{
				c.Data["json"] = Success("注册成功")
			}
		}
	} else {
		c.Data["json"] = Error(400,err.Error())
	}
	c.ServeJSON()
}

// @router /info [get]
func (c *UsersController) UserInfo() {

	//用户校验 获取用户信息
	//token := c.Ctx.Input.Header("Authorization")
	//v,err := utils.ValidateToken(token)
	name := c.GetSession("name")
	if name == nil {
		c.Data["json"] = Error(422,"未登录或登录状态失效")
		c.ServeJSON()
	}
	beego.Info("sessionName:",name)
	//if err != nil {
	//	c.Data["json"] = Error("token err")
	//	c.ServeJSON()
	//}

	// 查询用户信息
	//u,_ := models.GetUsersByName(name.(string))
	var u models.Users
	orm := orm.NewOrm()
	err := orm.QueryTable("users").Filter("name",name).One(&u,"id","name","email","phone","avatar_url","intro","is_admin","created_at","updated_at")
	if err != nil {
		c.Data["json"] = Error(400,"查询出错")
		c.ServeJSON()
	}

	// 判断是否是admin
	var admin bool
	if u.IsAdmin == "1" {
		admin = true
	}else{
		admin = false
	}

	// 定义返回字段
	user := make(map[string]interface{})
	user["id"]          = u.Id
	user["name"]        = u.Name
	user["phone"]       = u.Phone
	user["email"]       = u.Email
	user["avatar_url"]  = u.AvatarUrl
	user["captcha"]     = u.Captcha
	user["intro"]       = u.Intro
	user["admin"]       = admin
	user["created_at"]  = u.CreatedAt
	user["updated_at"]  = u.UpdatedAt

	c.Data["json"] = Success(user)
	c.ServeJSON()
}

// 用户退出接口
// @router /logout [post]
func (c *UsersController) Logout() {

	c.DelSession("name")
	c.DelSession("admin")
	c.Data["json"] = Success("退出成功")
	c.ServeJSON()

}