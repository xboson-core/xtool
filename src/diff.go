package xtool

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)


type LoadMetaDiff struct {
	zip *zip.Writer
	files map[string]Meta
}


func (l *LoadMetaDiff) OnFile(new *Meta) error {
	old, has := l.files[new.FullPath]
	if has {
		if new.Size == old.Size && new.Unix <= old.Unix {
			return nil
		}
	}

	f, err := l.zip.Create(new.Path)
	if err != nil {
		return err
	}
	in, err := os.Open(new.FullPath)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(f, in)
	return err
}


func (l *LoadMetaDiff) OnFinish() {
	if err := l.zip.Close(); err != nil {
		fmt.Println(err)
	}
}


func read_csv(csvpath, delfile string) (map[string]Meta, error) {
	csvfile, err := os.Open(csvpath)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	csv := csv.NewReader(csvfile)
	files := make(map[string]Meta)

	var del *os.File
	if delfile != "" {
		del, err = os.OpenFile(delfile, os.O_CREATE|os.O_TRUNC, 0700)
		if err != nil {
			return nil, err
		}
		defer del.Close()
	}

	for {
		row, err := csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		m := Meta{}
		m.FromArr(row)
		files[m.FullPath] = m

		if del != nil {
			if _, err := os.Stat(m.FullPath); os.IsNotExist(err) {
				del.WriteString(m.FullPath)
				del.Write([]byte{ '\n' })
			}
		}
	}
	return files, nil
}


func NewLoadMetaDiff(csvpath, zippath, del string) (*LoadMetaDiff, error) {
	zipfile, err := os.OpenFile(zippath, os.O_CREATE, 0700)
	if err != nil {
		return nil, err
	}

	files, err := read_csv(csvpath, del)
	if err != nil {
		return nil, err
	}

	return &LoadMetaDiff{
		files 	: files,
		zip 		: zip.NewWriter(zipfile),
	}, nil
}


func DiffMode(o *Options) error {
	e := EachOption{
		Base: o.Input,
	}
	if o.ExcludeFrom != "" {
		if err := e.ReadExclude(o.ExcludeFrom); err != nil {
			return err
		}
	}
	diff, err := NewLoadMetaDiff(o.Meta, o.Output, o.DeleteFile)
	if err != nil {
		return err
	}
	return EachFile(&e, o.MakeWrap(diff))
}