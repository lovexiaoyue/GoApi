package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

type Response struct {
	Status int `json:"status"`
	ErrorCode int `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Data interface{} `json:"data"`
}

// 返回成功结果
func Success(data interface{}) Response {
	return Response{200,200,"",data}
}

// 返回失败结果
func Error(code int, message string) Response {
	return Response{200,code,message,nil}
}