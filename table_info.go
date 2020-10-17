package binrpt

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/siddontang/go-mysql/client"
)

const (
	TableInfoDefaultExpiration = 5 * time.Minute
	TableInfoCleanupInterval   = 10 * time.Minute
)

type TableInfo struct {
	Cache *cache.Cache
	Conn  *client.Conn
}

type Table struct {
	Schema      string
	Name        string
	ColumnNames []string
	ColumnCount uint64
}

func NewTableInfo(conn *client.Conn) *TableInfo {
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
	r, err := tableInfo.Conn.Execute(`
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

	rowNum := r.RowNumber()

	if rowNum < 1 {
		return nil, nil
	}

	columnNames := make([]string, rowNum)

	for i := 0; i < rowNum; i++ {
		columnNames[i], err = r.GetString(i, 0)

		if err != nil {
			return nil, err
		}
	}

	return columnNames, nil
}
