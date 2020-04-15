package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shivamanipatil/GoPing/ping"
)

func main() {

	//Command line flags
	hostPtr := flag.String("host", "localhost", "Host name or IP address")
	countPtr := flag.Int("c", -1, "Stop after sending this number of requests. Default -1 is for infinite.")
	intervalPtr := flag.Float64("i", 1, "Interval seconds between two requests.")
	sizePtr := flag.Int("s", 56, "Size of packet.")
	protocolPtr := flag.Int("p", 4, "Protocol number 4 for ipv4 and 6 for ipv6.")
	flag.Parse()

	if *protocolPtr != 4 && *protocolPtr != 6 {
		fmt.Printf("Enter 4 for ipv4 and 6 for ipv6.\n")
		os.Exit(1)
	}

	//Go routine to listen for SIGINT signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	//Variable initialization
	totalReq := 0
	receivedReq := 0
	var durations []time.Duration
	//Generic type to hold ping.Ping4 and ping.Ping6
	var Ping func(domain string, seqNumber int, size int) (*net.IPAddr, time.Duration, error)
	if *protocolPtr == 4 {
		Ping = ping.Ping4
	} else {
		Ping = ping.Ping6
	}
	fmt.Printf("PING %v sending %v bytes of data\n", *hostPtr, *sizePtr)

	//ping loop
Loop:
	for {
		select {
		case <-done:
			break Loop
		default:
			if *countPtr == totalReq {
				break Loop
			}
			time.Sleep(time.Duration(*intervalPtr) * time.Second)
			_, duration, err := Ping(*hostPtr, receivedReq, *sizePtr)
			totalReq++
			if err != nil {
				fmt.Print(err)
				break Loop
			}
			fmt.Printf("time=%v\n", duration)
			durations = append(durations, duration)
			receivedReq++
		}
	}
	statistics(receivedReq, totalReq, durations)
}

func statistics(receivedReq, totalReq int, durations []time.Duration) {
	fmt.Printf("\n----- Statistics -----\n")
	fmt.Printf("Sent: %v Received: %v Loss: %v percent\n", totalReq, receivedReq, (float32(totalReq-receivedReq)*100)/float32(receivedReq))
	if len(durations) > 0 {
		min, max, avg := durations[0], durations[0], time.Duration(0)
		for _, duration := range durations {
			if duration < min {
				min = duration
			}
			if duration > max {
				max = duration
			}
			avg += duration
		}
		fmt.Printf("rtt --- min / max / avg : %v / %v / %v \n", min, max, avg/time.Duration(len(durations)))
	}
}
