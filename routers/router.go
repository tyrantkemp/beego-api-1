// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego-api-1/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),

		beego.NSNamespace("/profile",
			beego.NSInclude(
				&controllers.ProfileController{},
			),
		),

		beego.NSNamespace("/resource",
			beego.NSInclude(
				&controllers.ResourceController{},
			),
		),

		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),

		beego.NSNamespace("/role_resource_rel",
			beego.NSInclude(
				&controllers.RoleResourceRelController{},
			),
		),

		beego.NSNamespace("/role_user_rel",
			beego.NSInclude(
				&controllers.RoleUserRelController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
