package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/DoubleCircle-Salt/ping"
)

func main() {

	var (
		host     string
		port     int
		count    int
		interval int
		timeout  int
		dns      string
	)
	
	flag.StringVar(&host, "h", "", "dest host")
	flag.IntVar(&port, "p", 0, "dest port")
	flag.IntVar(&count, "c", 10, "ping count")
	flag.IntVar(&interval, "i", 1000, "ping interval(ms)")
	flag.IntVar(&timeout, "t", 1, "ping timeout(s)")
	flag.StringVar(&dns, "d", "", "dns server")

	flag.Parse()

	if host == "" {
		fmt.Printf("host not found\n")
		return
	}

	resolver := net.DefaultResolver
	if dns != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second,
				}
				return d.DialContext(ctx, "udp", fmt.Sprintf("%s:53", dns))
			},
		}
	}

	config := &ping.Config{
		Count:    count,
		Interval: time.Millisecond,
		Timeout:  time.Duration(timeout) * time.Second,
		Resolver: resolver,
	}

	for {
		packetLoss, avgRtt, err := ping.PingHost(host, port, config)
		if err != nil {
			fmt.Printf("ping address %s:%d failed, err: %s\n", host, port, err.Error())
			time.Sleep(time.Duration(interval) * time.Millisecond)
			continue
		}

		fmt.Printf("ping address %s:%d success, packageLoss: %.3f %%, avgRtt: %.3f ms\n", host, port, packetLoss, avgRtt)

		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}