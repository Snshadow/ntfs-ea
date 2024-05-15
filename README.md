# ntfs-ea
Access NTFS(new technology file system) Extended Attributes(EA) with golang.

This package provides functions for writing and query Extended Attributes for files in NTFS which can be shown by using "fsutil file queryea [file_path] in cmd in Windows."

_example(cp949)_
~~~
C:\test>fsutil file queryea test.txt

Extended Attributes (EA) information for file C:\test\test.txt:

Total Ea Size: 0xbd

Ea Buffer Offset: 0
Ea Name: TEST
Ea Value Length: 16
0000:  74 68 69 73 20 69 73 20  61 6e 20 45 41 20 63 6f  this is an EA co
0010:  6e 74 65 6e 74 2e                                 ntent.

Ea Buffer Offset: 24
Ea Name: テスト
Ea Value Length: 21
0000:  45 41 20 77 72 69 74 65  20 74 65 73 74 2e 20 e3  EA write test. .
0010:  83 86 e3 82 b9 e3 83 88  20 6a 61 70 61 6e 65 73  ........ japanes
0020:  65                                                e

Ea Buffer Offset: 54
Ea Name: 테스트
Ea Value Length: 1f
0000:  45 41 20 77 72 69 74 65  20 74 65 73 74 2e 20 ed  EA write test. .
0010:  85 8c ec 8a a4 ed 8a b8  20 6b 6f 72 65 61 6e     ........ korean

Ea Buffer Offset: 84
Ea Name: 試驗
Ea Value Length: 2c
0000:  45 41 20 77 72 69 74 65  20 74 65 73 74 2e 20 e8  EA write test. .
0010:  a9 a6 e9 a9 97 20 43 4a  4b 20 55 6e 69 66 69 65  ..... CJK Unifie
0020:  64 20 49 64 65 6f 67 72  61 70 68 73              d Ideographs
~~~

## Writing EA
For writing EA into a file, EaAddFile or WriteEaWithFile can be used to add EA.

_write EA with byte slice_
```go
import (
	"github.com/Snshadow/ntfs-ea"
)

func main() {
	eaVal := "this is a value for ea"
	targetPath := "C:\\test\\test.txt"
	
	ea := ntfs_ea.EaInfo{
		Flags: 0,
		EaName: "eaFromBytes",
		EaValues: []byte(eaVal),
	}
	
	err := ntfs_ea.EaWriteFile(targetPath, ea)
	if err != nil {
		panic(err)
	}	
}
```

_write EA with file content_
```go
import (
	"github.com/Snshadow/ntfs-ea"
)

func main() {
	targetPath := "C:\\test\\test.txt"
	sourcePath := "C:\\test\\source.txt"
	flags := 0
	
	err := ntfs_ea.WriteEaWithFile(targetPath, sourcePath, flags, "eaFromFile")
	if err != nil {
		panic(err)
	}
}
```

## Querying EA
For querying EA within the file. QueryEaFile can be used with target file path.

_querying ea from file and writing it to console_
```go
import (
	"encoding/hex"
	"fmt"
	
	"github.com/Snshadow/ntfs-ea"
)

func main() {
    targetPath := "C:\\test\\test.txt"
	
	eaList, err := ntfs_ea.QueryEaFile(targetPath)
	if err != nil {
		panic err
	}
	
	for _, ea := range eaList {
		fmt.Printf("Flags: 0x%x\nEa Name: %s\nEa Value\n:%s\n", ea.Flags, ea.EaName, hex.Dump(ea.EaValue))
	}
}
```

## Executables
This package has two executables for accessing EA from file. Source files can be found in cmd directory.

※ _"github.com/josephspurrier/goversioninfo/cmd/goversioninfo" package used to set information for exe._

_usage_
~~~
write_file_ea.exe writes EA(Extended Attribute) info a file in NTFS(New Technology File System) with the content of a given source file, if the source file is empty the EA with ea-name is removed if exists.
Usage: write_file_ea.exe [target file] [source file] [ea name]
 or
 write_file_ea.exe -target-path [target path] -source-path [source path] -ea-name [ea-name] -need-ea
This program only works in Windows.

  -ea-name string
        name of the EA
  -need-ea
        set flag if file needs to be interpreted with EA
  -source-path string
        path of source file to be used as content for EA
  -target-path string
        path of target file to write EA
~~~
~~~
query_file_ea.exe queries EA(Extended Attribute) from a file in NTFS(New Technology File System).
Usage: query_file_ea.exe -query-name [eaName1],[eaName2],... -extract [target path]
 or  query_file_ea.exe -target-path [target path] -query-name [eaName1],[eaName2],... -dump -extract
This program is supposed to work only in Windows.

  -dump
        dump EA to console, this is enabled by default if no action is given
  -extract
        extract EA to file(s) with according EaName
  -query-name string
        names of EA to query, split by comma, if it is not given, query all EA from the file
  -target-path string
        path of the file to query EA
~~~
