package driver

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"stage/app/common"
	"strconv"
	"strings"
	"unsafe"
)

type Query struct {
	db    **gorm.DB
	data  map[string]interface{}
	limit uint
	field string
	key   string
	num   uint16
}

func NewQuery(db **gorm.DB, data map[string]interface{}) Query {
	var limit uint
	if common.MakeUint(data["limit"]) != 0 {
		limit = common.MakeUint(data["limit"])
	} else {
		limit = 20
	}
	query := Query{db: db, data: data, limit: limit}
	if order, ok := query.data["order"].(string); ok {
		if order != "" && order != "null" {
			arr := strings.Split(order, " ")
			columns := strings.Split(arr[0], ",")
			for k, v := range columns {
				//如果不是json类型排序，则添加反引号
				if !strings.Contains(v, "->") {
					columns[k] = "`" + v + "`"
				}
			}
			*query.db = (*query.db).Order(common.Implode(columns, ",") + " " + common.Implode(arr[1:], " "))
		}
	} else {
		*query.db = (*query.db).Order("id desc")
	}
	return query
}

func (q *Query) Eq(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" = ?", q.data[q.key])
	}
}

func (q *Query) Gt(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" > ?", q.data[q.key])
	}
}

func (q *Query) Gte(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" >= ?", q.data[q.key])
	}
}

func (q *Query) Lt(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" < ?", q.data[q.key])
	}
}

func (q *Query) Lte(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" <= ?", q.data[q.key])
	}
}

func (q *Query) Like(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" like ?", "%"+q.data[q.key].(string)+"%")
	}
}

// AwLike after wildcard后通配
func (q *Query) AwLike(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" like ?", q.data[q.key].(string)+"%")
	}
}

func (q *Query) In(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" in (?)", q.data[q.key])
	}
}

func (q *Query) NotIn(args ...string) {
	if q.checkKey(args...) {
		*q.db = (*q.db).Where(q.field+" not in (?)", q.data[q.key])
	}
}

func (q *Query) EqZero(args ...string) {
	q.checkKey(args...)
	*q.db = (*q.db).Where(q.field + " = 0")
}

func (q *Query) GtZero(args ...string) {
	q.checkKey(args...)
	*q.db = (*q.db).Where(q.field + " > 0")
}

func (q *Query) IsEmpty(args ...string) {
	if q.checkKey(args...) {
		parseBool, _ := strconv.ParseBool(common.MakeString(q.data[q.key]))
		if parseBool {
			*q.db = (*q.db).Where(q.field + " <> ''")
		} else {
			*q.db = (*q.db).Where(q.field + " is null or " + q.field + " = ''")
		}
	}
}

func (q *Query) Null(args ...string) {
	q.checkKey(args...)
	*q.db = (*q.db).Where(q.field + " is null")
}

func (q *Query) NotNull(args ...string) {
	q.checkKey(args...)
	*q.db = (*q.db).Where(q.field + " is not null")
}

//wildcard 通配
func (q *Query) Wc(args ...string) {
	if q.checkKey(args...) {
		if q.data[q.key] == "*" {
			*q.db = (*q.db).Where(q.field + " > 0")
		} else {
			*q.db = (*q.db).Where(q.field+" = ?", q.data[q.key])
		}
	}
}

//原生where语句
func (q *Query) Raw(query interface{}, args ...interface{}) {
	var ok bool
	var value interface{}
	var values []interface{}
	if len(args) > 0 {
		for _, v := range args {
			value, ok = q.data[v.(string)]
			values = append(values, value)
		}
		if ok {
			*q.db = (*q.db).Where(query, values...)
		}
	}
}

func (q *Query) Pages(value interface{}) *Query {
	t := reflect.TypeOf(value).Elem()
	v := reflect.ValueOf(value)
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		switch name {
		case "Limit":
			*(*uint)(unsafe.Pointer(v.Elem().FieldByName(name).Addr().Pointer())) = q.limit
		case "Page":
			*(*uint)(unsafe.Pointer(v.Elem().FieldByName(name).Addr().Pointer())) = common.MakeUint(q.data["page"])
		case "Count":
			var count uint
			(*q.db).Count(&count)
			*(*uint)(unsafe.Pointer(v.Elem().FieldByName(name).Addr().Pointer())) = count
		}
	}
	return q
}

func (q *Query) List(value interface{}) *Query {
	(*q.db).Limit(q.limit).Offset(common.GetOffset(q.data["page"], q.limit)).Find(value)
	return q
}

func (q *Query) checkKey(args ...string) bool {
	arr := strings.Split(args[0], ".")
	if len(args) == 1 {
		if len(arr) == 1 {
			q.field = "`" + args[0] + "`"
			q.key = arr[0]
		} else {
			q.field = args[0]
			q.key = arr[1]
		}
	} else {
		if len(arr) == 1 {
			q.field = "`" + args[0] + "`"
		} else {
			q.field = args[0]
		}
		q.key = args[1]
	}
	if q.data[q.key] == nil {
		return false
	}
	q.num++
	return true
}
