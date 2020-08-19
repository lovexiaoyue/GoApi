package controllers

import (
	"MyGoApi/models"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// ArticlesController operations for Articles
type ArticlesController struct {
	beego.Controller
}

var PageSize int = 10

var BasePath string = "http://127.0.0.1:8080/v1/article/list"
// URLMapping ...
func (c *ArticlesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Articles
// @Param	body		body 	models.Articles	true		"body for Articles content"
// @Success 201 {int} models.Articles
// @Failure 403 body is empty
// @router / [post]
func (c *ArticlesController) Post() {
	var v models.Articles
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddArticles(&v); err == nil {
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
// @Description get Articles by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Articles
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ArticlesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetArticlesById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Articles
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Articles
// @Failure 403
// @router / [get]
func (c *ArticlesController) GetAll() {
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

	l, err := models.GetAllArticles(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Articles
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Articles	true		"body for Articles content"
// @Success 200 {object} models.Articles
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ArticlesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Articles{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateArticlesById(&v); err == nil {
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
// @Description delete the Articles
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ArticlesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteArticles(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}


// @router /list [post]
func (c *ArticlesController) List() {
	var v map[string]interface{}
	var page Paginator
	var articles []models.Articles
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//page := v["page"].(float64)
		o := orm.NewOrm()
		num,err := o.QueryTable("articles").Count()
		if err != nil {
			c.Data["json"] = Error(err.Error())
		}else{
			CurrentPage := int(v["page"].(float64))
			From := (CurrentPage-1)*PageSize + 1
			To := 10
			if (CurrentPage-1)*PageSize <= int(num){
				To   = (CurrentPage-1)*PageSize
			}else{
				To   = int(num)
			}
			LastPage:= int(math.Ceil((float64(num)/float64(PageSize))))
			o.QueryTable("articles").OrderBy("-created_at").Limit(PageSize,From).All(&articles)

			page.CurrentPage = CurrentPage
			page.Data = articles
			page.FirstPageUrl = BasePath+"?page=1"
			page.LastPage = LastPage
			page.LastPageUrl = BasePath+"?page="+string(LastPage)
			page.NextPageUrl = BasePath+"?page="+string(CurrentPage+1)
			page.Path = BasePath
			page.PerPage = 10
			if CurrentPage-1 <= 0{
				page.PrevPageUrl = ""
			}else{
				page.PrevPageUrl = BasePath+"?page="+string(CurrentPage+1)
			}
			page.Total = int(num)
			page.From = From
			page.To = To
			c.Data["json"] = Success(page)
		}
	}else{
		c.Data["json"] = Error(err.Error())
	}
	c.ServeJSON()
}
