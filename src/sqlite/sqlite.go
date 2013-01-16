package sqlite

import (
	. "github.com/kuroneko/gosqlite3"
	"reflect"
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

func (db *DB) QueryResult(v interface{}, sql string, args ...interface{}) {
	rv := reflect.ValueOf(v)
	pv := rv
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		panic("Invalid Unmarshal Error")
	}

	st, err := db.connection.Prepare(sql)
	if err != nil {
		panic(err)
	}
	rv = rv.Addr()
	rv.Set(reflect.New(rv.Type().Elem()))
	st.All(func(st *Statement, row ...interface{}) {

	})

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
