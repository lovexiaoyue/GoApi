package controllers

import (
	"MyGoApi/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// AdsController operations for Ads
type TagsController struct {
	beego.Controller
}

// @router /list [post]
func (c *TagsController) List() {
	//var v map[string]interface{}
	//var page Paginator
	//var articles []models.Articles
	//if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
	//	//计算总量
	//	o := orm.NewOrm()
	//	num,err := o.QueryTable("articles").Count()
	//	if err != nil {
	//		c.Data["json"] = Error(err.Error())
	//	}else{
	//		// 当前页码和其实查询偏移量
	//		CurrentPage := int(v["page"].(float64))
	//		From := (CurrentPage-1)*PageSize + 1
	//		//数据查询
	//		o.QueryTable("articles").OrderBy("-created_at").Limit(PageSize,From).All(&articles)
	//
	//		// 数据返回
	//		page.Data = articles
	//		page.Total = int(num)
	//		c.Data["json"] = Success(page)
	//	}
	//}else{
	//	c.Data["json"] = Error(err.Error())
	//}
	//c.ServeJSON()
	//var ats models.Articles
	var v map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		o := orm.NewOrm()
		var posts []*models.Articles
		CurrentPage := int(v["page"].(float64))
		From := (CurrentPage-1)*PageSize
		num, _ := o.QueryTable("articles").Filter("Tags__Tag", v["tag"]).Filter("deleted_at__isnull", true).Count()
		if num == 0 {
			c.Data["json"] = Error(400,"该标签下的文章暂时下架")
			c.ServeJSON()
		}
		_, err := o.QueryTable("articles").Filter("Tags__Tag", v["tag"]).Filter("deleted_at__isnull", true).OrderBy("-created_at").Limit(PageSize,From).All(&posts)
		var res []interface{}
		if err == nil {
			for _, post := range posts {
				info := make(map[string]interface{})
				info["article"] = post
				info["article_id"] = post.Id
				res = append(res, info)
			}
		}
		page := make(map[string]interface{})
		page["total"] = num
		page["data"] = res
		c.Data["json"] = Success(page)
		c.ServeJSON()
	}
}