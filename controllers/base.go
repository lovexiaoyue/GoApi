package controllers

import (
	"MyGoApi/models"
	"MyGoApi/utils"
	"errors"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/bcrypt"
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
	Message string         `json:"message"`
	Data interface{}       `json:"data"`
}

// 分页数据


type Paginator struct {
	//CurrentPage int          `json:"current_page"`
	Data []models.Articles     `json:"data"`
	//FirstPageUrl string      `json:"first_page_url"`
	//LastPage int             `json:"last_page"`
	//LastPageUrl  string      `json:"last_page_url"`
	//NextPageUrl  string      `json:"next_page_url"`
	//Path         string      `json:"path"`
	//PerPage int              `json:"per_page"`
	//PrevPageUrl string       `json:"prev_page_url"`
	Total  int                 `json:"total"`
	//From   int               `json:"from"`
	//To     int               `json:"to"`
}
// 返回成功结果
func Success(data interface{}) Response {
	return Response{200,"success","",data}
}

// 返回失败结果
func Error(message string) Response {
	return Response{200,"fail",message,""}
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
	beego.Error(err)
	c.ServeJSON()
}

// 比对密码
func ValidatePassWd(src string, passWd string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(passWd), []byte(src)); err != nil {
		return false, errors.New("密码错误")
	}
	return true, nil
}

//生成密码
func GeneratePassWd(src string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
}