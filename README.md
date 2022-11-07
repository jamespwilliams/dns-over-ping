# dns-over-icmp

```
$ ping localhost -4 -p "$(printf "cloudflare.com?" | xxd -p)" -c1
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

```
  68 10 85 e5 68 10 84 e5
= 68 10 85 e5, 68 10 84 e5 
= 0x68.0x10.0x85.0xe5, 0x68.0x10.0x84.0xe5 
= 0x68.0x10.0x85.0xe5, 0x68.0x10.0x84.0xe5 
= 104.16.133.229, 104.16.132.229
```
