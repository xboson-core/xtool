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
		Name: sch.Name(node.Table.Name.String()),
		ColDef: make(map[string]*sqlparser.ColumnDefinition),
		index: make(map[string]int),
	}
	t.Columns = make([]string, 0, len(node.TableSpec.Columns))
	for _, col := range node.TableSpec.Columns {
		t.Columns = append(t.Columns, col.Name.String())
		t.ColDef[col.Name.String()] = col
	}
	for _, idx := range node.TableSpec.Indexes {
		if idx.Info != nil &&
			idx.Info.Type == sqlparser.IndexTypePrimary {
			for _, col := range idx.Columns {
				t.PrimaryKey = append(t.PrimaryKey, col.Column.String())
			}
		}
	}
  return t
}


func parseInsert(node *sqlparser.Insert, t *Table) {
	var cols []string
	if node.Columns == nil {
		cols = t.Columns
	} else {
		for _, col := range node.Columns {
			cols = append(cols, col.String())
		}
	}

	switch rows := node.Rows.(type) {
	case sqlparser.Values:
		for _, row := range rows {
			_row := make([]string, 0, len(rows))
			for _, val := range row {
				_row = append(_row, sqlparser.String(val))
			}
			rowNumber := len(t.Rows)
			t.Rows = append(t.Rows, _row)

			key := t.rowKey(_row)
			if _, has := t.index[key]; has {
				panic(fmt.Errorf("Primary key conflict %s", key))
			}
      t.index[key] = rowNumber
		}

	default:
		panic(fmt.Errorf("unsupported INSERT rows type: %T\n", node.Rows))
	}
}
