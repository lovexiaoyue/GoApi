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
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/ads",
			beego.NSInclude(
				&controllers.AdsController{},
			),
		),

		beego.NSNamespace("/articles",
			beego.NSInclude(
				&controllers.ArticlesController{},
			),
		),

		beego.NSNamespace("/comments",
			beego.NSInclude(
				&controllers.CommentsController{},
			),
		),

		beego.NSNamespace("/links",
			beego.NSInclude(
				&controllers.LinksController{},
			),
		),

		beego.NSNamespace("/messages",
			beego.NSInclude(
				&controllers.MessagesController{},
			),
		),

		beego.NSNamespace("/migrations",
			beego.NSInclude(
				&controllers.MigrationsController{},
			),
		),

		beego.NSNamespace("/telescope_entries",
			beego.NSInclude(
				&controllers.TelescopeEntriesController{},
			),
		),

		beego.NSNamespace("/user_auths",
			beego.NSInclude(
				&controllers.UserAuthsController{},
			),
		),

		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UsersController{},
			),
		),

		beego.NSNamespace("/webinfos",
			beego.NSInclude(
				&controllers.WebinfosController{},
			),
		),
	)

	beego.Router("/v1/login",&controllers.UsersController{},"post:Login")
	beego.Router("/v1/register",&controllers.UsersController{},"post:Register")
	beego.AddNamespace(ns)
}
