// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 272.

// Package memotest provides common functions for
// testing various designs of the memo package.
package memotest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
	"math/rand"
)

//!+httpRequestBody
func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//!-httpRequestBody

var HTTPGetBody = httpGetBody

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"http://gopl.io",
			"http://www.baidu.com",
			"http://www.baidu.com",
			"http://gopl.io",
			"http://gopl.io",
			"https://music.163.com",
			"http://www.baidu.com",
			"http://gopl.io",
			"http://www.baidu.com",
			"http://www.baidu.com",
			"http://www.baidu.com",
			"https://music.163.com",
			"https://music.163.com",
			"http://gopl.io",
			"http://gopl.io",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

//stefaman
var inputs = []string{
	// "http://gopl.io",
	"http://www.baidu.com",
	"http://music.163.com",
	"https://www.zhihu.com/",
	"http://www.sina.com.cn/",
}
const nInputs = 10000
func IncomingURLs() <-chan string {
	l := int32(len(inputs))
	ch := make(chan string)
	go func() {
		for i := 0; i < nInputs; i++ {
			ch <- inputs[uint(rand.Int31n(l))]
		}
		close(ch)
	}()
	return ch
}
type M interface {
	Get(key string) (interface{}, error)
}

/*
//!+seq
	m := memo.New(httpGetBody)
//!-seq
*/

func Sequential(t *testing.T, m M) {
	//!+seq
	for url := range IncomingURLs() {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
	//!-seq
}

/*
//!+conc
	m := memo.New(httpGetBody)
//!-conc
*/

func Concurrent(t *testing.T, m M) {
	//!+conc
	var n sync.WaitGroup
	for url := range IncomingURLs() {
		n.Add(1)
		go func(url string) {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n",
				url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	n.Wait()
	//!-conc
}
