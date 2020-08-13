package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Webinfos struct {
	Id          int       `orm:"column(id);auto"`
	Title       string    `orm:"column(title);size(191);null"`
	Keyword     string    `orm:"column(keyword);size(191);null"`
	Description string    `orm:"column(description);null"`
	Personinfo  string    `orm:"column(personinfo);null"`
	Github      string    `orm:"column(github);size(191);null"`
	Icp         string    `orm:"column(icp);size(191);null"`
	Weixin      string    `orm:"column(weixin);size(191);null"`
	Zhifubao    string    `orm:"column(zhifubao);size(191);null"`
	Qq          string    `orm:"column(qq);size(191);null"`
	Phone       string    `orm:"column(phone);size(191);null"`
	Email       string    `orm:"column(email);size(191);null"`
	StartTime   time.Time `orm:"column(startTime);type(date);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *Webinfos) TableName() string {
	return "webinfos"
}

func init() {
	orm.RegisterModel(new(Webinfos))
}

// AddWebinfos insert a new Webinfos into database and returns
// last inserted Id on success.
func AddWebinfos(m *Webinfos) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetWebinfosById retrieves Webinfos by Id. Returns error if
// Id doesn't exist
func GetWebinfosById(id int) (v *Webinfos, err error) {
	o := orm.NewOrm()
	v = &Webinfos{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWebinfos retrieves all Webinfos matches certain condition. Returns empty list if
// no records exist
func GetAllWebinfos(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Webinfos))
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

	var l []Webinfos
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

// UpdateWebinfos updates Webinfos by Id and returns error if
// the record to be updated doesn't exist
func UpdateWebinfosById(m *Webinfos) (err error) {
	o := orm.NewOrm()
	v := Webinfos{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWebinfos deletes Webinfos by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWebinfos(id int) (err error) {
	o := orm.NewOrm()
	v := Webinfos{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Webinfos{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
