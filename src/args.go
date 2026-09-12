package xtool

import (
	"flag"
	"fmt"
	"os"
)


type Options struct {
	Mode        string
	Input       string
	Output      string
	Base        string
	Meta        string
	ExcludeFrom string
	DeleteFile  string
	Verbose     bool
	Line        bool
}


func (o *Options) MakeWrap(fp FileProcess) FileProcess {
	if o.Verbose {
		return &Show{ Sub: fp, Line: o.Line }
	} 
	return fp
}


func ParseArgs() *Options {
	opt := &Options{}
	var pack, diff, sqld bool

	flag.BoolVar(&pack, "pack", false, "build pack zip")
	flag.BoolVar(&diff, "diff", false, "build diff zip")
	flag.BoolVar(&sqld, "sqld", false, "build sql diff")
	flag.StringVar(&opt.Input, "i", "", "input directory")
	flag.StringVar(&opt.Output, "o", "", "output zip")
	flag.StringVar(&opt.Base, "b", "", "base zip")
	flag.StringVar(&opt.Meta, "meta", "", "meta csv")
	flag.StringVar(&opt.ExcludeFrom, "exclude-from", "", "exclude file")
	flag.StringVar(&opt.DeleteFile, "del", "", "save removed file")
	flag.BoolVar(&opt.Verbose, "v", false, "show progress")
	flag.BoolVar(&opt.Line, "l", false, "show progress in one line")
	flag.Parse()

	n := 0
	if pack {
		opt.Mode = "pack"
		n++
	}
	if diff {
		opt.Mode = "diff"
		n++
	}
	if sqld {
		opt.Mode = "sqld"
		n++
	}
	if n != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if opt.Mode == "pack" && opt.Input == "" {
		fmt.Fprintln(os.Stderr, "-i is required in pack mode")
		os.Exit(2)
	}
	if opt.Mode == "diff" {
		if opt.Input == "" {
			fmt.Fprintln(os.Stderr, "-i is required in diff mode")
			os.Exit(2)
		}
		if opt.Base == "" {
			fmt.Fprintln(os.Stderr, "-b is required in diff mode")
			os.Exit(2)
		}
	}
	
	if opt.Meta == "" {
		if pack {
			opt.Meta = opt.Output + ".csv"
		} else if diff {
			opt.Meta = opt.Base + ".csv"
		}
	}
	return opt
}