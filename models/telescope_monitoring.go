package models

type TelescopeMonitoring struct {
	Tag string `orm:"column(tag);size(191)"`
}
