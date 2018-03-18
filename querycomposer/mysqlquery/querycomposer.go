package mysqlquery

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type MySqlUpdateQueryComposer struct {
	MySqlInsertQueryComposer
	whereClause string
}

func NewMySqlUpdateQueryComposer() *MySqlUpdateQueryComposer {
	return new(MySqlUpdateQueryComposer)
}

func (q *MySqlUpdateQueryComposer) Columns(input interface{}) *MySqlUpdateQueryComposer {
	columns := createColumnClauseArray(input)
	var tmp string
	for i, c := range columns {
		tmp = c + " = ?"
		columns[i] = tmp
	}

	q.columnClause = strings.Join(columns, ", ")

	return q
}

func (q *MySqlUpdateQueryComposer) Where(whereClause string) *MySqlUpdateQueryComposer {
	q.whereClause = whereClause
	return q
}

func (q *MySqlUpdateQueryComposer) PersistenceNames(names []string) *MySqlUpdateQueryComposer {
	q.persistenceClause = strings.Join(names, ",")
	return q
}

func (q *MySqlUpdateQueryComposer) Compose() string {
	query := bytes.NewBufferString("UPDATE ")
	query.WriteString(q.persistenceClause)
	query.WriteString(" SET " + q.columnClause + " ")
	query.WriteString("WHERE " + q.whereClause)

	return query.String()
}

type MySqlInsertQueryComposer struct {
	columnClause      string
	persistenceClause string
}

func NewMySqlInsertQueryComposer() *MySqlInsertQueryComposer {
	return new(MySqlInsertQueryComposer)
}

func (q *MySqlInsertQueryComposer) Columns(input interface{}) *MySqlInsertQueryComposer {

	q.columnClause = createColumnClause(input)

	return q
}

func(q *MySqlInsertQueryComposer) PersistenceNames(names []string) *MySqlInsertQueryComposer {
	if len(names) == 0 || len(names) != 1 {
		return q
	}
	q.persistenceClause = names[0]
	return q
}

func (q *MySqlInsertQueryComposer) Compose() string {
	query := bytes.NewBufferString("INSERT INTO ")
	query.WriteString(q.persistenceClause)
	query.WriteString("( " + q.columnClause + " ) ")
	query.WriteString("VALUES ( ")
	ssplit := strings.Split(q.columnClause, ",")
	for i:=range ssplit {
		if i == 0 {
			query.WriteString("?")
		} else {
			query.WriteString(",?")
		}
	}
	query.WriteString(" )")

	return query.String()
}

type MySqlSelectQueryComposer struct {
	columnClause      string
	whereClause       string
	orderClause       string
	pagingClause      string
	persistenceClause string
}

func NewMySqlSelectQueryComposer() *MySqlSelectQueryComposer {
	return new(MySqlSelectQueryComposer)
}

func (q *MySqlSelectQueryComposer) Where(whereClause string) *MySqlSelectQueryComposer {
	q.whereClause = fmt.Sprintf("WHERE %s", whereClause)
	return q
}

func (q *MySqlSelectQueryComposer) OrderBy(orderByClause string) *MySqlSelectQueryComposer {
	q.orderClause = "ORDER BY " + orderByClause
	return q
}

func (q *MySqlSelectQueryComposer) Paginate(page, limit int) *MySqlSelectQueryComposer {
	q.pagingClause = fmt.Sprintf("LIMIT %d OFFSET %d", limit, (page-1)*limit)
	return q
}

// PersistenceNames creates the table and inner join part of the query
// names contain the name of table, separated by coma
func (q *MySqlSelectQueryComposer) PersistenceNames(names []string) *MySqlSelectQueryComposer {

	if len(names) == 0 {
		return q
	}

	sfrom := fmt.Sprintf("FROM %s", strings.Join(names, ", "))

	q.persistenceClause = sfrom

	return q
}

// Columns creates/composes column clause.
// input only three types are accepted for input, those are []string, string, or simple struct
// simple struct mean, none of field have type struct, slice, or map
// Columns only read field that defined by tag
func (q *MySqlSelectQueryComposer) Columns(input interface{}) *MySqlSelectQueryComposer {

	switch input.(type) {
	case []string:
		sstr := input.([]string)
		q.columnClause = "SELECT " + strings.Join(sstr, ", ")
	case string:
		q.columnClause = "SELECT " + input.(string)
	default:

		q.columnClause = "SELECT " + createColumnClause(input)
	}

	return q
}

func (q *MySqlSelectQueryComposer) Compose() string {
	query := bytes.NewBufferString("")
	query.WriteString(q.columnClause + " ")
	query.WriteString(q.persistenceClause + " ")
	query.WriteString(q.whereClause + " ")
	query.WriteString(q.orderClause + " ")
	query.WriteString(q.pagingClause)

	return query.String()
}

func createColumnClauseArray(input interface{}) []string {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()
	if ivalue.Kind() != reflect.Struct {
		return nil
	}

	clauses := make([]string, 0)

	for i := 0; i < ivalue.NumField(); i++ {
		ftype := itype.Field(i)

		if ftype.Type.Kind() == reflect.Struct || ftype.Type.Kind() == reflect.Slice ||
			ftype.Type.Kind() == reflect.Map {
			continue // unsupported type
		}

		tags := ftype.Tag.Get("query")
		tsplit := strings.Split(tags, ",")

		if len(tsplit) > 0 {
			var skip = len(tsplit) > 1 && tsplit[1] == "primary"
			if !skip {
				clauses = append(clauses, tsplit[0])
			}

		}
	}

	return clauses
}

func createColumnClause(input interface{}) string {
	clauses := createColumnClauseArray(input)

	return strings.Join(clauses, ", ")
}
