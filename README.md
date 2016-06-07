# redirectresolver

A quick and dirty implementation for resolving a list of urls to their destination.

## Input

A text file containing one URL per line.

```
http://example.org
http://en.wikipedia.org
```

## Output

```
{"from":"http://example.org","to":"http://example.org","error":""}
{"from":"http://en.wikipedia.org","to":"https://en.wikipedia.org/wiki/Main_Page","error":""}
```

## Run it

```
$ go build .
$ ./redirectresolver --help
```
