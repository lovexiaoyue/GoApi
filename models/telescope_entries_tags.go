package models

type TelescopeEntriesTags struct {
	EntryUuid *TelescopeEntries `json:"entry_uuid"    orm:"column(entry_uuid);rel(fk)"`
	Tag       string            `json:"tag"           orm:"column(tag);size(191)"`
}
