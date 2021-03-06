package controllers

import (
	"MyGoApi/models"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
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
	info := make(map[string]interface{})
	o := orm.NewOrm()
	idStr := c.Ctx.Input.Param(":id")
	num,err := o.QueryTable("comments").Filter("article_id",idStr).Count()
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetArticlesById(id)
	if err == nil {
		info["classify"] = v.Classify
		info["clicks"]   = v.Clicks
		info["comment"]  = num
		info["content"]  = v.Content
		info["created_at"] = v.CreatedAt
		info["deleted_at"] = v.DeletedAt
		info["desc"] = v.Desc
		info["id"] = v.Id
		info["img"] = v.Img
		info["like"] = v.Like
		info["title"] = v.Title
	}
	nv, err1 := models.GetNextArticlesById(id)
	if err1 == nil{
		nextAticle := make(map[string]interface{})
		nextAticle["id"] = nv.Id
		nextAticle["title"] = nv.Title
		info["nextrAticle"] = nextAticle
	}

	pv, err2 := models.GetPreArticlesById(id)
	if err2 == nil {
		prevAticle := make(map[string]interface{})
		prevAticle["id"] = pv.Id
		prevAticle["title"] = pv.Title
		info["prevArticle"] = prevAticle
	}

	if err != nil {
		c.Data["json"] = Error(400,"该文章已下架或删除了！！")
	} else {
		c.Data["json"] = Success(info)
	}
	//o := orm.NewOrm()
	//o.QueryTable("tags").Filter("classify")
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
		//计算总量
		o := orm.NewOrm()
		num,err := o.QueryTable("articles").Filter("deleted_at__isnull", true).Count()
		if err != nil {
			c.Data["json"] = Error(400,err.Error())
		}else{
			// 当前页码和其实查询偏移量
			CurrentPage := int(v["page"].(float64))
			From := (CurrentPage-1)*PageSize
			//数据查询
			o.QueryTable("articles").Filter("deleted_at__isnull", true).OrderBy("-created_at").Limit(PageSize,From).All(&articles)

			// 数据返回
			page.Data = articles
			page.Total = int(num)
			c.Data["json"] = Success(page)
		}
	}else{
		c.Data["json"] = Error(400,err.Error())
	}
	c.ServeJSON()
}


// @router /classify [get]
func (c *ArticlesController) Classify() {
	//var data interface{}
	o:=orm.NewOrm()
	var ats []*models.Articles
	var tags []*models.Tags
	qs := o.QueryTable("articles")
	n, err := qs.GroupBy("classify").All(&ats)

	if err == nil && n > 0 {
		var res []interface{}
		for i := 0; i < len(ats); i++ {
			// 每次都申请新地址 防止数据覆盖
			classify := make(map[string]interface{})
			o.QueryTable("tags").Filter("classify",ats[i].Classify).GroupBy("tag").All(&tags);
			classify["name"] = ats[i].Classify
			var tag []interface{}
			for i := 0; i < len(tags); i++ {
				tag = append(tag, tags[i].Tag)
				beego.Info(tags[i].Tag)
			}
			classify["tags"] = tag
			res = append(res, classify)
		}
		c.Data["json"] = Success(res)
	} else {
		beego.Info(err.Error())
		c.Data["json"] = Error(400,err.Error())
	}
	c.ServeJSON()
}