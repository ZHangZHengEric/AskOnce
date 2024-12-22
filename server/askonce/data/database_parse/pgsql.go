// Package database_parse -----------------------------
// @file      : pgsql.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:57
// -------------------------------------------
package database_parse

import "fmt"
import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgreSQLHandler 实现 DatabaseHandler 接口
type PostgreSQLHandler struct {
	DB *sqlx.DB
}

func (h *PostgreSQLHandler) Ping() error {
	return h.DB.Ping()
}

func (h *PostgreSQLHandler) GetTables() ([]TableInfo, error) {
	query := " SELECT t.table_name, COALESCE(d.description, '') AS table_comment  FROM information_schema.tables t LEFT JOIN pg_class c ON c.relname = t.table_name LEFT JOIN pg_namespace n ON n.oid = c.relnamespace LEFT JOIN pg_description d ON c.oid = d.objoid  WHERE t.table_schema = 'public' AND n.nspname = 'public';`"
	var tables []TableInfo
	err := h.DB.Select(&tables, query)
	return tables, err
}

func (h *PostgreSQLHandler) GetColumns(table string) ([]ColumnInfo, error) {
	query := fmt.Sprintf(`
		SELECT column_name, data_type, col_description((SELECT oid FROM pg_class WHERE relname='%s'), ordinal_position::int) AS column_comment 
		FROM information_schema.columns 
		WHERE table_name='%s';`, table, table)

	rows, err := h.DB.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		column := make(map[string]interface{})
		if err := rows.MapScan(column); err != nil {
			return nil, err
		}
		// 转换为 map[string]string
		columnData := ColumnInfo{
			ColumnName:    fmt.Sprintf("%v", column["column_name"]),
			ColumnType:    fmt.Sprintf("%v", column["data_type"]),
			ColumnComment: fmt.Sprintf("%v", column["column_comment"]),
		}
		columns = append(columns, columnData)
	}
	return columns, nil
}

func (h *PostgreSQLHandler) GetSampleData(table, column string) ([]string, error) {
	query := fmt.Sprintf("SELECT %s FROM %s LIMIT 1000;", column, table)
	var samples []string
	err := h.DB.Select(&samples, query)
	return samples, err
}

func (h *PostgreSQLHandler) Close() error {
	h.DB.Close()
	return nil
}
