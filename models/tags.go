package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Tags struct {
	Id        int       `json:"id"               orm:"column(id);auto"`
	Tag       string    `json:"tag"              orm:"column(tag);size(191)"`
	Classify  string    `json:"classify"         orm:"column(classify);size(191);null"`
	Article *Articles `json:"article_id"         orm:"rel(fk)"`
	DeletedAt time.Time `json:"deleted_at"       orm:"column(deleted_at);type(timestamp);null"`
}

func (t *Tags) TableName() string {
	return "tags"
}

func init() {
	orm.RegisterModel(new(Tags))
}