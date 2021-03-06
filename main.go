package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joanbono/akt/modules/rotate"
	"github.com/joanbono/akt/modules/writer"
)

// Defining colors
var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var cyan = color.New(color.FgCyan)
var bold = color.New(color.FgHiWhite, color.Bold)

var (
	accessKey, secretKey, username string
	profileFlag, userFlag          string
	rotateFlag                     bool
	saveFlag                       bool
	versionFlag                    bool
	version                        string
)

func init() {
	flag.StringVar(&profileFlag, "profile", "default", "Profile to use")
	flag.StringVar(&userFlag, "user", "", "User to generate keys")
	flag.BoolVar(&rotateFlag, "rotate", false, "Rotate keys")
	flag.BoolVar(&saveFlag, "save", false, "Save new keys to .aws/credentials")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()
}

func main() {
	if versionFlag {
		fmt.Printf("\n\tâï¸  ğ AKT %v\n\n", bold.Sprintf(version))
		return
	}

	if !rotateFlag {
		fmt.Printf("\n âï¸  ğ AKT %v\n", bold.Sprintf(version))
		fmt.Printf("%v Insufficient options!\n", yellow.Sprintf("[!]"))
		fmt.Printf("%v Try with %v\n\n", cyan.Sprintf("[i]"), bold.Sprintf("akt -h"))
		return
	} else {
		//Check that .aws/credentials file is accessible before rotating the keys
		//This will prevent rotated keys not being written anywhere
		if _, err := os.Stat(writer.Reader()); err != nil {
			fmt.Printf("%v .aws/credentials file not found\n", red.Sprintf("[-]"))
			os.Exit(2)
		}
		accessKey, secretKey, username = rotate.Rotate(profileFlag, userFlag)
		if saveFlag {
			writer.Profiler(profileFlag, accessKey, secretKey)
		} else {
			writer.Printer(profileFlag, accessKey, secretKey)
		}
	}
}
