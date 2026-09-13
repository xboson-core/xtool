package xtool

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)


type FileProcess interface {
	OnFile(*Meta) error
	OnFinish()
}


type Meta struct {
	// 本地完整路径
	FullPath string
	// 剪切后的路径
	Path string
	// 文件修改时间
	Unix int64
	// 文件大小
	Size int64
}


func (m *Meta) ToArr() []string {
	return []string {
		m.FullPath, 
		strconv.FormatInt(m.Unix, 10), 
		strconv.FormatInt(m.Size, 10),
	}
}


func (m *Meta) FromArr(v []string) error {
  if len(v) < 3 {
    return fmt.Errorf("invalid meta record")
  }
  unix, err := strconv.ParseInt(v[1], 10, 64)
  if err != nil {
    return err
  }
  size, err := strconv.ParseInt(v[2], 10, 64)
  if err != nil {
    return err
  }
  m.FullPath = v[0]
  m.Unix = unix
  m.Size = size
  return nil
}


type EachOption struct {
	Base string
	Root *Leaf
}


func (e *EachOption) Ignore(pathName string) bool {
	if e.Root == nil {
		return false
	}
	return e.Root.Find(pathName)
}


func (e *EachOption) ReadExclude(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	e.Root = &Leaf{
		Name: "",
		Subs: make(map[string]*Leaf),
	}
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line=="" || strings.HasPrefix(line, "#") {
			continue
		}
		e.Root.Put(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}


// 迭代目录, 用 cut 剪掉前缀
func EachFile(opt *EachOption, fp FileProcess) error {
	err := filepath.WalkDir(opt.Base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		absPath := path[len(opt.Base)+1:]
		if opt.Ignore(absPath) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return fp.OnFile(&Meta{
			FullPath: path,
			Path:     absPath,
			Unix:     info.ModTime().Unix(),
			Size:     info.Size(),
		})
	})
	if err == nil {
		fp.OnFinish()
	}
	return err
}