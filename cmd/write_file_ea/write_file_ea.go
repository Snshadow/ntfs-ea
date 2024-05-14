//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Snshadow/ntfs-ea"
)

func main() {
	var srcPath, targetPath, eaName string
	var needEa bool

	flag.StringVar(&targetPath, "target-path", "", "path of target file to write EA")
	flag.StringVar(&srcPath, "source-path", "", "path of source file to be used as content for EA")
	flag.StringVar(&eaName, "ea-name", "", "name of the EA")

	flag.BoolVar(&needEa, "need-ea", false, "set flag if file needs to be interpreted with EA")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s writes EA(Extended Attribute) info a file in NTFS(New Technology File System) with the content of a given source file, if the source file is empty the EA with ea-name is removed if exists.\nUsage: %s [target file] [source file] [ea name]\n or\n %s -target-path [target path] -source-path [source path] -ea-name [ea-name] -need-ea\nThis program only works in Windows.\n\n", os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if targetPath == "" {
		targetPath = flag.Arg(0)
	}
	if srcPath == "" {
		srcPath = flag.Arg(1)
	}
	if eaName == "" {
		eaName = flag.Arg(2)
	}

	var flags uint8 = 0
	if needEa {
		flags |= ntfs_ea.NeedEa
	}

	if flag.NArg() < 3 {
		flag.Usage()
		os.Exit(1)
	}

	err := ntfs_ea.AddEaFileWithFile(targetPath, srcPath, eaName, flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write EA into file: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Written EA into file \"%s\" using \"%s\" with ea name \"%s\"\n", targetPath, srcPath, eaName)
}
