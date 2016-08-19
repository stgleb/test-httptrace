package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
)

func main() {
	ctx := httptrace.WithClientTrace(context.Background(), &httptrace.ClientTrace{
		PutIdleConn: func(err error) {
			if err == nil {
				log.Println("Put connection to idle pool\n")
			} else {
				log.Println("Put connection to idle pool: %s\n",
					err.Error())
			}
		},
		GetConn: func(hostPort string) {
			log.Printf("Take connection on %s\n", hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Printf("Got reused connection: %t\n",
				connInfo.Reused)
		},
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			log.Printf("Start DNS resolving on %s\n",
				dnsInfo.Host)
		},
		DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
			log.Printf("Done DNS resolving %s\n",
				dnsDoneInfo.Addrs[0].String())
		},
		ConnectStart: func(network, addr string) {
			log.Printf("Connection start on network %s addr %s\n",
				network,
				addr)
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("Connection done on network %s  addr %s with error: %s\n",
					network,
					addr,
					err.Error())
			} else {
				log.Printf("Connection done on network %s  addr %s\n",
					network,
					addr)
			}
		},
		WroteHeaders: func() {
			log.Println("Wrote headers")
		},
		WroteRequest: func(wroteRequestInfo httptrace.WroteRequestInfo) {
			if wroteRequestInfo.Err != nil {
				log.Println("Wrote request with error %s\n",
					wroteRequestInfo.Err.Error())
			} else {
				log.Println("Wrote request with no error")
			}
		},
		Wait100Continue: func() {
			log.Println("Wait response code 100 Continue")
		},
		Got100Continue: func() {
			log.Println("Got response 100 Continue")
		},
		GotFirstResponseByte: func() {
			log.Println("Got first response byte")
		},
	})
	req, err := http.NewRequest("GET", "http://www.golang.org/", nil)
	if err != nil {
		log.Fatal("new req", err)
	}
	req = req.WithContext(ctx)
	_, err = http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(os.Stderr, "fetch: %v\n", err)
	}
}
