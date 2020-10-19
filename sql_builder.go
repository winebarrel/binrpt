package binrpt

import (
	"strings"

	"github.com/siddontang/go-mysql/replication"
)

type SQLBuilder struct {
	Header *replication.EventHeader
	Event  *replication.RowsEvent
	Table  *Table
}

type SQL struct {
	Statement string
	Params    []interface{}
}

func NewSQLBuilder(table *Table, header *replication.EventHeader, ev *replication.RowsEvent) *SQLBuilder {
	return &SQLBuilder{
		Table:  table,
		Header: header,
		Event:  ev,
	}
}

func (sqlBld *SQLBuilder) SQLs() []*SQL {
	rowNum := len(sqlBld.Event.Rows)
	sqls := []*SQL{}

	switch sqlBld.Header.EventType {
	case replication.WRITE_ROWS_EVENTv2:
		for i := 0; i < rowNum; i++ {
			sql, params := sqlBld.buildInsertSQL(sqlBld.Event.Rows[i])
			sqls = append(sqls, &SQL{Statement: sql, Params: params})
		}
	case replication.UPDATE_ROWS_EVENTv2:
		for i := 0; i < rowNum; i += 2 {
			whereVals := sqlBld.Event.Rows[i]
			setVals := sqlBld.Event.Rows[i+1]
			sql, params := sqlBld.buildUpdateSQL(whereVals, setVals)
			sqls = append(sqls, &SQL{Statement: sql, Params: params})
		}
	case replication.DELETE_ROWS_EVENTv2:
		for i := 0; i < rowNum; i++ {
			sql, params := sqlBld.buildDeleteSQL(sqlBld.Event.Rows[i])
			sqls = append(sqls, &SQL{Statement: sql, Params: params})
		}
	}

	return sqls
}

func (sqlBld *SQLBuilder) buildInsertSQL(values []interface{}) (string, []interface{}) {
	var builder strings.Builder
	builder.WriteString("INSERT INTO `")
	builder.WriteString(sqlBld.Table.Schema)
	builder.WriteString("`.`")
	builder.WriteString(sqlBld.Table.Name)
	builder.WriteString("` (`")

	for i := 0; i < len(values); i++ {
		builder.WriteString(sqlBld.Table.ColumnNames[i])

		if i < len(values)-1 {
			builder.WriteString("`, `")
		}
	}

	builder.WriteString("`) VALUES (")

	for i := 0; i < len(values); i++ {
		builder.WriteString("?")

		if i < len(values)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString(")")

	return builder.String(), values
}

func (sqlBld *SQLBuilder) buildUpdateSQL(whereVals, setVals []interface{}) (string, []interface{}) {
	params := make([]interface{}, 0, len(setVals)+len(whereVals))
	params = append(params, setVals...)

	var builder strings.Builder
	builder.WriteString("UPDATE `")
	builder.WriteString(sqlBld.Table.Schema)
	builder.WriteString("`.`")
	builder.WriteString(sqlBld.Table.Name)
	builder.WriteString("` SET `")

	for i := 0; i < len(setVals); i++ {
		builder.WriteString(sqlBld.Table.ColumnNames[i])
		builder.WriteString("` = ?")

		if i < len(setVals)-1 {
			builder.WriteString(", `")
		}
	}

	builder.WriteString(" ")
	whereStmt, whereParams := sqlBld.buildWhereClause(whereVals)
	builder.WriteString(whereStmt)

	if len(whereParams) > 0 {
		params = append(params, whereParams...)
	}

	return builder.String(), params
}

func (sqlBld *SQLBuilder) buildDeleteSQL(values []interface{}) (string, []interface{}) {
	whereStmt, params := sqlBld.buildWhereClause(values)

	var builder strings.Builder
	builder.WriteString("DELETE FROM `")
	builder.WriteString(sqlBld.Table.Schema)
	builder.WriteString("`.`")
	builder.WriteString(sqlBld.Table.Name)
	builder.WriteString("` ")
	builder.WriteString(whereStmt)

	return builder.String(), params
}

func (sqlBld *SQLBuilder) buildWhereClause(values []interface{}) (string, []interface{}) {
	params := make([]interface{}, 0, len(values))

	var builder strings.Builder
	builder.WriteString("WHERE ")

	for i, v := range values {
		builder.WriteString(sqlBld.Table.ColumnNames[i])

		if v != nil {
			builder.WriteString(" = ?")
			params = append(params, v)
		} else {
			builder.WriteString(" IS NULL")
		}

		if i < len(values)-1 {
			builder.WriteString(" AND ")
		}
	}

	return builder.String(), params
}
