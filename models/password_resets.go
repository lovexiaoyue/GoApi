package models

import "time"

type PasswordResets struct {
	Email     string    `json:"email"        orm:"column(email);size(191)"`
	Token     string    `json:"token"        orm:"column(token);size(191)"`
	CreatedAt time.Time `json:"created_at"   orm:"column(created_at);type(timestamp);null"`
}
