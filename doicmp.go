package doicmp

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	icmpv4ChecksumLength = 16

	pingMinimumDataLength = 64
	pingLeadingDataLength = 16

	defaultSnaplen = 1600
	defaultNetwork = "ip4:icmp"
	defaultAddress = "0.0.0.0"
)

type Server struct {
	snaplen          int
	network, address string
}

func NewServer() Server {
	return Server{
		snaplen: defaultSnaplen, network: defaultNetwork, address: defaultAddress,
	}
}

func (s Server) WithSnaplen(snaplen int) {
	s.snaplen = snaplen
}

func (s Server) WithNetwork(network string) {
	s.network = network
}

func (s Server) Address(address string) {
	s.address = address
}

func (s Server) Serve() error {
	conn, err := icmp.ListenPacket(s.network, s.address)
	if err != nil {
		return fmt.Errorf("doicmp: listening failed: %w", err)
	}
	defer conn.Close()

	for {
		request := make([]byte, s.snaplen)

		n, addr, err := conn.ReadFrom(request)
		if err != nil {
			return fmt.Errorf("doicmp: reading packets failed: %w", err)
		}

		response, err := handleBytes(request[:n])
		if err != nil {
			return err
		}

		if response == nil {
			continue
		}

		if _, err := conn.WriteTo(response, addr); err != nil {
			fmt.Printf("doicmp: warn: failed to write response: %v\n", err)
		}
	}
}

func handleBytes(requestBytes []byte) (response []byte, err error) {
	parsed, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), requestBytes)
	if err != nil {
		return nil, fmt.Errorf("doicmp: parsing message failed: %w", err)
	}

	switch parsed.Type {
	case ipv4.ICMPTypeEcho:
	default:
		return nil, nil
	}

	icmpEchoRequest, ok := parsed.Body.(*icmp.Echo)
	if !ok {
		return nil, fmt.Errorf("packet wasn't icmp echo?")
	}

	icmpMessageResponse, err := handleICMPEcho(icmpEchoRequest)
	if err != nil {
		return nil, err
	}

	return icmpMessageResponse.Marshal(nil)
}

func handleICMPEcho(request *icmp.Echo) (response *icmp.Message, err error) {
	name, err := extractNameFromPayload(request.Data[icmpv4ChecksumLength:])
	if err != nil {
		return nil, fmt.Errorf("failed to extract name from icmp payload: %w", err)
	}

	ipv4s, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve name: %w", err)
	}

	return &icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Code: 0,
		Body: &icmp.Echo{
			ID:   request.ID,
			Seq:  request.Seq,
			Data: prepareResponseData(ipv4s),
		},
	}, nil
}

func prepareResponseData(ipv4s []net.IP) []byte {
	// When ping presents the "wrong data" message, it chops off `pingLeadingDataLength`
	// of leading data. Pad that out with zeroes, so our IPs get shown:
	responseData := make([]byte, pingLeadingDataLength)
	responseData = append(responseData, flatten(ipv4sToByteSlices(ipv4s))...)

	padding := make([]byte, pingMinimumDataLength-len(responseData))
	responseData = append(responseData, padding...)

	return responseData
}

// this is kind of shit, there's probably a better way
func extractNameFromPayload(payload []byte) (string, error) {
	// TODO check if -1
	start := findIndex(payload, '?')
	if start == -1 {
		return "", nil
	}

	end := findIndex(payload[start+1:], '?')
	if end == -1 {
		return "", nil
	}

	return string(payload[start+1 : start+1+end]), nil
}

func findIndex(bytes []byte, delim byte) int {
	for i := range bytes {
		if bytes[i] == delim {
			return i
		}
	}
	return -1
}

func flatten(slices [][]byte) []byte {
	var res []byte
	for _, slice := range slices {
		res = append(res, slice...)
	}
	return res
}

func ipv4sToByteSlices(ips []net.IP) [][]byte {
	var slices [][]byte
	for _, ip := range ips {
		if ip == nil || len(ip) != 16 {
			continue
		}

		slices = append(slices, ip[12:])
	}
	return slices
}
