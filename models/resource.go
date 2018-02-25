package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"beego-api-1/utils"
)

type Resource struct {
	Id       int    `orm:"column(id);auto"`
	Rtype    int    `orm:"column(rtype)"`
	Name     string `orm:"column(name);size(64)"`
	ParentId int    `orm:"column(parent_id);null"`
	Seq      int    `orm:"column(seq)"`
	Icon     string `orm:"column(icon);size(32)"`
	UrlFor   string `orm:"column(url_for);size(256)"`
}

func (t *Resource) TableName() string {
	return "resource"
}

func init() {
	orm.RegisterModel(new(Resource))
}

// AddResource insert a new Resource into database and returns
// last inserted Id on success.
func AddResource(m *Resource) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetResourceById retrieves Resource by Id. Returns error if
// Id doesn't exist
func GetResourceById(id int) (v *Resource, err error) {
	o := orm.NewOrm()
	v = &Resource{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllResource retrieves all Resource matches certain condition. Returns empty list if
// no records exist
func GetAllResource(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Resource))
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

	var l []Resource
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

// UpdateResource updates Resource by Id and returns error if
// the record to be updated doesn't exist
func UpdateResourceById(m *Resource) (err error) {
	o := orm.NewOrm()
	v := Resource{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteResource deletes Resource by Id and returns error if
// the record to be deleted doesn't exist
func DeleteResource(id int) (err error) {
	o := orm.NewOrm()
	v := Resource{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Resource{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 根据用户id获取包含的权限url
func GetResouceUrlByUserId(userId int) []*Resource {
	resources := make([]*Resource, 0)

	key := fmt.Sprintf("test_userresourcese_urls:%v", userId)
	if err := utils.GetCache(key, resources); err == nil {
		return resources
	}

	user, err := GetUserById(userId)
	if err != nil || user == nil {
		return resources
	}
	o := orm.NewOrm()

	if user.IsAdmin == true {
		sql := fmt.Sprintf("SELECT  * from resource ORDER BY seq ASC,id asc ")
		o.Raw(sql).QueryRows(resources)
	} else {
		sql := fmt.Sprintf("select T0.* from %s as T0 INNER JOIN %s as T1 ON T0.id=T1.resource_id INNER JOIN %s as T2 ON T1.role_id=T2.role_id WHERE T2.backend_user_id=%v "+
			"ORDER BY T0.seq ASC ,T0.id ASC","resource", "role_resource_rel", "role_user_rel", userId)
		o.Raw(sql).QueryRows(resources)
	}

	// 保存进redis
	utils.SetCache(key, resources, 30)
	return resources

}
