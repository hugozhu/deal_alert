package sqlite

import (
	"fmt"
	"github.com/hugozhu/log4go"
	. "github.com/kuroneko/gosqlite3"
	"os"
	"reflect"
	"strings"
)

type DB struct {
	connection *Database
}

var log = log4go.New(os.Stdout)

func convertRow(st *Statement, row []interface{}) map[string]interface{} {
	a := make(map[string]interface{})
	for i := 0; i < st.Columns(); i++ {
		a[st.ColumnName(i)] = row[i]
	}
	return a
}

func (db *DB) Query(v interface{}, sql string, args ...interface{}) {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x, ":", sql)
		}
	}()

	rv := reflect.ValueOf(v)
	pv := rv
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		panic("Invalid Unmarshal Error, must be pointer")
	}
	rv = rv.Elem()
	isSlice := true
	if rv.Kind() != reflect.Slice {
		rv.Set(reflect.New(rv.Type()).Elem())
		isSlice = false
	} else {
		rv.Set(reflect.MakeSlice(rv.Type(), 10, 10))
	}

	st, err := db.connection.Prepare(sql)
	if err != nil {
		panic(err)
	}
	i := 0
	st.All(func(st *Statement, row ...interface{}) {
		var r reflect.Value
		if isSlice {
			if i >= rv.Cap() {
				//grow slice if necessary
				newcap := rv.Cap() + rv.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(rv.Type(), rv.Len(), newcap)
				reflect.Copy(newv, rv)
				rv.Set(newv)
			}
			if i >= rv.Len() {
				rv.SetLen(i + 1)
			}
			r = rv.Index(i)
		} else {
			r = rv
		}
		for j := 0; j < st.Columns(); j++ {
			name := convertColumnNameToFieldName(st.ColumnName(j))
			field := r.FieldByName(name)
			if field.IsValid() {
				switch field.Kind() {
				case reflect.Int:
				case reflect.Int64:
					field.SetInt(row[j].(int64))
					break
				case reflect.String:
					field.SetString(row[j].(string))
				default:
					panic(fmt.Sprintf("Failed to convert: %+v", st.ColumnName(j)))
				}
			}
		}
		i++
	})
	if isSlice {
		rv.SetLen(i)
	}
}

func convertColumnNameToFieldName(s string) string {
	bytes1 := []byte(s)
	bytes2 := make([]byte, len(bytes1))
	length := 0
	for i := 0; i < len(bytes1); i++ {
		x := (i == 0)
		if bytes1[i] == '_' {
			i++
			x = true
		}
		if x {
			bytes2[length] = []byte(strings.ToUpper(s[i : i+1]))[0]
		} else {
			bytes2[length] = bytes1[i]
		}
		length++
	}
	return string(bytes2[0:length])
}

func (db *DB) QueryForMap(sql string, args ...interface{}) []map[string]interface{} {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x, ":", sql)
		}
	}()

	var result []map[string]interface{}
	st, err := db.connection.Prepare(sql)
	if err != nil {
		panic(err)
		return result
	}
	st.All(func(st *Statement, row ...interface{}) {
		result = append(result, convertRow(st, row))
	})
	return result
}

func (db *DB) Execute(sql string, args ...interface{}) (c int, e error) {
	var st *Statement
	if st, e = db.connection.Prepare(sql, args...); e == nil {
		c, e = st.All(func(st *Statement, row ...interface{}) {
			log.Debug("hello")
		})
	}
	return

}

func Run(file string, f func(*DB)) {
	Session(file, func(db *Database) {
		f(&DB{db})
	})
}
