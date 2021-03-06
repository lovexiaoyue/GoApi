package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Ads struct {
	Id        int       `json:"id"              orm:"column(id);auto"`
	Title     string    `json:"title"           orm:"column(title);size(191)"`
	Url       string    `json:"url"             orm:"column(url);size(191)"`
	Type      string    `json:"type"            orm:"column(type);size(191);null"`
	CreatedAt time.Time `json:"created_at"      orm:"column(created_at);type(timestamp);null"`
	UpdatedAt time.Time `json:"updated_at"      orm:"column(updated_at);type(timestamp);null"`
}

func (t *Ads) TableName() string {
	return "ads"
}

func init() {
	orm.RegisterModel(new(Ads))
}

// AddAds insert a new Ads into database and returns
// last inserted Id on success.
func AddAds(m *Ads) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAdsById retrieves Ads by Id. Returns error if
// Id doesn't exist
func GetAdsById(id int) (v *Ads, err error) {
	o := orm.NewOrm()
	v = &Ads{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAds retrieves all Ads matches certain condition. Returns empty list if
// no records exist
func GetAllAds(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Ads))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Ads
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateAds updates Ads by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdsById(m *Ads) (err error) {
	o := orm.NewOrm()
	v := Ads{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAds deletes Ads by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAds(id int) (err error) {
	o := orm.NewOrm()
	v := Ads{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Ads{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func ListAds(page,count int) (ads []Ads,err error) {
	o := orm.NewOrm()
	beego.Info(count)
	var v []Ads
	qs := o.QueryTable("ads")
	qs.Limit(count).Offset((page-1)*count).All(&v)
	beego.Info(ads)
	return v,nil
}