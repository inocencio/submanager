# SubManager

`Submanager` is a tiny program to sync subtitles in SRT files.

## Instalation

Considering you already have [Go](https://golang.org) installed, just do it:

```bash
$ go get -u github.com/inocencio/submanager
```

## Usage

First off, make sure `submanger` is on your path system to be called anywhere. Both following ways are valids to make it works:

- Passing arguments directly using `file` and `time` flags.
- Or just calling `submanager` executable without arguments (some options will pop up in CLI prompt).

For first option, follow these instructions down below: 

- Usage: ./submanager -file:<subtitle_file.srt> -time=<time_in_ms>
- E.g:   ./submanager -file:"lawrence of arabia.srt" -time=-2000

`time` flag is an integer input by milleseconds (ms);
`file` flag is a string input of srt file name;

For second option, the program will find the SRT file corresponding to the video. A menu with an input time will shows up. Just pick one or chose custom option.
