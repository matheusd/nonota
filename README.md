# nonota

Nonota is my No Nonsense Tasking command line tool for task and time tracking.

This is a personal tool, that I wrote mostly for personal usage.

Here's an asciinema video on how it works:

[![asciicast](https://asciinema.org/a/AEXNVzvrd88wmJhMFWnSqGxzR.svg)](https://asciinema.org/a/AEXNVzvrd88wmJhMFWnSqGxzR)

Use at your own risk.

## Building and Running

Use go 1.11+ and build outside `$GOPATH` (or with `GO111MODULES=on`).

Install it with:

```
$ go get https://github.com/matheusd/nonota/cmd/nonota
```

Then just execute `nonota`. It will open or create a file named `nonota-board`.yml
on the current dir to store the data. You can manipulate this file with a 
text editor if you need to do something nonota doesn't support (like deleting
tasks).

## Exporting tasks

You can export the list of tasks for the previous month by running `nonota-csv`.
By default, it will export the tasks for the previous billable month.

