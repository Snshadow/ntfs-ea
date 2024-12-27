//go:build windows
// +build windows

package ntfs_ea

import (
	"os"
	"path/filepath"
	"testing"
)

var tempDir string

func TestMain(m *testing.M) {
	var err error
	tempDir, err = os.MkdirTemp("", "ntfs_ea_test")
	if err != nil {
		panic("Failed to create temp directory")
	}

	code := m.Run()

	os.RemoveAll(tempDir)

	os.Exit(code)
}

func TestWriteFileEa(t *testing.T) {
	testFile := filepath.Join(tempDir, "writetest.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	eaInfos := []EaInfo{
		{EaName: "TESTEA¼", EaValue: []byte("test value 1")},
		{EaName: "TESTEA2", EaValue: []byte("test value 2")},
		{EaName: "TESTEA³", EaValue: []byte("test value 3")},
		{Flags: NeedEa, EaName: "TEST±A4", EaValue: []byte("test required value 4")},
	}

	err = EaWriteFile(testFile, false, eaInfos...)
	if err != nil {
		t.Fatalf("EaWriteFile failed: %v", err)
	}

	eas, err := QueryFileEa(testFile, false)
	if err != nil {
		t.Fatalf("QueryFileEa failed: %v", err)
	}

	if len(eas) != len(eaInfos) {
		t.Fatalf("Expected %d EAs, got %d", len(eaInfos), len(eas))
	}

	for i, eaInfo := range eaInfos {
		if eas[i].EaName != eaInfo.EaName || string(eas[i].EaValue) != string(eaInfo.EaValue) {
			t.Fatalf("EA data mismatch: got %v, expected %v", eas[i], eaInfo)
		}
	}
}

func TestWriteEaWithFile(t *testing.T) {
	testFile := filepath.Join(tempDir, "writefile.txt")
	srcFile := filepath.Join(tempDir, "srcfile.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = os.WriteFile(srcFile, []byte("src content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create src file: %v", err)
	}

	err = WriteEaWithFile(testFile, false, srcFile, "FILEEA", 0)
	if err != nil {
		t.Fatalf("WriteEaWithFile failed: %v", err)
	}

	eas, err := QueryFileEa(testFile, false)
	if err != nil {
		t.Fatalf("QueryFileEa failed: %v", err)
	}

	if len(eas) != 1 {
		t.Fatalf("Expected 1 EA, got %d", len(eas))
	}

	if eas[0].EaName != "FILEEA" || string(eas[0].EaValue) != "src content" {
		t.Fatalf("EA data mismatch: got %v", eas[0])
	}
}

func TestQueryFileEa(t *testing.T) {
	testFile := filepath.Join(tempDir, "testquery.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	eaInfo := []EaInfo{
		{Flags: 0, EaName: "TESTEA¼", EaValue: []byte("test value 1")},
		{EaName: "TESTEA2", EaValue: []byte("test value 2")},
	}

	err = EaWriteFile(testFile, false, eaInfo...)
	if err != nil {
		t.Fatalf("EaWriteFile failed: %v", err)
	}

	eas, err := QueryFileEa(testFile, false, "TESTEA2")
	if err != nil {
		t.Fatalf("QueryFileEa failed: %v", err)
	}

	if len(eas) != 1 {
		t.Fatalf("Expected 1 EA, got %d", len(eas))
	}

	if eas[0].EaName != "TESTEA2" || string(eas[0].EaValue) != "test value 2" {
		t.Fatalf("EA data mismatch: got %v", eas[0])
	}
}
