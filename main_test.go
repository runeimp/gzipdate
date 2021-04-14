package main

import "testing"

func reset() {
	deleteSource = false
	files = []string{}
	useFileDate = false
	useTimeZone = false
}

func TestPosixOptionGroup1(t *testing.T) {
	reset()

	args := []string{"-dt", "test-file.txt"}
	if deleteSource && useTimeZone {
		t.Errorf("deleteSource and useTimeZone already set to true")
	} else if deleteSource {
		t.Errorf("deleteSource already set to true")
	} else if useTimeZone {
		t.Errorf("useTimeZone already set to true")
	}

	argParse(args)

	if deleteSource == false && useTimeZone == false {
		t.Errorf("deleteSource and useTimeZone were not set to true")
	} else if deleteSource == false {
		t.Errorf("deleteSource was not set to true")
	} else if useTimeZone == false {
		t.Errorf("useTimeZone was not set to true")
	}

}

func TestPosixOptionGroup2(t *testing.T) {
	reset()

	args := []string{"-td", "test-file.txt"}
	argParse(args)

	if deleteSource == false && useTimeZone == false {
		t.Errorf("deleteSource and useTimeZone were not set to true")
	} else if deleteSource == false {
		t.Errorf("deleteSource was not set to true")
	} else if useTimeZone == false {
		t.Errorf("useTimeZone was not set to true")
	}
}

func TestPosixOptionGroup3(t *testing.T) {
	reset()

	args := []string{"-tx", "test-file.txt"}
	expect := "unknown option -x"
	err := argParse(args)
	got := ""
	if err != nil {
		got = err.Error()
	}

	if got != expect {
		t.Errorf("expected %q but got %q", expect, got)
	}
}

func TestPosixOptionGroup4(t *testing.T) {
	reset()

	args := []string{"-tv", "test-file.txt"}
	expect := "invalid use of -v in a POSIX group"
	err := argParse(args)
	got := ""
	if err != nil {
		got = err.Error()
	}

	if got != expect {
		t.Errorf("expected %q but got %q", expect, got)
	}
}

func TestOptionParsing(t *testing.T) {
	reset()

	args := []string{"-t", "test-file.txt", "--", "--version"}
	expect := []string{"test-file.txt", "--version"}
	argParse(args)

	for i, got := range files {
		if i < len(expect) && got != expect[i] {
			t.Errorf("expected [%d] %q but got %q", i, expect[i], got)
		}
	}
}

func TestFilenameGeneration(t *testing.T) {
	reset()

	filename := "test-file.txt"

	args := []string{"-ft", filename}
	argParse(args)
	expect := "test-file.txt_2021-04-13_133403_PDT.gz"
	_, modTime := fileContentAndModTime(filename)
	got := datetimeFilename(filename, modTime)

	if got != expect {
		t.Errorf("expected %q but got %q", expect, got)
	}
}
