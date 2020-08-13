package models

type TelescopeEntriesTags struct {
	EntryUuid *TelescopeEntries `orm:"column(entry_uuid);rel(fk)"`
	Tag       string            `orm:"column(tag);size(191)"`
}
