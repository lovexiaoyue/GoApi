package controllers

import "github.com/astaxie/beego"

type VerifyController struct {
	beego.Controller
}

// 登录校验
type LoginVerify struct {
	Name      string
	Password  string
}


// 注册校验
type RegisterVerify struct {
	Name        string
	Email       string
	Password    string
	Repassword  string
}