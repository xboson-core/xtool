package main

import (
	"fmt"

	xtool "com.xboson.com/xtool/v2/src"
)


func diff(opt *xtool.Options) error {
	return nil
}


func sql_diff(opt *xtool.Options) error {
	return nil
}


func main() {
  opt := xtool.ParseArgs()
	var err error = nil

  switch opt.Mode {
  case "pack":
    err = xtool.PackMode(opt)
  case "diff":
    err = xtool.DiffMode(opt)
  case "sqld":
    err = xtool.SqlDiff(opt)
  }

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Done")
	}
}