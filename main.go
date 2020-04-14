package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shivamanipatil/Goping/ping"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	host := os.Args[1]
	totalReq := 0
	successfulReq := 0

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Starting")
Loop:
	for {
		select {
		case <-done:
			break Loop

		default:
			time.Sleep(time.Second)
			addr, duration, err := ping.Ping(host)
			totalReq++
			if err != nil {
				fmt.Print("Not succefull")
				continue
			}
			fmt.Printf("Host is %v and duration is %v\n", addr, duration)
			successfulReq++
		}
	}
	fmt.Print(totalReq, successfulReq)
}
