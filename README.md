GzipDate
========

Simple command line app for the specific creation and handling of Gzip archives.


Features
--------

* Create gzip archives for each of a list of files, adding a file system safe [ISO-8601][ISO 8601 - Wikipedia] datatime value to the filename when creating the archive. The format is `YYYY-MM-DD_HHmmss`.
* Extract the archived file from Gzip archives restoring their original names if available. The original name is always available for archives created by GzipDate.
* Does not delete source files by default. Though they can be automatically deleted with a command line switch.
* Always uses maximum compression when creating archives
* 100% compatible with gzip and gunzip tools
* Written in Go so it's easy to cross-compile binaries for many platforms


Usage Examples
--------------

### Help

```text
$ gzipdate -h
GzipDate v1.0.0

USAGE: gzipdate [OPTIONS] [FILENAMES]

OPTIONS:
  -d | -del | -delete   Delete the source file after successful processing
  -h | -help            Display this help info
  -v | -ver | -version  Display this apps version number

Options may be interspersed with file names if so desired.
They are not position dependent.
```

### Brogue Save Game Backup

```text
runeimp$ ls -hl
total 112
-rw-r--r--  1 runeimp  staff    54K Apr 16 10:33 Saved game.broguesave

runeimp$ gzipdate -d *
Saving 21400 bytes from "Saved game.broguesave" to "Saved game.broguesave_2020-04-16_104846.gz"
    Deleting source: "Saved game.broguesave"

runeimp$ ls -hl
total 48
-rw-r--r--  1 runeimp  staff    21K Apr 16 10:48 Saved game.broguesave_2020-04-16_104846.gz

runeimp$ gzipdate *.gz
54984 B written to "Saved game.broguesave"

runeimp$ ls -hl
total 160
-rw-r--r--  1 runeimp  staff    54K Apr 16 10:33 Saved game.broguesave
-rw-r--r--  1 runeimp  staff    21K Apr 16 10:48 Saved game.broguesave_2020-04-16_104846.gz
```

Note that the size and date of "Saved game.broguesave" is the same in the original as in the version restored from the archive.

Side Note: I added an extra line before each successive command in the second example for better readability.


Rational
--------

I like making archives with the date in the name. Especially for game saves that get auto deleted like in Roguelike type games. It's also super annoying to me that `gunzip` doesn't use the stored filename within the archive by default. I also like my compression tools to be smart about if a file is already compressed or not. Gzip will give an error if you try to recompress a file that already ends in `.gz`. GzipDate just uncompresses the file instead. So just one command to use in either case.


Installation
------------

See the [installation docs](INSTALL.md)


[ISO 8601 - Wikipedia]: https://en.wikipedia.org/wiki/ISO_8601

