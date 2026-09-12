package xtool

import (
	"os"
	"strings"
)

// 如果 put 'a/b/c' 'a/b'
// a/b 节点的end=true, Find(a/b/xxx) 返回true
type Leaf struct {
	Name string
	Subs map[string]*Leaf
	// 这可以是一个终结节点
	end bool
}


func (l *Leaf) Put(name string) {
  parts := strings.FieldsFunc(name, func(r rune) bool {
    return r == '/' || r == '\\'
  })
  cur := l
  for _, part := range parts {
    if cur.Subs == nil {
      cur.Subs = make(map[string]*Leaf)
    }
    next := cur.Subs[part]
    if next == nil {
      next = &Leaf{Name: part}
      cur.Subs[part] = next
    }
    cur = next
  }
  cur.end = true
}


func (l *Leaf) Find(name string) bool {
  parts := strings.FieldsFunc(name, func(r rune) bool {
    return r == os.PathSeparator
  })
  cur := l
  for _, part := range parts {
    if cur.end {
      return true
    }
    cur = cur.Subs[part]
    if cur == nil {
      return false
    }
  }
  return cur.end
}