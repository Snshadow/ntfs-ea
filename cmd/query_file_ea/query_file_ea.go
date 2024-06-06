//go:build windows
// +build windows

//go:generate go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo query_file_ea.json

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Snshadow/ntfs-ea"
)

func main() {
	var targetPath, queryName string
	var dump, extract, stdout bool

	flag.StringVar(&targetPath, "target-path", "", "path of the file to query EA")
	flag.StringVar(&queryName, "query-name", "", "names of EA to query, split by comma, if it is not given, query all EA from the file")

	flag.BoolVar(&dump, "dump", false, "dump EA to console, this is enabled by default if no action is given")
	flag.BoolVar(&extract, "extract", false, "extract EA to file(s) with according EaName")
	flag.BoolVar(&stdout, "stdout", false, "extract EA into stdout")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s queries EA(Extended Attribute) from a file in NTFS(New Technology File System).\nUsage: %s -query-name [eaName1],[eaName2],... -extract [target path]\n or %s -target-path [target path] -query-name [eaName1],[eaName2],... -dump -extract\nWrite EA value to stdout(for piping output): %s -stdout -extract -target-path [target path] -query-name [eaName] | (process output)\n\n", os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	var queryList []string
	if len(queryName) != 0 {
		queryList = strings.Split(queryName, ",")
	}

	if !(dump || extract) {
		dump = true
	}

	if targetPath == "" {
		targetPath = flag.Arg(0)
		if targetPath == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	eaList, err := ntfs_ea.QueryFileEa(targetPath, queryList...)
	if err != nil {
		fmt.Printf("Error querying EA: %v\n", err)
		os.Exit(2)
	}

	for _, ea := range eaList {
		// If EaName with missing EA is queried, the EA with that name is returned with empty EaValue instead of not returning EA with the specific name.
		if len(ea.EaValue) == 0 {
			fmt.Fprintf(os.Stderr, "EA with name \"%s\" does not exist\n", ea.EaName)
			continue
		}

		if dump {
			fmt.Printf("Flags: 0x%x\nEa Name: %s\nEa Value:\n%s\n", ea.Flags, ea.EaName, hex.Dump(ea.EaValue))
		}

		if extract {
			if stdout {
				os.Stdout.Write(ea.EaValue)
			} else {
				err := os.WriteFile(ea.EaName, ea.EaValue, 0777)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to write EA for %s into file: %v\n", ea.EaName, err)
				} else {
					fmt.Printf("Extracted EaValue in \"%s\"", ea.EaName)
				}
			}
		}
	}
}
