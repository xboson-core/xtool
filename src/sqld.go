package xtool

import (
	"fmt"
	"io"
	"os"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)


type Render func(io.Writer)(int, error)


type Schema struct {
	schema string
}


func (s *Schema) Name(t string) string {
	if s.schema == "" {
		return t
	}
	return s.schema +"."+ t
}


func (s *Schema) SafeName(t string) string {
	if s.schema == "" {
		return "`"+ t +"`"
	}
	return "`"+ s.schema +"`.`"+ t +"`"
}


type Table struct {
	Name 				string
	schema      Schema
	Columns 		[]string
	// [列索引]`列名`
	SafeCol     []string
	// [主键名]列索引
	PrimaryKey 	map[string]int
	// [主键列索引]
	PrimaryIndex []int
	// 任何字面值都不可能是空字符串
	Rows 				[][]string
	// [主键值]Rows索引
	pkrow_index map[string]int
	// [列名]列索引
	col_index   map[string]int
	// [列名]列定义
	ColDef      map[string]*sqlparser.ColumnDefinition
}


func (t *Table) SafeName() string {
	return t.schema.SafeName(t.Name)
}


func (t *Table) rowPKeyValue(row []string) string {
  var key string
	for _, pi := range t.PrimaryIndex {
		key += row[pi]
		key += "\x00"
	}
  return key
}


type TableDataBuilder struct {
	tables 			map[string]*Table
	renders 	 	[]Render
	Render 		 	bool
	nextInsert 	func(*Table, string)error
	nextCreate  func(*Table, string)error
	Schema      Schema
}


func (s *TableDataBuilder) AddTable(t *Table) error {
	if _, has := s.tables[t.Name]; has {
		return fmt.Errorf("fail: table conflict '%s'", t.Name)
	}
	s.tables[t.Name] = t
	return nil
}


func (s *TableDataBuilder) GetTable(fullname string) (*Table, error) {
	t, has := s.tables[fullname]
	if !has {
		return nil, fmt.Errorf("Table not defined: %s", fullname)
	}
	return t, nil
}


func (s *TableDataBuilder) OnSql(st sqlparser.Statement, sql string) error {
	switch node := st.(type) {

	case *sqlparser.Insert:
		table, err := s.GetTable(s.Schema.Name(node.Table.TableNameString()))
		if err != nil {
			return err
		}
		parseInsert(node, table)
		if s.nextInsert != nil {
			if err := s.nextInsert(table, sql); err != nil {
				return err
			}
		}

	case *sqlparser.CreateTable:
		table := parseCreateTable(node, s.Schema)
		if err := s.AddTable(table); err != nil {
			return err
		}
		if s.nextCreate != nil {
			if err := s.nextCreate(table, sql); err != nil {
				return err
			}
		}


	case *sqlparser.Update:
		panic(fmt.Errorf("Not support: %s", sql))

	case *sqlparser.Use:
		s.Schema = Schema{ node.DBName.String() }
		s.putrs(sql)

	case *sqlparser.DropTable:
	case *sqlparser.DropDatabase:
	case *sqlparser.DropView:
		// donothing

	default:
		s.putrs(sql)
	}

	return nil
}


func (s *TableDataBuilder) put(r Render) {
	if s.Render {
		s.renders = append(s.renders, r)
	}
}


func (s *TableDataBuilder) putrs(str string) {
	if str=="" || (!s.Render) {
		return
	}
	s.put(func(w io.Writer) (int, error) {
		l := len(str)
		if l>0 && (str[l-1] != '\n') {
			defer w.Write([]byte{'\n'})
		}
		return w.Write([]byte(str))
	})
}


func (s *TableDataBuilder) RenderFile(file string) error {
	if !s.Render {
		panic("cannot render file")
	}
	fout, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fout.Close()
	total := 0
	for _, render := range s.renders {
		if len, err := render(fout); err != nil {
			return err
		} else {
			total += len
		}
	}
	fmt.Println("Write", file, total, "bytes,", len(s.renders ), "renders")
	return nil
}


func (s *TableDataBuilder) OnFinish() {
	fmt.Println("sql load finished")
}


func NewTableDataBuilder() *TableDataBuilder {
	return &TableDataBuilder{
		Render: false,
		tables: make(map[string]*Table),
	}
}


type DiffDataBuilder struct {
	TableDataBuilder
	base *TableDataBuilder
}


func (d *DiffDataBuilder) writeInsert(t *Table, ins string) error {
	baset, _ := d.base.GetTable(t.Name)
	// 如果是新建表 则全部输出
	if baset == nil {
		d.putrs(ins)
		t.Rows = nil
		return nil
	}
	if err := d.diffColumns(baset, t); err != nil {
    return err
  }
  d.diffRows(baset, t)
  t.Rows = nil
  return nil
}


