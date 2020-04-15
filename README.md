# GoPing

## Added
1. Host or ip address support.
1. ipv4 and ipv6 support.
1. TTL value can be changed.
1. Statistics of request.
1. RTT duration.
1. Number of requests flag.
1. Size of packet can be set.
1. Time delay between packets can be set.
1. Number of request can be set.
1. Sequence numbers shown.

Go >= 1.13 preferred 
## Usage

``` bash
#Build 
go build

#ping google.com with ICMPipv6 and data of 100 bytes
sudo ./GoPing -host=google.com -p=6 -s=100

#Get help menu with flags
sudo ./GoPing -help 
```
