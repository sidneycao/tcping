# tcping
The tcping is similarly to 'ping', but over tcp connection, And written with Golang.
  
```
Usage:  
  tcping [hostname/ip] [port] [flags]  

Examples:  
  
  1. ping www.google.com with default port 80  
    > tcping www.google.com  
  2. ping www.google.com with a custom port  
    > tcping www.google.com 443  
  

Flags:
  -c, --counters int      ping counter(if <= 0 will ping continuously) (default 4)
  -h, --help              help for tcping
  -i, --interval string   ping interval, the unit is seconds (default "1")
  -t, --timeout string    ping timeout, the unit is seconds (default "3")
```