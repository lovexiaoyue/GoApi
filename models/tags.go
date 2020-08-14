package models

import "time"

type Tags struct {
	Tag       string    `json:"tag"              orm:"column(tag);size(191)"`
	Classify  string    `json:"classify"         orm:"column(classify);size(191);null"`
	ArticleId *Articles `json:"article_id"       orm:"column(article_id);rel(fk)"`
	DeletedAt time.Time `json:"deleted_at"       orm:"column(deleted_at);type(timestamp);null"`
}
