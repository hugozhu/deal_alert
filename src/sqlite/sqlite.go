package sqlite

import (
	"fmt"
	. "github.com/kuroneko/gosqlite3"
	"log"
	"reflect"
	"strings"
)

type DB struct {
	connection *Database
}

func convertRow(st *Statement, row []interface{}) map[string]interface{} {
	a := make(map[string]interface{})
	for i := 0; i < st.Columns(); i++ {
		a[st.ColumnName(i)] = row[i]
	}
	return a
}

func (db *DB) QueryResults(v interface{}, sql string, args ...interface{}) {
	rv := reflect.ValueOf(v)
	pv := rv
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		panic("Invalid Unmarshal Error, must be pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Slice {
		panic("Invalid Unmarshal Error, must be pointer of slice")
	}

	st, err := db.connection.Prepare(sql)
	if err != nil {
		panic(err)
	}
	rv.Set(reflect.MakeSlice(rv.Type(), 10, 10))
	i := 0
	st.All(func(st *Statement, row ...interface{}) {
		if i >= rv.Cap() {
			//grow slice if necessary
			newcap := rv.Cap() + rv.Cap()/2
			if newcap < 4 {
				newcap = 4
			}
			log.Println(i, rv.Cap(), newcap)
			newv := reflect.MakeSlice(rv.Type(), rv.Len(), newcap)
			reflect.Copy(newv, rv)
			rv.Set(newv)
		}
		if i >= rv.Len() {
			rv.SetLen(i + 1)
		}
		r := rv.Index(i)
		for j := 0; j < st.Columns(); j++ {
			name := convertColumnNameToFieldName(st.ColumnName(j))
			field := r.FieldByName(name)
			if field.IsValid() {
				switch field.Kind() {
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
	rv.SetLen(i)
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

func (db *DB) Query(sql string, args ...interface{}) []map[string]interface{} {
	// defer func() {
	// 	if x := recover(); x != nil {
	// 		log.Error(x, ":", sql)
	// 	}
	// }()

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

func (db *DB) Execute(sql string) {
	db.connection.Execute(sql)
}

func Run(file string, f func(*DB)) {
	Session(file, func(db *Database) {
		f(&DB{db})
	})
}
