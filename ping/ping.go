package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func Ping(domain string) (*net.IPAddr, time.Duration, error) {
	//Resolve IP address
	addr, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		return nil, 0, err
	}

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return addr, 0, err
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
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
		return addr, 0, fmt.Errorf("Couldn't write whole message...")
	}

	receiveBuffer := make([]byte, 2000)
	n, receiveAddr, err := conn.ReadFrom(receiveBuffer)
	if err != nil {
		return addr, 0, err
	}
	span := time.Since(start)
	parsedMsg, err := icmp.ParseMessage(1, receiveBuffer[:n])
	if err != nil {
		return addr, 0, err
	}
	if parsedMsg.Type == ipv4.ICMPTypeEchoReply {
		return addr, span, nil
	} else {
		return addr, 0, fmt.Errorf("Received %+v from %v", parsedMsg, receiveAddr)
	}

}