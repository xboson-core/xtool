package xtool

import (
	"fmt"
	"io"
	"os"

	"vitess.io/vitess/go/vt/sqlparser"
)


type Render func(io.Writer)(int, error)


type Schema struct {
	s string
}


func (s *Schema) Name(t string) string {
	if s.s == "" {
		return t
	}
	return s.s +"."+ t
}


type Table struct {
	Name 				string
	Columns 		[]string
	PrimaryKey 	[]string
	Rows 				[][]string
	index       map[string]int
	ColDef      map[string]*sqlparser.ColumnDefinition
}


type TableDataBuilder struct {
	tables 			map[string]*Table
	renders 	 	[]Render
	Render 		 	bool
	nextInsert 	func(*Table, string)error
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
		s.putrs(sql)

	case *sqlparser.Update:
			fmt.Println(">>> UPDATE")

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
	if !s.Render {
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
	// 这里应该检测表结构不同做什么
	// 如果baset有的列t没有,则输出 ALTER TABLE DROP COLUMN
	// 如果t有的列baset没有则输出 ALTER TABLE ADD COLUMN
	// 用 Table {ColDef []*sqlparser.ColumnDefinition}
	baset, _ := d.base.GetTable(t.Name)
	// 如果是新建表 则全部输出
	if baset == nil {
		d.putrs(ins)
		t.Rows = nil
		return nil
	}
	//如果baset有行的t没有则输出 delete 
	//如果t有的baset没有则输出 insert
	//如果都有则按列比较不同的数据进行update,如果表结构不同以t为基础迭代
	//比较完成把 t.rows清空
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
        "ALTER TABLE `%s` DROP COLUMN `%s`;",
        cur.Name, name,
      ))
    }
  }

  for name, col := range cur.ColDef {
    if _, ok := base.ColDef[name]; !ok {
      d.putrs(fmt.Sprintf(
        "ALTER TABLE `%s` ADD COLUMN %s;",
        cur.Name,
        sqlparser.String(col),
      ))
    }
  }
  return nil
}


func (d *DiffDataBuilder) diffRows(base, cur *Table) {
  for _, row := range base.Rows {
    key := base.rowKey(row)
    if _, ok := cur.index[key]; !ok {
      d.putrs(d.makeDelete(cur, row))
    }
  }

  for _, row := range cur.Rows {
    key := cur.rowKey(row)
    i, ok := base.index[key]
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


func (t *Table) rowKey(row []string) string {
  var key string
  for _, name := range t.PrimaryKey {
    for i, col := range t.Columns {
      if col == name {
        key += row[i]
        key += "\x00"
        break
      }
    }
  }
  return key
}


func sameRow(base *Table, old []string, cur *Table, row []string) bool {
  for name, _ := range cur.ColDef {
    if _, ok := base.ColDef[name]; !ok {
      continue
    }

    oi := columnIndex(base.Columns, name)
    ni := columnIndex(cur.Columns, name)
    if oi < 0 || ni < 0 {
      continue
    }

    if old[oi] != row[ni] {
      return false
    }
  }
  return true
}


func columnIndex(cols []string, name string) int {
  for i, col := range cols {
    if col == name {
      return i
    }
  }
  return -1
}


func (d *DiffDataBuilder) makeUpdate(t *Table, old, new []string) string {
	//TODO
	return fmt.Sprintf("Update %s Set %s Where %s", t.Name, new, old)
}


func (d *DiffDataBuilder) makeDelete(t *Table, row []string) string {
	//TODO
	return fmt.Sprintf("Delete %s where %s", t.Name, row)
}


func (d *DiffDataBuilder) makeInsert(t *Table, row []string) string {
	//TODO
	return fmt.Sprintf("Insert into %s (%s) Values (%s)", t.Name, t.Columns, row)
}


func (d *DiffDataBuilder) OnFinish() {
	// 迭代 base中的表, 在d中找不到就输出 
	for name := range d.base.tables {
    if _, ok := d.tables[name]; !ok {
      d.putrs(fmt.Sprintf("DROP TABLE IF EXISTS %s;", name))
    }
  }
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