package xtool

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)


type SaveZipCsv struct {
	Csv string
	Zip string
	wcsv *csv.Writer
	wzip *zip.Writer
}


func NewSaveZipCsv(csvpath, zippath string) (*SaveZipCsv, error) {
	csvfile, err := os.OpenFile(csvpath, os.O_CREATE, 0700)
	if err != nil {
		return nil, err
	}
	zipfile, err := os.OpenFile(zippath, os.O_CREATE, 0700)
	if err != nil {
		return nil, err
	}

	return &SaveZipCsv{
		Csv : csvpath,
		Zip : zippath,
		wcsv : csv.NewWriter(csvfile),
		wzip : zip.NewWriter(zipfile),
	}, nil
}


func (s *SaveZipCsv) OnFile(m *Meta) error {
	s.wcsv.Write(m.ToArr())

	f, err := s.wzip.Create(m.Path)
	if err != nil {
		return err
	}
	in, err := os.Open(m.FullPath)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(f, in)
	return err
}


func (s *SaveZipCsv) OnFinish() {
	if err := s.wzip.Close(); err != nil {
		fmt.Println(err)
	}
	s.wcsv.Flush()
	if err := s.wcsv.Error(); err != nil {
		fmt.Println(err)
	}
}


func PackMode(o *Options) error {
	e := EachOption{
		Base: o.Input,
	}
	if o.ExcludeFrom != "" {
		if err := e.ReadExclude(o.ExcludeFrom); err != nil {
			return err
		}
	}
	zc, err := NewSaveZipCsv(o.Meta, o.Output)
	if err != nil {
		return err
	}
	return EachFile(&e, o.MakeWrap(zc))
}