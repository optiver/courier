package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

type createDirHashTest struct {
	dir    string
	ignore string
}

var testsExpectedDifferent = []createDirHashTest{
	//Folder with one file (empty)
	{"./test/filesystem/test_single_file", ""},
	//Folder with one file (not empty)
	{"./test/filesystem/test_single_file_nonempty", ""},
	//Folder with multiple files
	{"./test/filesystem/test_multiple_files", ""},
	//Test presence of subfolder (one file)
	{"./test/filesystem/test_subfolder", ""},
	//Test the ignore feature (pt1)
	{"./test/filesystem/test_ignored_folder", ""},
	//Test the ignore feature (pt2)
	{"./test/filesystem/test_ignored_folder", "ignored"},
	//Test the ignore feature where ignored folder is in subfolder (pt1)
	{"./test/filesystem/test_ignored_subfolder", ""},
	//Test the ignore feature where ignored folder is in subfolder (pt2)
	{"./test/filesystem/test_ignored_subfolder", "ignored"},
}

var testsExpectedSame = [][]createDirHashTest{
	{ //Test different paths to same folders have the same hash
		{"./test/filesystem/test_different_paths/a", ""},
		{"./test/filesystem/test_different_paths/b", ""},
	},
	{ //Test changes inside ignored folder
		{"./test/filesystem/test_ignored_folder_changes/a", "ignored"},
		{"./test/filesystem/test_ignored_folder_changes/b", "ignored"},
	},
}

func TestCreateDirHashNonExistent(t *testing.T) {
	if _, err := CreateDirHash("some non existent file", ""); !os.IsNotExist(err) {
		t.Errorf("CreateDirHash: Expected error on non-existent file")
	}
}

func TestCreateDirHashExpectDifferent(t *testing.T) {
	seen := make(map[string]string)
	for _, test := range testsExpectedDifferent {
		hash, err := CreateDirHash(test.dir, test.ignore)
		if err != nil {
			t.Errorf("CreateDirHash: %q [%q]: Error: %v", test.dir, test.ignore, err)
		}
		strHash := fmt.Sprintf("%x", hash)

		if _, present := seen[strHash]; present {
			t.Errorf("CreateDirHash: %q [%q]: Got same hash as %q (%q), expected a difference",
				test.dir, test.ignore, seen[strHash], strHash)
		}
		seen[strHash] = test.dir
	}
}

func TestCreateDirHashExpectSame(t *testing.T) {
	for _, tests := range testsExpectedSame {
		var expected []byte
		for _, test := range tests {
			hash, err := CreateDirHash(test.dir, test.ignore)
			if err != nil {
				t.Errorf("CreateDirHash: %q [%q]: Error: %v", test.dir, test.ignore, err)
			}

			if expected == nil {
				expected = hash
			} else if !bytes.Equal(hash, expected) {
				t.Errorf("CreateDirHash: %q [%q]: Got hash %x, expected same as %x",
					test.dir, test.ignore, hash, expected)
			}
		}
	}
}

func TestCreateDirHashPermissions(t *testing.T) {
	infoNoX, err := os.Stat("./test/filesystem/test-x/file1")
	if err != nil {
		t.Errorf("CreateDirHash: %v", err)
	}
	infoWithX, err := os.Stat("./test/filesystem/test+x/file1")
	if err != nil {
		t.Errorf("CreateDirHash: %v", err)
	}

	//Check for the execution bit
	if infoNoX.Mode() == infoWithX.Mode() {
		t.Logf("Detected that execution bits are not supported; skipping test.")
	} else {
		hashNoX, err := CreateDirHash("./test/filesystem/test-x", "")
		if err != nil {
			t.Errorf("CreateDirHash: %q: Error: %v", "test-x", err)
		}
		hashWithX, err := CreateDirHash("./test/filesystem/test+x", "")
		if err != nil {
			t.Errorf("CreateDirHash: %q: Error: %v", "test+x", err)
		}
		if bytes.Equal(hashNoX, hashWithX) {
			t.Errorf("CreateDirHash: Expected differing hashes due to change in permissions bit")
		}
	}
}
