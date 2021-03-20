package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pilot114/scanner/proto"
)

// ResponseHTTPInfo : ответ
type ResponseHTTPInfo struct {
	headers map[string]string
	time    time.Duration
	ip      string
	error   string
}

// ResponseInfo : ответ
type ResponseInfo struct {
	ip   string
	resv int
	sent int
	avg  time.Duration
}

func getHeaders(url string) ResponseHTTPInfo {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	timeout := time.Duration(3 * time.Second)
	client := &http.Client{Transport: transport, Timeout: timeout}

	start := time.Now()
	response, err := client.Get(fmt.Sprintf("http://%s", url))
	duration := time.Since(start)

	info := ResponseHTTPInfo{make(map[string]string), duration, url, ""}

	if err != nil {
		info.error = fmt.Sprintf("Error download: %s", err)
		return info
	}

	if response.StatusCode != http.StatusOK {
		info.error = fmt.Sprintf("Error HTTP Status: %s", response.Status)
		return info
	}

	for k, v := range response.Header {
		info.headers[strings.ToLower(k)] = string(v[0])
	}
	return info
}

func icmp(url string) ResponseInfo {

	sent := 3 // сколько отправили
	resv := 0 // сколько получили
	avg := time.Duration(0)

	for i := 1; i <= sent; i++ {
		duration, err := proto.Ping(url)
		avg = time.Duration(int64(avg+duration) / int64(i))
		if err == nil {
			resv = resv + 1
		}
		// быстро или надежно? =)
		// time.Sleep(time.Millisecond * 100)
	}

	return ResponseInfo{
		ip:   url,
		resv: resv,
		sent: sent,
		avg:  avg,
	}
}

func worker(wid int, ips <-chan string, responses chan<- ResponseInfo) {
	for ip := range ips {
		// fmt.Printf("worker %d get ip %s\n", wid, ip)
		responses <- icmp(ip)
	}
}

func main() {
	a := os.Args[1]
	b := os.Args[2]
	workerLimit, _ := strconv.Atoi(os.Args[3])

	// каналы: источник адресов и получатель заголовков
	ips := make(chan string, 66000)
	resInfo := make(chan ResponseInfo, 3)

	// стартуем воркеров
	for wid := 1; wid <= workerLimit; wid++ {
		go worker(wid, ips, resInfo)
	}

	start := time.Now()
	ip := ""

	// https://ant.isi.edu/address/
	go func() {
		for c := 0; c <= 255; c++ {
			for d := 0; d <= 255; d++ {
				ip = fmt.Sprintf("%s.%s.%s.%s", a, b, strconv.Itoa(c), strconv.Itoa(d))
				ips <- ip
			}
		}
		close(ips)
	}()

	total := 256 * 256
	successCount := 0

	for total > 0 {
		info := <-resInfo

		if info.resv > 0 {
			// json, _ := json.Marshal(info)
			fmt.Printf("%s %d %d %s\n", info.ip, info.resv, info.sent, info.avg)
			successCount = successCount + 1
		} else {
			// json, _ := json.Marshal(info)
			fmt.Fprintf(os.Stderr, "%s %d %d\n", info.ip, info.resv, info.sent)
		}
		total = total - 1
	}

	duration := time.Since(start)
	fmt.Printf("Total time: %s, found: %d\n", duration, successCount)
}
