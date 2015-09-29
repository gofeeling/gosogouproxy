package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	refreshDuration = 30 * time.Minute // Regular refresh
	waitDuration    = 10 * time.Second // Force refresh
)

func hostlistDaemon(proxy *SogouProxyHandler, web *WebHandler) {
	isHostValid := make([]bool, proxy.hostMax)
	hostlist := refreshHostlist(proxy, web, isHostValid)
	freshChan := make(chan []int)
	for {
		select {
		case <-time.After(refreshDuration):
			// Regular updating, we don't stop the world.
			go func() { freshChan <- refreshHostlist(proxy, nil, isHostValid) }()
		case getreq := <-proxy.getRequestChan:
			if len(hostlist) > 0 {
				i := rand.Intn(len(hostlist))
				getreq <- hostlist[i]
			} else {
				// Stop and refresh
				hostlist = refreshHostlist(proxy, web, isHostValid)
			}
		case delreq := <-proxy.disableReqestChan:
			if len(hostlist) > 0 {
				isHostValid[delreq] = false
				hostlist = getList(isHostValid)
			} else {
				// Stop and refresh
				hostlist = refreshHostlist(proxy, web, isHostValid)
			}
		case getlistReq := <-web.getlistReqChan:
			getlistReq <- hostlist
		case newlist := <-freshChan:
			// Regular update
			hostlist = newlist
		}
	}
}

func refreshHostlist(proxy *SogouProxyHandler, web *WebHandler, isHostValid []bool) []int {
	var getlistReqChan chan chan []int
	if web != nil {
		getlistReqChan = web.getlistReqChan
	} else {
		getlistReqChan = make(chan chan []int)
		close(getlistReqChan)
	}
	hostlist := refreshHostlistOnce(proxy, isHostValid)
	for {
		if len(hostlist) > 0 {
			log.Println("Available proxy host list is updated.")
			return hostlist
		} else {
			log.Println("All hosts are unavailable. Try again...")
		}
		select {
		case getlistReq, ok := <-getlistReqChan:
			if ok {
				getlistReq <- nil
			}
			hostlist = refreshHostlistOnce(proxy, isHostValid)
		case <-time.After(waitDuration):
			hostlist = refreshHostlistOnce(proxy, isHostValid)
		}
	}
}

func refreshHostlistOnce(proxy *SogouProxyHandler, isHostValid []bool) []int {
	log.Println("Updating available proxy host list...")
	log.Printf("%s -- %s\n", fmt.Sprintf(proxy.hostTemplate, 0), fmt.Sprintf(proxy.hostTemplate, proxy.hostMax-1))
	var waiter sync.WaitGroup
	waiter.Add(proxy.hostMax)
	for i := 0; i < proxy.hostMax; i++ {
		go func(ihost int) {
			proxyHost := fmt.Sprintf(proxy.hostTemplate, ihost)
			conn, err := net.DialTimeout("tcp", proxyHost, proxy.timeOut)
			if err != nil {
				log.Printf("Host %d unavailable: %s\n", ihost, err)
				isHostValid[ihost] = false
			} else {
				log.Printf("Host %d OK (%s).\n", ihost, conn.RemoteAddr())
				conn.Close()
				isHostValid[ihost] = true
			}
			waiter.Done()
		}(i)
	}
	waiter.Wait()
	return getList(isHostValid)
}

func getList(isValid []bool) []int {
	list := make([]int, 0, len(isValid))
	for i := 0; i < len(isValid); i++ {
		if isValid[i] {
			list = append(list, i)
		}
	}
	return list
}