func (d *DiffDataBuilder) diffColumns(base, cur *Table) error {
  for name := range base.ColDef {
    if _, ok := cur.ColDef[name]; !ok {
      d.putrs(fmt.Sprintf(
        "ALTER TABLE %s DROP COLUMN %s;",
        cur.SafeName(), name,
      ))
    }
  }

  for name, col := range cur.ColDef {
    if _, ok := base.ColDef[name]; !ok {
      d.putrs(fmt.Sprintf(
        "ALTER TABLE %s ADD COLUMN %s;",
        cur.SafeName(), sqlparser.String(col),
      ))
    }
  }
  return nil
}


func (d *DiffDataBuilder) diffRows(base, cur *Table) {
  for _, row := range base.Rows {
    key := base.rowPKeyValue(row)
    if _, ok := cur.pkrow_index[key]; !ok {
      d.putrs(d.makeDelete(cur, row))
    }
  }

  for _, row := range cur.Rows {
    key := cur.rowPKeyValue(row)
    i, ok := base.pkrow_index[key]
    if !ok {
      d.putrs(d.makeInsert(cur, row))
      continue
    }

    old := base.Rows[i]
    if !sameRow(base, old, cur, row) {
      d.putrs(d.makeUpdate(cur, old, row))
    }
  }
}


func sameRow(base *Table, old []string, cur *Table, row []string) bool {
  for name, _ := range cur.ColDef {
    if _, ok := base.ColDef[name]; !ok {
      return false
    }

    oi, ohas := base.col_index[name]
    ni, chas := cur.col_index[name]
    if (!ohas) || (!chas) {
      return false
    }

    if old[oi] != row[ni] {
      return false
    }
  }
  return true
}


func makeWhereWithPK(t *Table, row []string) string {
	where := strings.Builder{}
	wi := 0
	for _, i := range t.PrimaryIndex {
		if wi > 0 {
			where.WriteString(" AND ")
		}
		where.WriteString(t.SafeCol[i])
		where.WriteString(" = ")
		where.WriteString(row[i])
		wi += 1
	}
	return where.String()
}


func (d *DiffDataBuilder) makeUpdate(t *Table, old, new []string) string {
	where := makeWhereWithPK(t, new)
	set := strings.Builder{}
	si := 0
	for i, col := range t.Columns {
		if new[i] == "" {
			continue
		}
		if _, has := t.PrimaryKey[col]; has {
			continue
		}
		ot, err := d.base.GetTable(t.Name)
		if err != nil {
			continue
		}
		oi := ot.col_index[col]
		if oi>=len(old) || i>=len(new) || old[oi]==new[i] {
			continue
		}

		if si > 0 {
			set.WriteString(", ")
		}
		set.WriteString(t.SafeCol[i])
		set.WriteString(" = ")
		set.WriteString(new[i])
		si += 1
	}
	if si == 0 {
		return ""
	}
	return fmt.Sprintf(
		"UPDATE %s \n\tSET %s \n\tWHERE %s", 
		t.SafeName(), set.String(), where)
}


func (d *DiffDataBuilder) makeDelete(t *Table, row []string) string {
	where := makeWhereWithPK(t, row)
	return fmt.Sprintf("DELETE FROM %s where %s", t.SafeName(), where)
}


func (d *DiffDataBuilder) makeInsert(t *Table, row []string) string {
	_cols := strings.Builder{}
	_rows := strings.Builder{}
	c := 0
	for i, v := range row {
		if v == "" {
			continue
		}
		if c > 0 {
			_rows.WriteString(", ")
			_cols.WriteString(", ")
		}
		_rows.WriteString(v)
		_cols.WriteString(t.SafeCol[i])
		c += 1
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES \n\t(%s)", 
		t.Name, _cols.String(), _rows.String())
}


func (d *DiffDataBuilder) OnFinish() {
	// 迭代 base中的表, 在d中找不到就输出 
	for name := range d.base.tables {
    if _, ok := d.tables[name]; !ok {
			safename := d.base.tables[name].SafeName()
      d.putrs(fmt.Sprintf("DROP TABLE IF EXISTS %s;", safename))
    }
  }
}


func (d *DiffDataBuilder) writeCreate(t *Table, sql string) error {
	if _, has := d.base.tables[t.Name]; !has {
		d.putrs(sql)
	}
	return nil
}


func NewDiffDataBuilder(b *TableDataBuilder) *DiffDataBuilder {
	ret := &DiffDataBuilder{
		TableDataBuilder: TableDataBuilder{
			Render: true,
			tables: make(map[string]*Table),
		},
		base: b,
	}
	ret.nextInsert = ret.writeInsert
	ret.nextCreate = ret.writeCreate
	return ret
}


func SqlDiff(o *Options) error {
	base := NewTableDataBuilder()
	if err := EachSqlFrom(o.Base, base); err != nil {
		return err
	}
	update := NewDiffDataBuilder(base)
	if err := EachSqlFrom(o.Input, update); err != nil {
		return err
	}
	return update.RenderFile(o.Output)
}