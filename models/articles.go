package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Articles struct {
	Id        int       `json:"id"           orm:"column(id);auto"`
	Title     string    `json:"title"        orm:"column(title);size(191)"`
	Desc      string    `json:"desc"         orm:"column(desc);null"`
	Img       string    `json:"img"          orm:"column(img);null"`
	Content   string    `json:"content"      orm:"column(content);null"`
	Clicks    int       `json:"clicks"       orm:"column(clicks);null"`
	Classify  string    `json:"classify"     orm:"column(classify);size(191);null"`
	Like      int       `json:"like"         orm:"column(like);null"`
	DeletedAt time.Time `json:"deleted_at"   orm:"column(deleted_at);type(timestamp);null"`
	CreatedAt time.Time `json:"created_at"   orm:"column(created_at);type(timestamp);null"`
	UpdatedAt time.Time `json:"updated_at"   orm:"column(updated_at);type(timestamp);null"`
	//Tags []*Tags `orm:"reverse(many)"` // fk 的反向关系
	Tags []*Tags        `json:"tags"         orm:"reverse(many)"` // fk 的反向关系
}

func (t *Articles) TableName() string {
	return "articles"
}

func init() {
	orm.RegisterModel(new(Articles))
}

// AddArticles insert a new Articles into database and returns
// last inserted Id on success.
func AddArticles(m *Articles) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetArticlesById retrieves Articles by Id. Returns error if
// Id doesn't exist
func GetArticlesById(id int) (v *Articles, err error) {
	o := orm.NewOrm()
	v = &Articles{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllArticles retrieves all Articles matches certain condition. Returns empty list if
// no records exist
func GetAllArticles(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Articles))
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

	var l []Articles
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

// UpdateArticles updates Articles by Id and returns error if
// the record to be updated doesn't exist
func UpdateArticlesById(m *Articles) (err error) {
	o := orm.NewOrm()
	v := Articles{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteArticles deletes Articles by Id and returns error if
// the record to be deleted doesn't exist
func DeleteArticles(id int) (err error) {
	o := orm.NewOrm()
	v := Articles{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Articles{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
