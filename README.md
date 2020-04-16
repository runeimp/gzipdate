GzipDate
========

Simple command line app for the specific creation and handling of Gzip archives.


Features
--------

* Create gzip archives for each of a list of files, adding a file system safe ISO-601 datatime value to the filename when creating the archive
* Extract the archived file from Gzip archives restoring their original names if available. The original name is always available for archives created by GzipDate.
* Does not delete source files by default. Though they can be automatically deleted with a command line switch.
* Always uses maximum compression when creating archives
* 100% compatible with gzip and gunzip tools


Rational
--------

I like making archives with the date in the name. Especially for game saves that get auto deleted like in Roguelike type games. It's also supper annoying to me that `gunzip` doesn't use the stored filename within the archive by default. I also like my tools to be smart about if a file is already compressed or not. Gzip will give an error if you try to recompress a file that already ends in `.gz`. GzipDate just uncompresses the file instead.


