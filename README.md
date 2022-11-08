# dns-over-ping(8)

You've heard of DNS-over-HTTP, DNS-over-TLS, DNS-over-GRPC... Now get ready for
DNS-over-ping(8)!

Resolve names straight from the standard inetutils/iptools `ping` tool:

```
$ ping localhost -4 -p "$(printf "cloudflare.com?" | xxd -p -c0)" -c1
PATTERN: 0x636c6f7564666c6172652e636f6d3f
PING  (127.0.0.1) 56(84) bytes of data.
72 bytes from localhost (127.0.0.1): icmp_seq=1 ttl=64 time=1667863702756 ms
wrong data byte #16 should be 0x6c but was 0x68
#16	68 10 85 e5 68 10 84 e5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
#48	0 0 0 0 0 0 0 0

---  ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 1667863702755.909/1667863702755.909/1667863702755.909/0.000 ms
```

Simply read your answer off from the `wrong data` hexdump:

```
   68 10 85 e5          68 10 84 e5
=> 0x68.0x10.0x85.0xe5  0x68.0x10.0x84.0xe5 
=  104.16.133.229       104.16.132.229
```

### Limitations

* Only `A` lookups supported, sorry.
    * Including both IPv4 and IPv6 addresses would necessitate a more
      complicated output format to avoid ambiguity.
    * The choice is then between IPv4 and IPv6 - IPv4 wins because we can
      display more IPs (`ping` will only show data 56 bytes in its errors).

* Names can be at most 15 bytes long.
    * This is because `ping` only lets you specify 16 byte data patterns -
      everything beyond 16 bytes is ignored - and a byte more is required for
      the delimiter (question mark)

* At most 12 IPs can be returned
    * `ping` will always display 56 bytes of hexdumped wrong data, regardless
      of how much is in the response packet

None of these are inherent limitations of ICMP, rather they are limitations of
the `ping` tool and its output. DNS-over-ICMP could actually be made to work
pretty well (but would be less fun).

### Running

Optionally, prevent your machine sending its own ICMP responses to incoming
ICMP echo requests:

```
# echo "1" > /proc/sys/net/ipv4/icmp_echo_ignore_all
```

Then:

```
$ go build ./cmd/pingdns
$ sudo setcap cap_net_raw+ep pingdns
$ ./pingdns
```

Then, in another shell, for example:

```
$ ping localhost -4 -p "$(printf "jameswillia.ms?" | xxd -p -c0)" -c1
```
