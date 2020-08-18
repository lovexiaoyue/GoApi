package utils

import (
	"fmt"
	"github.com/astaxie/beego/validation"
)


// 校验登录数据
func CheckLogin(name,password string)(errMessage string)  {

	valid := validation.Validation{}
	valid.Required(name, "name").Message("用户名必填")
	valid.Required(password, "Password").Message("密码必填")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		for _, err := range valid.Errors {
			return fmt.Sprintf("%s", err.Message)
		}
	}
	return fmt.Sprintf("%s", "ok")
}

// 校验登录数据
func CheckList(page,count int)(errMessage string)  {

	valid := validation.Validation{}
	valid.Required(page, "page").Message("页码")
	valid.Required(count, "sumCount").Message("分页数量")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		for _, err := range valid.Errors {
			return fmt.Sprintf("%s", err.Message)
		}
	}
	return fmt.Sprintf("%s", "ok")
}

// 校验注册数据
func CheckRegister(name,email,password,repasswod string)(errMessage string)  {

	valid := validation.Validation{}
	valid.Required(name,  "name").Message("用户名必填")
	valid.Required(email,     "Email").Message("邮箱必填")
	valid.Required(password,  "Password").Message("密码必填")
	valid.Required(repasswod, "Repassword").Message("确认密码必填")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		for _, err := range valid.Errors {
			return fmt.Sprintf("%s", err.Message)
		}
	}
	return fmt.Sprintf("%s", "ok")
}