# GREDO
### Go redis exporter
Only `SRANDMEMBER` is supported at the moment

Example
```
$ go build && ./gredo my_set 30
```

Will export 30 random members of my_set to a file called my_set.txt

NOTE: For now it only connects to redis in `localhost:6379`
