# waybackrobots

Collect old robots.txt files from waybackmachine and download the Disallow paths

Install:
```shell
$ go install github.com/ogow/waybackrobots@latest 
```

Usage:

Just give the flag `-domain` a valid domain and it will start downloading all the archived responses from 2015

```shell
$ waybackrobots -domain google.com 
```

Get all response from a custom year:
```shell
$ waybackrobots -domain google.com -fd 2020
```

Sometimes the wayback api can return alot of results which will take along time to download, to avoid this the `-strat` flag can be used.
The `-strat` flag takes one of these values `day, month, digest`, digest is default.
Usually `digest` will be the go to value, but if we look at a domain like google.com that has been archived alot the `digest` filter will still return alot of results.
If this is the case we can try to use the `day` filter which gets one snapshot each day.

filters in use explanation: https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#collapsing

comparing digest with day and month
```shell
$ go run . -domain google.com -strat digest
[i] found 38261 old robots.txt files

$ go run . -domain google.com -strat day 
[i] found 473 old robots.txt files

$ go run . -domain google.com -strat month
[i] found 122 old robots.txt files
```


Help:
```plain
Usage of waybackrobots:
  -domain string
        which domain to find old robots for
  -fd int
        choose date from when to get robots from format: 2015 (default 2015)
  -strat string
        interval to get robots for, possible values: digest, day, month (default "digest") 
```
