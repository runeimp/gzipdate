package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

/*
 * CONSTANTS
 */
const (
	AppName         = "GzipDate"
	AppVersion      = "1.0.0"
	argDeleteUsage  = "Delete the source file after successful processing"
	argHelpUsage    = "Display this help info"
	argVersionUsage = "Display this apps version number"
)

/*
 * DERIVED CONSTANTS
 */
var (
	AppLabel = fmt.Sprintf("%s v%s", AppName, AppVersion)
)

/*
 * VARIABLES
 */
var (
	deleteSource = false
	files        []string
	showHelp     = false
)

func main() {

	// Bespoke CLI Handler because flags is lame
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-d", "-del", "-delete":
			deleteSource = true
		case "-h", "-help":
			helpOutput()
			os.Exit(0)
		case "-v", "-ver", "-version":
			fmt.Println(AppLabel)
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}

	if len(os.Args[1:]) == 0 {
		helpOutput()
		os.Exit(1)
	}

	// fmt.Printf("Delete Source File? %t\n", deleteSource)

	var err error

	if len(files) > 0 {
		for _, filename := range files {
			// j := i + 1
			// fmt.Printf("%-3d %s\n", j, filename)
			if path.Ext(filename) == ".gz" {
				// fmt.Printf("File #%-3d gunzip %q\n", j, filename)
				content, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				err = gzipDecode(content)
				if err == nil && deleteSource {
					fmt.Printf("    Deleting source: %q\n", filename)
					os.Remove(filename)
				}
			} else {
				// fmt.Printf("File #%-3d gzip %q\n", j, filename)
				err = gzipEncode(filename)
				if err == nil && deleteSource {
					fmt.Printf("    Deleting source: %q\n", filename)
					os.Remove(filename)
				}
			}
		}
	}
}

func gzipDecode(content []byte) error {
	newFilename := "filename.unknown"

	zr, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())
	// if _, err := io.Copy(os.Stdout, zr); err != nil {
	// 	log.Fatal(err)
	// }

	if len(zr.Name) > 0 {
		newFilename = zr.Name
	}
	destination, err := os.Create(newFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, zr)
	if err != nil {
		log.Fatal(err)
	}

	mtime := zr.ModTime.UTC()
	atime := zr.ModTime.UTC()
	if err := os.Chtimes(newFilename, atime, mtime); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d B written to %q\n", nBytes, zr.Name)

	return err
}

func gzipEncode(filename string) error {
	var (
		destination *os.File
		err         error
		file        *os.File
		fileinfo    os.FileInfo
		nBytes      int64
	)
	// content, err := ioutil.ReadFile(filename)
	file, err = os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	fileinfo, err = file.Stat()
	content := make([]byte, fileinfo.Size())
	_, err = file.Read(content)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	datetime := now.Format("2006-01-02_150405")
	newFilename := fmt.Sprintf("%s_%s.gz", filename, datetime)
	// fmt.Printf("gzipEncode() | filename = %q | content length = %d B | newFilename = %q\n", filename, len(content), newFilename)

	var buf bytes.Buffer
	zw, zwError := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if zwError != nil {
		log.Fatal(zwError)
	}

	// Setting the Header fields is optional.
	zw.Name = filename
	zw.Comment = fmt.Sprintf("Compressed with %s", AppLabel)
	zw.ModTime = fileinfo.ModTime()

	if _, err := zw.Write(content); err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	destination, err = os.Create(newFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	nBytes, err = buf.WriteTo(destination)
	fmt.Printf("Saving %d bytes from %q to %q\n", nBytes, filename, newFilename)

	return err
}

func helpOutput() {
	// flag.Usage()
	// fmt.Println("----")
	fmt.Printf("USAGE: gzipdate [OPTIONS] FILENAMES\n\nOPTIONS:\n")
	// flag.PrintDefaults()
	// fmt.Println()

	fmt.Printf("  -d | -del | -delete   %s\n", argDeleteUsage)
	fmt.Printf("  -h | -help            %s\n", argHelpUsage)
	fmt.Printf("  -v | -ver | -version  %s\n", argVersionUsage)
	fmt.Println()
	fmt.Printf("Options may be interspersed with file names if so desired.\nThey are not position dependent.\n\n")
}
