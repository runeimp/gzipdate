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
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

//
// CONSTANTS
//
const (
	AppName              = "GzipDate"
	AppVersion           = "1.1.0"
	cliUsageDelete       = "Delete the source file after processing"
	cliUsageDoubleHyphen = "Disable option parsing and consider all following arguments as file names only"
	cliUsageFileDate     = "Use the files modification time for the date instead of the current time"
	cliUsageHelp         = "Display this help info"
	cliUsageTimeZone     = "Turn the timezone feature on"
	cliUsageVersion      = "Display this apps version number"
	cliUsageHeader       = `%s

USAGE: gzipdate [OPTIONS] [FILENAMES]

OPTIONS:
`
)

//
// DERIVED CONSTANTS
//
var (
	AppLabel = fmt.Sprintf("%s v%s", AppName, AppVersion)
)

//
// VARIABLES
//
var (
	deleteSource = false
	files        []string
	showHelp     = false
	termWidth    int
	timeFormat   = "2006-01-02_150405"
	useFileDate  = false
	useTimeZone  = false
)

//
// MAIN ENTRY POINT
//
func main() {
	var err error

	termWidth, _, _ = terminal.GetSize(int(os.Stdin.Fd()))
	if termWidth == 0 {
		termWidth = 80
	}

	if len(os.Args[1:]) == 0 {
		helpOutput()
		os.Exit(1)
	}

	err = argParse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	gzdEnv := strings.ToUpper(os.Getenv("GZIPDATE"))
	if gzdEnv == "TIMEZONE" || gzdEnv == "TZ" {
		useTimeZone = true
	}

	// fmt.Printf("Delete Source File? %t\n", deleteSource)

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

//
// FUNCTIONS
//
func argParse(args []string) error {
	// Bespoke CLI Handler because flags is lame
	parseArgs := true
	for _, arg := range args {
		if parseArgs {
			switch arg {
			case "--":
				parseArgs = false
			case "-d", "-del", "-delete", "--delete":
				deleteSource = true
			case "-f", "-file", "-file-date", "--file-date":
				useFileDate = true
			case "-h", "-help", "--help":
				helpOutput()
				os.Exit(0)
			case "-t", "-tz", "-timezone", "--timezone":
				useTimeZone = true
			case "-v", "-ver", "-version", "--version":
				fmt.Println(AppLabel)
				os.Exit(0)
			default:
				if len(arg) > 0 {
					if arg[0:1] == "-" {
						// POSIX group processing
						for _, c := range arg[1:] {
							switch c {
							case 'd':
								deleteSource = true
							case 'f':
								useFileDate = true
							case 't':
								useTimeZone = true
							case 'h', 'v':
								return fmt.Errorf("invalid use of -%s in a POSIX group", string(c))
								// fmt.Fprintf(os.Stderr, "Invalid use of -%s in a POSIX group\n", string(c))
								// os.Exit(1)
							default:
								return fmt.Errorf("unknown option -%s", string(c))
								// fmt.Fprintf(os.Stderr, "Unknown POSIX group option: %q\n", c)
								// fmt.Fprintf(os.Stderr, "Unknown option -%s\n", string(c))
								// os.Exit(1)
							}
						}
					} else {
						// Arg files
						files = append(files, arg)
					}
				}
			}
		} else {
			files = append(files, arg)
		}
	}

	return nil
}

func datetimeFilename(filename string, t time.Time) string {
	if useTimeZone {
		timeFormat = "2006-01-02_150405_MST"
	}
	datetime := t.Format(timeFormat)
	return fmt.Sprintf("%s_%s.gz", filename, datetime)
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

func fileContentAndModTime(filename string) ([]byte, time.Time) {
	var (
		err      error
		file     *os.File
		fileinfo os.FileInfo
	)
	file, err = os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	fileinfo, err = file.Stat()
	content := make([]byte, fileinfo.Size())
	_, err = file.Read(content)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	return content, fileinfo.ModTime()
}

func gzipEncode(filename string) error {
	var (
		destination *os.File
		err         error
		modTime     time.Time
		nBytes      int64
	)
	// content, err := ioutil.ReadFile(filename)
	// file, err = os.Open(filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fileinfo, err = file.Stat()
	// content := make([]byte, fileinfo.Size())
	// _, err = file.Read(content)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	content, modTime := fileContentAndModTime(filename)

	now := time.Now()
	if useFileDate {
		now = modTime
	}
	newFilename := datetimeFilename(filename, now)
	// if useTimeZone {
	// 	timeFormat = "2006-01-02_150405_MST"
	// }
	// datetime := now.Format(timeFormat)
	// newFilename := fmt.Sprintf("%s_%s.gz", filename, datetime)
	// fmt.Printf("gzipEncode() | filename = %q | content length = %d B | newFilename = %q\n", filename, len(content), newFilename)

	var buf bytes.Buffer
	zw, zwError := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if zwError != nil {
		log.Fatal(zwError)
	}

	// Setting the Header fields is optional.
	zw.Name = filename
	zw.Comment = fmt.Sprintf("Compressed with %s", AppLabel)
	zw.ModTime = modTime

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

func helpWrap(msg, prefix string) string {
	if len([]rune(msg)) > termWidth {
		words := strings.Split(msg, " ")
		lines := []string{""}
		i := 0
		for _, word := range words {
			// fmt.Fprintf(os.Stderr, "helpWrap() | i: %d\n", i)
			// fmt.Fprintf(os.Stderr, "helpWrap() | word: %q\n", word)
			tmpMsg := lines[i] + word + " "
			if len([]rune(tmpMsg)) < termWidth {
				lines[i] = tmpMsg
			} else {
				last := len(lines[i]) - 1
				lines[i] = lines[i][:last]
				i++
				word = fmt.Sprintf("%s%s ", prefix, word)
				lines = append(lines, word)
			}
		}
		last := len(lines[i]) - 1
		lines[i] = lines[i][:last]
		// fmt.Fprintf(os.Stderr, "helpWrap() | len(lines): %d\n", len(lines))
		msg = strings.Join(lines, "\n")
	}
	msg += "\n"

	return msg
}

func helpOutput() {
	// fmt.Fprintf(os.Stderr, "helpOutput() | termWidth: %d\n", termWidth)
	fmt.Printf(cliUsageHeader, AppLabel)
	// fmt.Printf("  -d | -del  | -delete     %s\n", cliUsageDelete)
	// fmt.Printf("  -f | -file | -file-date  %s\n", cliUsageFileDate)
	// fmt.Printf("  -h | -help | -help       %s\n", cliUsageHelp)
	// fmt.Printf("  -t | -tz   | -timezone   %s\n", cliUsageTimeZone)
	// fmt.Printf("  -v | -ver  | -version    %s\n", cliUsageVersion)
	prefix := "                            "
	msg := helpWrap(fmt.Sprintf("  -d | -del  | --delete     %s", cliUsageDelete), prefix)
	msg += helpWrap(fmt.Sprintf("  -f | -file | --file-date  %s", cliUsageFileDate), prefix)
	msg += helpWrap(fmt.Sprintf("  -h | -help | --help       %s", cliUsageHelp), prefix)
	msg += helpWrap(fmt.Sprintf("  -t | -tz   | --timezone   %s", cliUsageTimeZone), prefix)
	msg += helpWrap(fmt.Sprintf("  -v | -ver  | --version    %s", cliUsageVersion), prefix)
	msg += helpWrap(fmt.Sprintf("  --                        %s", cliUsageDoubleHyphen), prefix)
	fmt.Println(msg)
	fmt.Println()
	// 	fmt.Printf(`Options are not position dependent and may be interspersed with file names. POSIX options in the first column can be grouped together with the exception of -h and -v which must be independent. Long options in the third column can use a Multics style single hyphen prefix as displayed or the Gnu style double hyphen prefix.

	// `)
	msg = "Options are not position dependent and may be interspersed with file names. POSIX options in the first column can be grouped together with the exception of -h and -v which must be independent. Long options in the third column can use a Multics style single hyphen prefix or the Gnu style double hyphen prefix as displayed."
	prefix = ""
	fmt.Printf("%s\n", helpWrap(msg, prefix))
}
