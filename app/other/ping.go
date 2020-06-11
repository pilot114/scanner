package other

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ProtocolICMP = 1
)

// Default to listen on all IPv4 interfaces
var ListenAddr = "0.0.0.0"
var Timeout = 1

func Receiver(addr string) *icmp.PacketConn {
	// слушаем входящие пакеты
	c, err := icmp.ListenPacket("ip4:icmp", addr)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	defer c.Close()
	return c
}

func Ping(receiver *icmp.PacketConn, addr string) (*net.IPAddr, time.Duration, error) {

	// резолвим DNS. TODO: можно не резолвить?
	dst, err := net.ResolveIPAddr("ip4", addr)

	// Создаем ICMP сообщение. Разные типы подразумевают разные Body
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Intn(65000),
			Seq:  1,
			Data: []byte(""),
		},
	}
	data, err := m.Marshal(nil)
	if err != nil {
		return dst, 0, err
	}

	// Отправка
	start := time.Now()
	n, err := receiver.WriteTo(data, dst)
	if err != nil {
		return dst, 0, err
	} else if n != len(data) {
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(data))
	}

	// Ожидаем ответ
	reply := make([]byte, 1500)
	err = receiver.SetReadDeadline(time.Now().Add(time.Duration(Timeout) * time.Second))
	if err != nil {
		return dst, 0, err
	}
	n, peer, err := receiver.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	// Pack it up boys, we're done here
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return dst, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return dst, duration, nil
	default:
		return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
}

func checkSum(msg []byte) uint16 {
	sum := 0

	// assume even for now
	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	ip := ""
	a := "24"

	//connection := Receiver("0.0.0.0")
	//p := func(addr string){
	//	dst, dur, err := Ping(connection, addr)
	//	if err != nil {
	//		fmt.Printf("%s (%s)\n", dst, err)
	//		return
	//	}
	//	fmt.Printf("%s %s\n", dst, dur)
	//}

	//p("192.30.253.113")

	addr, _ := net.ResolveIPAddr("ip", "192.30.253.113")
	conn, _ := net.DialIP("ip4:icmp", addr, addr)

	var msg [512]byte
	msg[0] = 8  // echo
	msg[1] = 0  // code 0
	msg[2] = 0  // checksum, fix later
	msg[3] = 0  // checksum, fix later
	msg[4] = 0  // identifier[0]
	msg[5] = 13 //identifier[1]
	msg[6] = 0  // sequence[0]
	msg[7] = 37 // sequence[1]
	len := 8

	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)

	_, err := conn.Write(msg[0:len])
	checkError(err)

	_, err = conn.Read(msg[0:])
	checkError(err)

	fmt.Println("Got response")
	if msg[5] == 13 {
		fmt.Println("identifier matches")
	}
	if msg[7] == 37 {
		fmt.Println("Sequence matches")
	}

	go func() {
		for b := 0; b <= 255; b++ {
			for c := 0; c <= 255; c++ {
				for d := 0; d <= 255; d++ {
					ip = fmt.Sprintf("%s.%s.%s.%s", a, strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d))
					//p(ip)
				}
			}
		}
	}()

	done := make(chan int)
	<-done
}
