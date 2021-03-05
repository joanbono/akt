package main

import (
	"flag"
	"fmt"

	"github.com/gocaio/akt/modules/rotate"
	"github.com/gocaio/akt/modules/setenv"
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

	if rotateFlag && !setenvFlag {
		println("rotate ", rotateFlag)
		rotate.GetNewPair()
	}
	if setenvFlag && !rotateFlag {
		println("setenv ", setenvFlag)
		setenv.GetVars()
	}

	if setenvFlag && rotateFlag {
		println("rotate ", rotateFlag)
		rotate.GetNewPair()
		println("setenv ", setenvFlag)
		setenv.GetVars()
	}
}
