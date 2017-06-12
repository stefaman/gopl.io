// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 272.

// Package memotest provides common functions for
// testing various designs of the memo package.
package memory_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
	// "math/rand"
	"gopl.io/ch9/memotest"
)

//!+httpRequestBody
func httpGetBody(url string, dones ...chan struct{}) (interface{}, error) {
	var done chan struct{}
	if len(dones) > 0 {
		done = dones[0]
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//!-httpRequestBody

var HTTPGetBody = httpGetBody

type M interface {
	Get(key string, dones ...chan struct{}) (interface{}, error)
}

/*
//!+seq
	m := memo.New(httpGetBody)
//!-seq
*/
const squenDuration = 100*time.Millisecond
const conDuration = 100*time.Millisecond
func Sequential(t *testing.T, m M) {
	//!+seq
	for in := range memotest.IncomingURLs() {
		done := make(chan struct{})
		go func(){
			select{
			case <- time.After(squenDuration):
				close(done)
			}
		}()
		start := time.Now()
		value, err := m.Get(in, done)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%s, %s, %d bytes\n",
			in, time.Since(start), len(value.([]byte)))
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
	for in := range memotest.IncomingURLs() {
		n.Add(1)
		go func(in string) {
			defer n.Done()
			done := make(chan struct{})
			go func(){
				select{
				case <- time.After(conDuration):
					close(done)
				}
			}()
			start := time.Now()
			value, err := m.Get(in, done)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n",
				in, time.Since(start), len(value.([]byte)))
		}(in)
	}
	n.Wait()
	//!-conc
}
