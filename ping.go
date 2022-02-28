package ping

import (
	"context"
	"fmt"
	"net"
	"time"

	gping "github.com/go-ping/ping"
)

type Config struct {
	Count    int
	Interval time.Duration
	Timeout  time.Duration
	Ports    []int
	Resolver *net.Resolver
}

var defaultConfig = &Config{
	Count:    10,
	Interval: time.Millisecond,
	Timeout:  time.Second,
	Ports:    []int{443, 80},
	Resolver: net.DefaultResolver,
}

func icmpPing(ip net.IP, config *Config) (packetLoss, avgRtt float64, err error) {
	pinger, err := gping.NewPinger(ip.String())
	if err != nil {
		return -1, -1, err
	}
	pinger.OnFinish = func(stats *gping.Statistics) {
		packetLoss = stats.PacketLoss
		avgRtt = float64(stats.AvgRtt.Nanoseconds())/1000000
	}
	pinger.Count = config.Count
	pinger.Interval = config.Interval
	pinger.Timeout = config.Timeout
	err = pinger.Run()
	if err != nil {
		return -1, -1, err
	}
	return
}

func IcmpPing(ip net.IP, config *Config) (packetLoss, avgRtt float64, err error) {
	if config == nil {
		config = defaultConfig
	}
	return icmpPing(ip, config)
}

func tcpPing(ip net.IP, port int, config *Config) (packetLoss, avgRtt float64, err error) {
	var avgRtt0 float64
	var count   int

	dialer := &net.Dialer{
		Timeout: config.Timeout,
	}

	packetBase := float64(100)/float64(config.Count)

	for i := 0; i < config.Count; i++ {
		startTime := time.Now()
		conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", ip.String(), port))
		if err != nil {
			packetLoss += packetBase
			continue
		}
		avgRtt0 += float64(time.Now().Sub(startTime).Nanoseconds())/1000000
		count++
		conn.Close()
	}
	if count != 0 {
		avgRtt = avgRtt0 / float64(count)
	}
	return
}

func TcpPing(ip net.IP, port int, config *Config) (packetLoss, avgRtt float64, err error) {
	if config == nil {
		config = defaultConfig
	}
	return tcpPing(ip, port, config)
}

func ping(ip net.IP, port int, config *Config) (packetLoss, avgRtt float64, err error) {
	packetLoss, avgRtt, err = icmpPing(ip, config)
	if err == nil && packetLoss != 100 {
		return packetLoss, avgRtt, nil
	}
	return tcpPing(ip, port, config)
}

func Ping(ip net.IP, port int, config *Config) (packetLoss, avgRtt float64, err error) {
	if config == nil {
		config = defaultConfig
	}
	return ping(ip, port, config)
}

func pingHost(host string, port int, config *Config) (packetLoss, avgRtt float64, err error) {
	var ips []string
	ip := net.ParseIP(host)
	if ip != nil {
		ips = append(ips, ip.String())
	} else {
		ips, err = config.Resolver.LookupHost(context.Background(), host)
		if err != nil {
			return -1, -1, err
		}
	}

	count := 0
	for i := 0; i < len(ips); i++ {
		ip := net.ParseIP(ips[i])
		if ip != nil {
			packetLoss0, avgRtt0, err := ping(ip, port, config)
			if err != nil || packetLoss0 == 100 {
				continue
			}
			count++
			packetLoss += packetLoss0
			avgRtt += avgRtt0
		}
	}

	if count != 0 {
		return packetLoss/float64(count), avgRtt/float64(count), nil
	} else {
		return -1, -1, fmt.Errorf("icmpPing/tcpPing %s:%d failed", host, port)
	}
}

func PingHost(host string, port int, config *Config) (packetLoss, avgRtt float64, err error) {
	if config == nil {
		config = defaultConfig
	}
	return pingHost(host, port, config)
}
