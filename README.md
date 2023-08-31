# tcping
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsidneycao%2Ftcping.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsidneycao%2Ftcping?ref=badge_shield)

The tcping is similarly to 'ping', but over tcp connection, And written with Golang.  
Download the release based on your platform, and move it to your environment PATH as necessary.  
If there is a desire to use it on other platforms, you can compile it yourself from the source code.  
</br>  
</br>  

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

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsidneycao%2Ftcping.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsidneycao%2Ftcping?ref=badge_large)