package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	yellowFmt = color.New(color.FgYellow).SprintFunc()
	greenFmt  = color.New(color.FgGreen).SprintFunc()
	redFmt    = color.New(color.FgRed).SprintFunc()
	grayFmt   = color.New(color.Faint).SprintFunc()
)

func Errorln(a ...interface{}) (int, error) {
	return Errorf("%s\n", fmt.Sprint(a...))
}

func Errorf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, redFmt("error: ")+format, a...)
}
