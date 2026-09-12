package xtool

import "fmt"


type Show struct {
	Sub  FileProcess
	Line bool
	i    int64
}

func (s *Show) OnFile(m *Meta) error {
	s.i += 1
	if s.Line {
		fmt.Printf("\r%5d %s", s.i, m.FullPath)
	} else {
		fmt.Printf("%5d %s\n", s.i, m.FullPath)
	}
	if s.Sub != nil {
		return s.Sub.OnFile(m)
	}
	return nil
}

func (s *Show) OnFinish() {
	if s.Line {
		fmt.Printf("\r-------------------------------------------------\n")
		fmt.Println("Total", s.i, "files")
	}
	if s.Sub != nil {
		s.Sub.OnFinish()
	}
}