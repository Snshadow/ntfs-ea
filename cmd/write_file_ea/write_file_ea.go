//go:build windows
// +build windows

//go:generate go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo write_file_ea.json

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Snshadow/ntfs-ea"
	"github.com/Snshadow/ntfs-ea/cmd/utils"
)

func main() {
	var srcPath, targetPath, eaName string
	var needEa, removeEa, stdin bool

	flag.StringVar(&targetPath, "target-path", "", "path of target file to write EA")
	flag.StringVar(&srcPath, "source-path", "", "path of source file to be used as content for EA")
	flag.StringVar(&eaName, "ea-name", "", "name of the EA")

	flag.BoolVar(&needEa, "need-ea", false, "set flag if file needs to be interpreted with EA")
	flag.BoolVar(&removeEa, "remove-ea", false, "remove the EA with the given name")
	flag.BoolVar(&stdin, "stdin", false, "use stdin as content for EA")

	progName := filepath.Base(os.Args[0])

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s writes EA(Extended Attribute) info a file in NTFS(New Technology File System) with the content of a given source file, if the source file is empty the EA with EaName is removed if exists.\nUsage: %s [target path] [source path] [EA name]\n or\n %s -target-path [target path] -source-path [source path] -ea-name [EA name] -need-ea\nWrite EA from stdin: echo \"[content for ea] | %s -stdin -target-path [target path] -ea-name [EA name]\nTo remove EA with specific name, use: %s -remove-ea [target path] [EA name]\n\n", progName, progName, progName, progName, progName)
		flag.PrintDefaults()

		// prevent window from closing immediately if the console was created for this process
		if utils.IsFromOwnConsole() {
			fmt.Println("\nPress enter to close...")
			fmt.Scanln()
		}
	}

	flag.Parse()

	requiredArg := 3

	if targetPath == "" {
		targetPath = flag.Arg(0)
	}
	if srcPath == "" && !removeEa {
		srcPath = flag.Arg(1)
	}
	if eaName == "" {
		if removeEa || stdin {
			eaName = flag.Arg(1)
			requiredArg = 2
		} else {
			eaName = flag.Arg(2)
		}
	}

	var flags uint8 = 0
	if needEa {
		flags |= ntfs_ea.NeedEa
	}

	if !(targetPath != "" && srcPath != "" && eaName != "") && !((removeEa || stdin) && targetPath != "" && eaName != "") && flag.NArg() < requiredArg {
		flag.Usage()
		os.Exit(1)
	}

	if removeEa {
		eaToRemove := ntfs_ea.EaInfo{
			EaName: eaName,
		}

		err := ntfs_ea.EaWriteFile(targetPath, eaToRemove)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove EA with name \"%s\" from file: %v\n", eaName, err)
			os.Exit(2)
		}

		fmt.Printf("Removed EA with name \"%s\" from file\n", eaName)

		return
	}

	if stdin {
		eaBuf := make([]byte, 65536)

		n, err := os.Stdin.Read(eaBuf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read from stdin:", err)
			os.Exit(2)
		}

		if n < len(eaBuf) {
			eaBuf = eaBuf[:n]
		}

		eaToWrite := ntfs_ea.EaInfo{
			Flags:   flags,
			EaName:  eaName,
			EaValue: eaBuf,
		}

		err = ntfs_ea.EaWriteFile(targetPath, eaToWrite)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write EA into file:", err)
			os.Exit(2)
		}

		fmt.Printf("Written EA info file \"%s\" from stdin with ea name \"%s\"\n", targetPath, eaName)

		return
	}

	err := ntfs_ea.WriteEaWithFile(targetPath, srcPath, eaName, flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write EA into file:", err)
		os.Exit(2)
	}

	fmt.Printf("Written EA into file \"%s\" using \"%s\" with ea name \"%s\"\n", targetPath, srcPath, eaName)
}
