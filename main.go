package main

import (
	"flag"
	"fmt"
)

var (
	profileFlag string
	setenvFlag  bool
	rotateFlag  bool
	versionFlag bool
	version     string
)

func init() {
	flag.StringVar(&profileFlag, "profile", "", "Profile to use")
	flag.BoolVar(&rotateFlag, "rotate", false, "Show version")
	flag.BoolVar(&setenvFlag, "setenv", false, "Show version")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()
}

func main() {
	if versionFlag {
		fmt.Printf("\tAKT %v\n\n", version)
		return
	}
	if profileFlag == "" {
		flag.PrintDefaults()
		return
	}
}
