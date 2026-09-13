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
	OnSql(st sqlparser.Statement, sql string, lineNum int) error
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

	fmt.Println("Loading file:", file)
	r := bufio.NewReader(f)
	var buf strings.Builder
	linenum := 1

	var parse = func(sql string) error {
		if sql != "" {
			stmt, err := parser.Parse(sql)
			if err != nil {
				fmt.Println("[WARN]", err)
				// fmt.Println(sql)
				return nil
			}
			if err := sp.OnSql(stmt, sql, linenum); err != nil {
				return err
			}
		}
		return nil
	}

	var quote byte = 0
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		if ch == '\n' {
			linenum += 1
			continue
		}

		buf.WriteByte(ch)
		if quote != 0 {
			// SQL 字符串中的反斜杠转义
			if ch == '\\' {
				if next, err := r.ReadByte(); err == nil {
					buf.WriteByte(next)
					if next == '\n' {
						linenum += 1
					}
				}
				continue
			}
			// ''、""、`` 表示转义/重复的 quote
			if ch == quote {
				next, err := r.ReadByte()
				if err == nil {
					if next != quote {
						// 不是连续 quote，把它放回去
						_ = r.UnreadByte()
						quote = 0
					}
					continue
				}
				quote = 0
			}
			continue
		}
		// 进入字符串
		if ch == '\'' || ch == '"' || ch == '`' {
			quote = ch
			continue
		}
		// 只有不在字符串中的 ; 才进行切割
		if ch == ';' {
			sql := buf.String()
			buf.Reset()

			if err := parse(sql); err != nil {
				return err
			}
		}
		quote = 0
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


func parseInsert(node *sqlparser.Insert, t *Table, lm int) {
	if len(t.PrimaryKey) < 1 {
		fmt.Printf("[WARN] Table %s has no primary key at:%d,\n", t.SafeName(), lm)
		fmt.Println("\tINSERT operations cannot be processed, and all data is ignored.")
		return
	}

	colsIndexIndex := make(map[int]int)
	tableName := node.Table.TableNameString()
	if t.Columns == nil {
		panic(fmt.Errorf("Insert Table %s but not has DDL at:%d", tableName, lm))
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
		panic(fmt.Errorf("unsupported INSERT rows type: %T, at:%d", node.Rows, lm))
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
				fmt.Errorf("Primary key conflict { %s=%s } (%s), Table: %s, at:%d", 
					t.PrimaryKey, pkValue, cf, t.Name, lm))
		}
		t.pkrow_index[pkValue] = rowNumber
	}
}
