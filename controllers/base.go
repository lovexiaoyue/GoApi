package controllers

import (
	"MyGoApi/utils"
	"github.com/astaxie/beego"
)

type Token struct {
	Token string   `json:"token"`
}

type BaseController struct {
	beego.Controller
}

type Response struct {
	Code int               `json:"code"`
	Status string          `json:"status"`
	Data interface{}       `json:"data"`
}

// 返回成功结果
func Success(data interface{}) Response {
	return Response{200,"success",data}
}

// 返回失败结果
func Error(message string) Response {
	return Response{200,"fail",message}
}

// @router / [post]
func (c *BaseController) RefreshToken (){
	if c.Ctx.Input.Header("Authorization") == "" {
		c.Data["json"] = Error("没有携带token")
	}
	token := c.Ctx.Input.Header("Authorization")
	token, err := utils.RefreshToken(token)
	retoken  := Token{token}

	if err != nil {
		c.Data["json"] = Error("刷新token失败")
	}else{
		c.Data["json"] = Success(retoken)
	}
}

