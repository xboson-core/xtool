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
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	opt := &Options{}
	var pack, diff, sqld bool

	f.BoolVar(&pack, "pack", false, "build pack zip")
	f.BoolVar(&diff, "diff", false, "build diff zip")
	f.BoolVar(&sqld, "sqld", false, "build sql diff")
	f.StringVar(&opt.Input, "i", "", "input directory / input sql file")
	f.StringVar(&opt.Output, "o", "", "output zip or sql file")
	f.StringVar(&opt.Base, "b", "", "base zip / base sql file")
	f.StringVar(&opt.Meta, "meta", "", "meta csv")
	f.StringVar(&opt.ExcludeFrom, "exclude-from", "", "exclude file")
	f.StringVar(&opt.DeleteFile, "del", "", "save removed file")
	f.BoolVar(&opt.Verbose, "v", false, "show progress")
	f.BoolVar(&opt.Line, "l", false, "show progress in one line")
	f.Parse(os.Args[1:])
	
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
		f.Usage()
		os.Exit(2)
	}

	if pack && opt.Input == "" {
		fmt.Fprintln(os.Stderr, "-i is required in pack mode")
		os.Exit(2)
	}

	if diff {
		if opt.Input == "" {
			fmt.Fprintln(os.Stderr, "-i is required in diff mode")
			os.Exit(2)
		}
		if opt.Base == "" {
			fmt.Fprintln(os.Stderr, "-b is required in diff mode")
			os.Exit(2)
		}
	}

	if sqld {
		if opt.Input == "" {
			fmt.Fprintln(os.Stderr, "-i is required in sqld mode")
			os.Exit(2)
		}
		if opt.Base == "" {
			fmt.Fprintln(os.Stderr, "-b is required in sqld mode")
			os.Exit(2)
		}
	}

	if opt.Output == "" {
		fmt.Fprintln(os.Stderr, "-o is required")
		os.Exit(2)
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