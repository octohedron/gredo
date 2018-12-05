# GREDO
### Go redis set exporter / importer

Only `SRANDMEMBER` is supported at the moment

TLDR: `go build && ./gredo {dump|load} {SET_IDENTIFIER} {AMOUNT}`

### Examples

To dump 30 random members of `my_set` to a file called `my_set.txt`
```
$ go build && ./gredo dump my_set 30
```

To load **up to** 30 items in a file called `my_set.txt` in the same directory
```
$ go build && ./gredo load my_set 30
```

Notes
+ When dumping it will overwrite the previous file, i.e. `my_set.txt`
+ It only connects to redis in `localhost:6379`
+ Loads/Dumps more than 10K items in less than a second, probably a lot more