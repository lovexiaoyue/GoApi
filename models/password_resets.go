package models

import "time"

type PasswordResets struct {
	Email     string    `orm:"column(email);size(191)"`
	Token     string    `orm:"column(token);size(191)"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null"`
}
