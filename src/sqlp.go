package xtool

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)


type SqlProcess interface {
	OnSql(st sqlparser.Statement, sql string) error
	OnFinish()
}


func EachSqlFrom(file string, sp SqlProcess) error {
	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var buf strings.Builder
	fmt.Println("Loading file:", file)

	var parse = func(sql string) error {
		if sql != "" {
			stmt, err := parser.Parse(sql)
			if err != nil {
				// fmt.Println(" --", err)
				return nil
			}

			if err := sp.OnSql(stmt, sql); err != nil {
				return err
			}
		}
		return nil
	}

	for {
		chunk, err := r.ReadString(';')
		if len(chunk) > 0 {
			buf.WriteString(chunk)
			if chunk[len(chunk)-1] == ';' {
				sql := buf.String()
				buf.Reset()
				if err := parse(sql); err != nil {
					return err
				}
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}

	if err := parse(buf.String()); err != nil {
		return err
	}
	sp.OnFinish()
	return nil
}


func parseCreateTable(node *sqlparser.CreateTable, sch Schema) *Table {
	t := &Table{
		Name: 				sch.Name(node.Table.Name.String()),
		PrimaryKey: 	make(map[string]int),
		pkrow_index: 	make(map[string]int),
		col_index: 		make(map[string]int),
		ColDef: 			make(map[string]*sqlparser.ColumnDefinition),
		schema:       sch,
	}
	t.Columns = make([]string, len(node.TableSpec.Columns))
	t.SafeCol = make([]string, len(node.TableSpec.Columns))
	for i, col := range node.TableSpec.Columns {
		colName := col.Name.String()
		t.Columns[i] = colName
		t.SafeCol[i] = "`"+ colName +"`"
		t.ColDef[colName] = col
		t.col_index[colName] = i
	}
	for _, idx := range node.TableSpec.Indexes {
		if idx.Info != nil &&
			idx.Info.Type == sqlparser.IndexTypePrimary {
			for _, col := range idx.Columns {
				colName := col.Column.String()
				colIndex := t.col_index[colName]
				t.PrimaryKey[colName] = colIndex
				t.PrimaryIndex = append(t.PrimaryIndex, colIndex)
			}
		}
	}
  return t
}


func parseInsert(node *sqlparser.Insert, t *Table) {
	colsIndexIndex := make(map[int]int)
	tableName := node.Table.TableNameString()
	if t.Columns == nil {
		panic(fmt.Errorf("Insert Table %s but not has DDL", tableName))
	}
	if len(node.Columns) < 1 {
		for i, _ := range t.Columns {
			colsIndexIndex[i] = i
		}
	} else {
		colsIndex := make(map[string]int)
		for i, col := range node.Columns {
			colsIndex[col.String()] = i
		}
		for i, n := range t.Columns {
			if c, has := colsIndex[n]; has {
				colsIndexIndex[i] = c
			} else {
				colsIndexIndex[i] = -1
			}
		}
	}

	rows, ok := node.Rows.(sqlparser.Values)
	if !ok {
		panic(fmt.Errorf("unsupported INSERT rows type: %T\n", node.Rows))
	}

	for _, row := range rows {
		_row := make([]string, len(t.Columns))
		
		for i, _ := range t.Columns {
			ci := colsIndexIndex[i]
			if ci < 0 || ci >= len(row) {
				continue
			}
			val := row[ci]
			if val == nil {
				continue
			} 
			_row[i] = sqlparser.String(val)
		}
		rowNumber := len(t.Rows)
		t.Rows = append(t.Rows, _row)

		pkValue := t.rowPKeyValue(_row)
		if cf, has := t.pkrow_index[pkValue]; has {
			panic(
				fmt.Errorf("Primary key conflict { %s=%s } (%s), Table: %s", 
					t.PrimaryKey, pkValue, cf, t.Name))
		}
		t.pkrow_index[pkValue] = rowNumber
	}
}
