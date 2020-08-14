package controllers

import "github.com/astaxie/beego"

type VerifyController struct {
	beego.Controller
}

type LoginVerify struct {
	username  string
	password  string
}