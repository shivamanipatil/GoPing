package ping

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

//Ping4 return checking code, duration and error
func Ping4(domain string, seqNumber int, size int, ttl int, wait int) (int, time.Duration, error) {
	//LookupIP gives Ip's slice we extract ipv6 address and convert to net.IPAddr from IP
	ipSlice, err := net.LookupIP(domain)
	if err != nil {
		return 0, 0, err
	}
	ip, err := findIpv4(ipSlice)
	if err != nil {
		return 0, 0, err
	}
	addr, _ := net.ResolveIPAddr("ip", ip.String())

	//Connection for listening packets
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	//This connection is needed for setting TTL flag and also accessing ip headers like TTL
	connNew := ipv4.NewPacketConn(conn)
	if err := connNew.SetControlMessage(ipv4.FlagTTL, true); err != nil {
		return 0, 0, err
	}

	//Set TTL
	err = connNew.SetTTL(ttl)
	if err != nil {
		return 0, 0, err
	}

	//Making ICMP message body and then marshalling
	data := make([]byte, size)
	data = bytes.Repeat([]byte("A"), size)
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
		return 0, 0, err
	}

	//Timer started just before writing message to connection
	start := time.Now()
	n, err := connNew.WriteTo(b, nil, addr)
	if err != nil {
		return 0, 0, err
	}
	if n != len(b) {
		return 0, 0, fmt.Errorf("Couldn't write whole message...\n")
	}

	//Setting read deadline
	err = connNew.SetReadDeadline(time.Now().Add(time.Duration(wait) * time.Second))
	if err != nil {
		return 1, 0, err
	}

	//Reading from connection
	receiveIPBuffer := make([]byte, 2000)
	n, controlMessage, receiveAddr, err := connNew.ReadFrom(receiveIPBuffer)
	if err != nil {
		return 1, 0, err
	}

	duration := time.Since(start)

	//Parsing the message
	parsedMsg, err := icmp.ParseMessage(1, receiveIPBuffer[:n])
	if err != nil {
		return 1, 0, err
	}
	switch parsedMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("%v bytes from %v: icmp_seq=%v ttl=%v time=%v\n",
			n, receiveAddr, parsedMsg.Body.(*icmp.Echo).Seq, controlMessage.TTL, duration)
		return 1, duration, nil
	case ipv4.ICMPTypeTimeExceeded:
		return 1, duration, fmt.Errorf("Time limit exceeded.")
	default:
		return 1, duration, fmt.Errorf("Other ICMP Error.")
	}

}

//Ping6 return checking code, duration and error
func Ping6(domain string, seqNumber int, size int, ttl, wait int) (int, time.Duration, error) {
	//LookupIP gives Ip's slice we extract ipv6 address and convert to net.IPAddr from IP
	ipSlice, err := net.LookupIP(domain)
	if err != nil {
		return 0, 0, err
	}
	ip, err := findIpv6(ipSlice)
	if err != nil {
		return 0, 0, err
	}
	addr, _ := net.ResolveIPAddr("ip", ip.String())

	//Connection for listening packets
	conn, err := net.ListenPacket("ip6:ipv6-icmp", "::")
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	//This connection is needed for setting TTL flag and also accessing ip headers like TTL
	connNew := ipv6.NewPacketConn(conn)
	if err := connNew.SetControlMessage(ipv6.FlagHopLimit, true); err != nil {
		return 0, 0, err
	}

	//Set Hop limit
	err = connNew.SetHopLimit(ttl)
	if err != nil {
		return 0, 0, err
	}

	//Making ICMP message body then marshalling
	data := make([]byte, size)
	data = bytes.Repeat([]byte("A"), size)
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
		return 0, 0, err
	}

	//Timer started just before writing message to connection
	start := time.Now()
	i, err := connNew.WriteTo(b, nil, addr)
	if err != nil {
		return 0, 0, err
	}
	if i != len(b) {
		return 0, 0, fmt.Errorf("Couldn't write whole message...\n")
	}

	//Setiting read deadline
	err = connNew.SetReadDeadline(time.Now().Add(time.Duration(wait) * time.Second))
	if err != nil {
		return 1, 0, err
	}

	//Reading from connection
	receiveIPBuffer := make([]byte, 2000)
	n, controlMessage, receiveAddr, err := connNew.ReadFrom(receiveIPBuffer)
	if err != nil {
		return 1, 0, err
	}

	duration := time.Since(start)

	//Parsing the message
	parsedMsg, err := icmp.ParseMessage(58, receiveIPBuffer[:n])
	if err != nil {
		return 1, 0, err
	}

	switch parsedMsg.Type {
	case ipv6.ICMPTypeEchoReply:
		fmt.Printf("%v bytes from %v: icmp_seq=%v ttl=%v time=%v\n",
			n, receiveAddr, parsedMsg.Body.(*icmp.Echo).Seq, controlMessage.HopLimit, duration)
		return 1, duration, nil
	case ipv6.ICMPTypeTimeExceeded:
		fmt.Println("Time limit exceeded.")
		return 1, duration, nil
	default:
		fmt.Println("Unknown ICMP message.")
		return 1, duration, nil
	}
}

// Return first ipv4 address in slice of net.IP else nil
func findIpv4(addresses []net.IP) (net.IP, error) {
	for _, address := range addresses {
		if address.To4() != nil {
			return address, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find ipv4 address of host.")
}

// Return first ipv6 address in slice of net.IP else nil
func findIpv6(addresses []net.IP) (net.IP, error) {
	for _, address := range addresses {
		if address.To4() == nil {
			return address, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find ipv6 address of host.")
}
