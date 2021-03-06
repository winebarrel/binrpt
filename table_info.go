package binrpt

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
)

const (
	TableInfoDefaultExpiration = 5 * time.Minute
	TableInfoCleanupInterval   = 10 * time.Minute
)

type TableInfo struct {
	Cache *cache.Cache
	Conn  *sql.DB
}

type Table struct {
	Schema      string
	Name        string
	ColumnNames []string
	ColumnCount uint64
}

func NewTableInfo(conn *sql.DB) *TableInfo {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &TableInfo{Cache: c, Conn: conn}
}

func (tableInfo *TableInfo) Get(schema string, name string) (*Table, error) {
	key := schema + "." + name
	t, found := tableInfo.Cache.Get(key)

	if found {
		return t.(*Table), nil
	}

	colmunNames, err := tableInfo.getColumnNames(schema, name)

	if err != nil {
		return nil, err
	}

	if colmunNames == nil {
		return nil, nil
	}

	newTbl := &Table{
		Schema:      schema,
		Name:        name,
		ColumnNames: colmunNames,
		ColumnCount: uint64(len(colmunNames)),
	}

	tableInfo.Cache.Set(key, newTbl, cache.DefaultExpiration)

	return newTbl, nil
}

func (tableInfo *TableInfo) getColumnNames(schema string, name string) ([]string, error) {
	rows, err := tableInfo.Conn.Query(`
		SELECT
			COLUMN_NAME
		FROM
			information_schema.COLUMNS
		WHERE
			TABLE_SCHEMA = ?
			AND TABLE_NAME = ?
		ORDER BY
			ORDINAL_POSITION
`, schema, name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	columnNames := []string{}

	for rows.Next() {
		var colName string
		err = rows.Scan(&colName)

		if err != nil {
			return nil, err
		}

		columnNames = append(columnNames, colName)
	}

	if len(columnNames) == 0 {
		return nil, nil
	}

	return columnNames, nil
}
