//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type queryList []string

func (q queryList) String() string {
	return strings.Join(q, ",")
}

func (q *queryList) Set(value string) error {
	*q = strings.Split(value, ",")
	return nil
}

func main() {
	var targetPath string
	var queryName queryList
	var dump, extract bool

	flag.StringVar(&targetPath, "target-path", "", "path of the file to query EA")
	flag.Var(&queryName, "query-name", "names of EA to query, split by comma, if it is not given, query all EA from the file")
	flag.BoolVar(&dump, "dump", false, "dump EA to console, this is the default action")
	flag.BoolVar(&extract, "extract", false, "extract EA to file(s) with according EaName")

	flag.Parse()

	if !(dump || extract) {
		fmt.Println("action isn't specified")
	}

	if targetPath == "" {
		fmt.Println("target path is required")
		flag.Usage()
		os.Exit(1)
	}

}
