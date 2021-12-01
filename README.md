# BENCH DNS

Simple script for DNS upload / requests processing tests. Written in Go.

Allow to observe static on the DNS servers, such as - errors, speed of response, cache, requests blocking, cache and etc.

## Features

* Using local domain list file with `-file` argument
* Download from URL and then using downloaded file with `-url` argument
* Using custom IP for DNS server - `-ip` argument
* Using iterations with `-iterate` argument

## Examples 

Using local file:
```
go run bench-dns.go -file "hosts.txt" -ip "1.1.1.1" -iterate 50
```

Using remote file (will saved locally as `domains.txt`):
```
go run bench-dns.go -ip "1.1.1.1" -iterate 100 -url "https://raw.githubusercontent.com/m0zgen/dns-hole/master/whitelist.txt"
```
## How to build

Clone repository:

```
git clone https://github.com/m0zgen/bench-dns.git
```

`cd` to repo caytalog and build to linux or own disoro/system:

```
env GOOS=linux GOARCH=amd64 go build bench-dns.go
``` 

Or just:

```
go build bench-dns.go
```

Then run:

```

./bench-dns                                                                                                                                              765 ms  master 
You need use one argement - url or file
  -file string
        Local domain list
  -ip string
        DNS server IP address
  -iterate int
        Repeat counts (default 1)
  -url string
        URL to raw domain list
```