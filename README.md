# waybackrobots

Collect old robots.txt files from waybackmachine and download the Disallow paths

Install:
```shell
$ go install github.com/dvvedz/waybackrobots@latest 
```

Example:

```shell
$ waybackrobots -domain google.com -from 2020

$ waybackrobots -domain google.com

$ cat domains.txt | waybackrobots -from 2020
```

Help:
```plain
Usage of waybackrobots:
  -domain string
        which domain to find old robots for
  -fd int
        choose date from when to get robots from format: 2015 (default 2015)
```
