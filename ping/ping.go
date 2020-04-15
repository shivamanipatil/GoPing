package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

//Ping4 Used with ipv4
func Ping4(domain string, seqNumber int, size int) (*net.IPAddr, time.Duration, error) {
	//LookupIP gives Ip's slice we extract ipv6 address and convert to net.IPAddr from IP
	ipSlice, err := net.LookupIP(domain)
	if err != nil {
		return nil, 0, err
	}
	ip, err := findIpv4(ipSlice)
	if err != nil {
		return nil, 0, err
	}
	addr, _ := net.ResolveIPAddr("ip", ip.String())

	//Connection for listening packets
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return addr, 0, err
	}
	defer conn.Close()

	//Making ICMP message body and then marshalling
	data := make([]byte, size)
	copy(data[:size], "abcd")
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seqNumber,
			Data: data,
		}}
	b, err := msg.Marshal(nil)
	if err != nil {
		return addr, 0, err
	}

	//Timer started just before writing message to connection
	start := time.Now()
	i, err := conn.WriteTo(b, addr)
	if err != nil {
		return addr, 0, err
	}
	if i != len(b) {
		return addr, 0, fmt.Errorf("Couldn't write whole message...\n")
	}

	//Make receive buffer and read to it
	receiveBuffer := make([]byte, 2000)
	n, receiveAddr, err := conn.ReadFrom(receiveBuffer)
	if err != nil {
		return addr, 0, err
	}
	span := time.Since(start)

	//Parsing the message
	parsedMsg, err := icmp.ParseMessage(1, receiveBuffer[:n])
	if err != nil {
		return addr, 0, err
	}
	if parsedMsg.Type == ipv4.ICMPTypeEchoReply {
		fmt.Printf("%v bytes from %v: icmp_seq=%v ", n, receiveAddr, parsedMsg.Body.(*icmp.Echo).Seq)
		return addr, span, nil
	} else {
		return addr, 0, fmt.Errorf("Received %+v from %v", parsedMsg, receiveAddr)
	}

}

//Ping6 Used with ipv6
func Ping6(domain string, seqNumber int, size int) (*net.IPAddr, time.Duration, error) {
	//LookupIP gives Ip's slice we extract ipv6 address and convert to net.IPAddr from IP
	ipSlice, err := net.LookupIP(domain)
	if err != nil {
		return nil, 0, err
	}
	ip, err := findIpv6(ipSlice)
	if err != nil {
		return nil, 0, err
	}
	addr, _ := net.ResolveIPAddr("ip", ip.String())

	//Connection for listening packets
	conn, err := net.ListenPacket("ip6:ipv6-icmp", "::")
	if err != nil {
		return addr, 0, err
	}
	defer conn.Close()

	//Making ICMP message body then marshalling
	data := make([]byte, size)
	copy(data[:size], "abcd")
	msg := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seqNumber,
			Data: data,
		}}
	b, err := msg.Marshal(nil)
	if err != nil {
		return addr, 0, err
	}

	start := time.Now()
	i, err := conn.WriteTo(b, addr)
	if err != nil {
		return addr, 0, err
	}
	if i != len(b) {
		return addr, 0, fmt.Errorf("Couldn't write whole message...\n")
	}

	//Make receive buffer and read to it
	receiveBuffer := make([]byte, 2000)
	n, receiveAddr, err := conn.ReadFrom(receiveBuffer)
	if err != nil {
		return addr, 0, err
	}
	span := time.Since(start)

	//Parsing the message
	parsedMsg, err := icmp.ParseMessage(58, receiveBuffer[:n])
	if err != nil {
		return addr, 0, err
	}
	if parsedMsg.Type == ipv6.ICMPTypeEchoReply {
		fmt.Printf("%v bytes from %v: icmp_seq=%v ", n, receiveAddr, parsedMsg.Body.(*icmp.Echo).Seq)
		return addr, span, nil
	} else {
		return addr, 0, fmt.Errorf("Received %+v from %v", parsedMsg, receiveAddr)
	}

}

func findIpv4(addresses []net.IP) (net.IP, error) {
	for _, address := range addresses {
		if address.To4() != nil {
			return address, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find ipv4 address of host.")
}

func findIpv6(addresses []net.IP) (net.IP, error) {
	for _, address := range addresses {
		if address.To4() == nil {
			return address, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find ipv6 address of host.")
}
