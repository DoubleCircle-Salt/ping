package ping

import (
	"fmt"
	"net"
	"testing"
)

func test_IcmpPing(ip net.IP) {

	packetLoss, avgRtt, err := IcmpPing(ip, nil)
	if err != nil {
		fmt.Println("icmpPing failed, err: ", err.Error())
	}

	fmt.Printf("icmpPing, packetLoss: %.2f%%, avgRtt: %.3fms\n", packetLoss, avgRtt)
}

func test_TcpPing(ip net.IP, port int) {

	packetLoss, avgRtt, err := TcpPing(ip, port, nil)
	if err != nil {
		fmt.Println("tcpPing failed, err: ", err.Error())
	}

	fmt.Printf("tcpPing, packetLoss: %.2f%%, avgRtt: %.3fms\n", packetLoss, avgRtt)
}

func test_Ping(ip net.IP, port int) {

	packetLoss, avgRtt, err := Ping(ip, port, nil)
	if err != nil {
		fmt.Println("ping failed, err: ", err.Error())
	}

	fmt.Printf("ping, packetLoss: %.2f%%, avgRtt: %.3fms\n", packetLoss, avgRtt)
}

func test_PingHost(host string, port int) {

	packetLoss, avgRtt, err := PingHost(host, port, nil)
	if err != nil {
		fmt.Println("ping host failed, err: ", err.Error())
	}

	fmt.Printf("ping host, packetLoss: %.2f%%, avgRtt: %.3fms\n", packetLoss, avgRtt)
}

func test_PingAddress(address string) {

	packetLoss, avgRtt, err := PingAddress(address, nil)
	if err != nil {
		fmt.Println("ping address failed, err: ", err.Error())
	}

	fmt.Printf("ping address, packetLoss: %.2f%%, avgRtt: %.3fms\n", packetLoss, avgRtt)
}

func Test_ping(t *testing.T) {
	test_IcmpPing(net.ParseIP("45.43.41.130"))
	test_TcpPing(net.ParseIP("45.43.41.130"), 80)
	test_Ping(net.ParseIP("45.43.41.130"), 80)
	test_PingHost("www.baidu.com", 80)
	test_PingAddress("www.baidu.com:80")
}
