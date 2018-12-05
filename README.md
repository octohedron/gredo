# GREDO
### Go redis exporter
Only `SRANDMEMBER` is supported at the moment

TLDR: `go build && ./gredo {dump|load} {SET_IDENTIFIER} {AMOUNT}`

Example
```
$ go build && ./gredo dump my_set 30
```

Will dump 30 random members of `my_set` to a file called `my_set.txt`

Notes
+ Be careful, it will overwrite the previous file, i.e. `my_set.txt`
+ For now it only connects to redis in `localhost:6379`
