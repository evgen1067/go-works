package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connecting to the server")
}

func main() {
	flag.Parse()
	flags := flag.Args()
	if len(flags) != 2 {
		log.Fatal("Error. Host or Port not defined")
		return
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	address := net.JoinHostPort(flags[0], flags[1])
	telnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := telnet.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		err := telnet.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
	OUTER:
		for {
			select {
			case <-sigs:
				break OUTER
			default:
				err := telnet.Send()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
	OUTER:
		for {
			select {
			case <-sigs:
				break OUTER
			default:
				err := telnet.Receive()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	wg.Wait()
}
