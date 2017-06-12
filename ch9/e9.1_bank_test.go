// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank
import (
	"fmt"
	"testing"

	// "gopl.io/ch9/bank1"
)

type drawStr struct {
	amount int
	ch chan bool
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdraw = make(chan chan int)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	// var draw = drawStr{
	// 	amount: amount,
	// 	ch: make(chan bool),
	// }
	//用来玩的，通道是值， channel is first class
	inC := make(chan int)
	withdraw <- inC
	inC <- amount
	outC := <-withdraw
	ret := <- outC
	if ret == -1 {
		return false
	}
	return true
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		// case draw := <-withdraw:
		// 	ok := false
		// 	if draw.amount <= balance {
		// 		balance -= draw.amount
		// 		ok = true
		// 	}
		// 	go func(){draw.ch <- ok}()
		case inC := <- withdraw:
			amount:= <-inC
			ret := -1
			if amount <= balance {
				balance -= amount
				ret = 0
				}
				outC := make(chan int)
				withdraw <- outC
			go func(){
				outC <- ret
			}()
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		Deposit(200)
		fmt.Println("=", Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		Deposit(150)
		done <- struct{}{}
	}()

	// Tom
	go func() {
		var ok bool
		for {
			ok = Withdraw(300)
			if ok {
				break
			}
		}
		fmt.Println("withraw ", ok)
		done <- struct{}{}
	}()


	// Wait for both transactions.
	<-done
	<-done
	<-done

	if got, want := Balance(), 50; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
