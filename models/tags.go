package models

import "time"

type Tags struct {
	Tag       string    `orm:"column(tag);size(191)"`
	Classify  string    `orm:"column(classify);size(191);null"`
	ArticleId *Articles `orm:"column(article_id);rel(fk)"`
	DeletedAt time.Time `orm:"column(deleted_at);type(timestamp);null"`
}
