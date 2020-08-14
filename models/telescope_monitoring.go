package models

type TelescopeMonitoring struct {
	Tag string `json:"tag"    orm:"column(tag);size(191)"`
}
