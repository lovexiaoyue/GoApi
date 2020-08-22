// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"MyGoApi/controllers"
	"github.com/astaxie/beego"
)




func init() {
	beego.Router("/v1/login",&controllers.UsersController{},"post:Login")
	beego.Router("/v1/register",&controllers.UsersController{},"post:Register")
	beego.Router("/v1/refresh",&controllers.BaseController{},"post:RefreshToken")
	beego.Router("/v1/list",&controllers.AdsController{},"post:List")
	beego.Router("/v1/list",&controllers.ArticlesController{},"post:List")
	beego.Router("/v1/classify",&controllers.ArticlesController{},"get:Classify")
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/ad",
			beego.NSInclude(
				&controllers.AdsController{},
			),
		),

		beego.NSNamespace("/article",
			beego.NSInclude(
				&controllers.ArticlesController{},
			),
		),

		beego.NSNamespace("/comment",
			beego.NSInclude(
				&controllers.CommentsController{},
			),
		),

		beego.NSNamespace("/link",
			beego.NSInclude(
				&controllers.LinksController{},
			),
		),

		beego.NSNamespace("/message",
			beego.NSInclude(
				&controllers.MessagesController{},
			),
		),

		beego.NSNamespace("/migration",
			beego.NSInclude(
				&controllers.MigrationsController{},
			),
		),

		beego.NSNamespace("/telescope_entrie",
			beego.NSInclude(
				&controllers.TelescopeEntriesController{},
			),
		),

		beego.NSNamespace("/user_auth",
			beego.NSInclude(
				&controllers.UserAuthsController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UsersController{},
			),
		),

		beego.NSNamespace("/webinfo",
			beego.NSInclude(
				&controllers.WebinfosController{},
			),
		),
		beego.NSNamespace("/bese",
			beego.NSInclude(
				&controllers.BaseController{},
			),
		),
	)
	//beego.InsertFilter("/*",beego.BeforeExec,utils.FilterToken)

	beego.AddNamespace(ns)
}
