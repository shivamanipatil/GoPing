package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shivamanipatil/GoPing/ping"
)

func main() {

	//Command line args and flags
	hostPtr := flag.String("host", "localhost", "Host name or IP address")
	countPtr := flag.Int("c", -1, "Stop after sending this number of requests.")
	flag.Parse()

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

	//Ping loop
	totalReq := 0
	successfulReq := 0
Loop:
	for {
		select {
		case <-done:
			break Loop
		default:
			if *countPtr == totalReq {
				break Loop
			}
			fmt.Println(*countPtr, totalReq)
			time.Sleep(time.Second)
			addr, duration, err := ping.Ping(*hostPtr, successfulReq)
			totalReq++
			if err != nil {
				fmt.Print(err)
				continue
			}
			fmt.Printf("Host is %v and duration is %v\n", addr, duration)
			successfulReq++
		}
	}
	fmt.Print(totalReq, successfulReq)
}
