package main

import (
	"encoding/json"
	"testing"
)

// Note: First column is postion 0 but offset is 1 based

func TestColAndPos(t *testing.T) {
	tests := []struct{
		name string
		input string
		offset int64
		wantLine, wantPos int
	}{
		{
			"TestErrorFirstLine", "123\n5678\n012345", 2, 1, 1,
		},
		{
			"TestErrorMiddleLine", "123\n5678\n012345", 7, 2, 2,
		},
		{
			"TestErrorLastLine", "123\n5678\n012345", 14, 3, 4,
		},
		{
			"TestWindowsLineEndings", "123\r\n6789\r\n234567", 14, 3, 2,
		},
		{
			"TestOffsetZero", "123\n5678\n012345", 0, 0, 0,
		},
		{
			"TestOffsetFstNewLine", "123\n568\n012345", 4, 2, 0,
		},
		{
			"TestOffsetSndNewLine", "123\n568\n346345", 9, 3, 0,
		},
		{
			"TestOffsetLastNewLine", "123\n5678\n012345\n", 16, 4, 0,
		},
		{
			"TestOffsetLastChar", "123\n5678\n012345", 15, 3, 5,
		},
		{
			"TestOffsetLastCharNewLine", "123\n5678\n012345\n", 15, 3, 5,
		},
		{
			"TestOffsetNegative", "123\n5678\n012345", -47, -1, -1,
		},
		{
			"TestOffsetLastTooBig", "123\n5678\n012345", 150, -1, -1,
		},
	}

	for _,test := range tests {
		err := json.SyntaxError{Offset: test.offset}
		line, pos, _ := findLineAndPos(&err, test.input)
		if line != test.wantLine {
			t.Errorf("%v - Wrong line number: %v; Expected: %v", test.name, line, test.wantLine)
		}
		if pos != test.wantPos {
			t.Errorf("%v - Wrong position number: %v; Expected: %v", test.name, pos, test.wantPos)
		}
	}
}

func TestCorrectLine(t *testing.T) {
		err := json.SyntaxError{Offset: 36}
		js := "{\n" +
					"    \"src/github.com/golang/test\" {\n" + // Missing colon
					"        \"vcs\": \"git\"\n" +
					"    }\n" +
					"}"
		_, _, ex := findLineAndPos(&err, js)
		correct := "\"src/github.com/golang/test\" {"
		if ex != correct {
			t.Errorf("TestCorrectLine - Wrong printed line: %v; Expected: %v", ex, correct)
		}
}

func TestCorrectErrorMessage(t *testing.T) {
	err := json.SyntaxError{Offset: 65}
	js := "{\n" +
				"    \"src/github.com/golang/test\": {\n" +
				"        \"vcs\": \"git\"\n" +
				"    }\n" // No closing bracket
	enrErr := enrichJSONError(&err, js)
	// Error message would normally be before the new line
	correct := "\nOccurred on line 5 at pos 0: "
	if enrErr.Error() != correct {
		t.Errorf("TestCorrectErrorMessage - Error was: %v\nError should be: %v", enrErr.Error(), correct)
	}
}
