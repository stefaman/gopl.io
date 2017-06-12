//e9.3
//stefaman 20170525
package memory

import(
	"fmt"
	"sync"
)


type Func func(string, ...chan struct{}) (interface{}, error)
type result struct{
	val interface{}
	err error
}
type entry struct {
	res result
	ready chan struct{}
	ok bool
}
type Memo struct {
	f Func
	mu sync.Mutex
	cache map[string]*entry
}


func New(f Func) *Memo {
	return &Memo{
		f: f,
		cache: make(map[string]*entry),
	}
}
/*version 1, bad performance
var mu sync.Mutex
func (m *Memo) Get(str string, dones ...chan struct{}) (interface{}, error) {
	var done chan struct{}
	if len(dones) > 0 {
		done = dones[0]
	}
	mu.Lock()
	e, ok := m.cache[str]
	if !ok {
		e = &entry{
			ready: make(chan struct{}),
		}
		e.res.val, e.res.err = m.f(str, done)
		select {
		case <- done:
			mu.Unlock()
			//  return e.res.val, fmt.Errorf("Get %s: request cancled while wating: %v", str, e.res.err)
			//go to return m.f's return values
		default:
			m.cache[str] = e
			mu.Unlock()
			close(e.ready)
		}
	}else {
		mu.Unlock()
		select {//如果同时closed done and e.ready, 随机性返回
		case <-done:
			return nil, fmt.Errorf("Get %s: request cancled while wating", str)
		case <-e.ready:
			//go to return
		}
	}
	return e.res.val, e.res.err
}
*/

//cancel channel done only has effect when waiting, for m.f calling or for e.ready

func (m *Memo) Get(str string, dones ...chan struct{}) (interface{}, error) {
	var done chan struct{}
	if len(dones) > 0 {
		done = dones[0]
	}
	RETRY:
	m.mu.Lock()
	e := m.cache[str]
	if e == nil {//"if e== nil || !e.ok{" //bug, reace: reading e.ok and e.ok = true
		e = &entry{ready: make(chan struct{})}
		m.cache[str] = e
		m.mu.Unlock()
		e.res.val, e.res.err = m.f(str, done)

		select {
		case <- done:
			//set e.ok false
			//go to return m.f's return values
		default:
			e.ok = true
			close(e.ready)
			// go to return
		}
	}else {
		m.mu.Unlock()
		select {//如果同时closed done and e.ready, 随机性返回
		case <-done:// canled while waiting
			return nil, fmt.Errorf("Get %s: request canceled while wating", str)
		case <-e.ready:
			if !e.ok {//m.f() calling goroutine canceled
				goto RETRY//这种处理方式待商榷
				//done如果采用延时关闭的方式，延时关闭后，所有在排队的协程 go to “RETRY”, 如果同一str的延时是一致的，那么就会造成雪崩，需要考虑采用延时递增的方式的，或者出错。
				return nil, fmt.Errorf("Get %s: request canceled as first goroutine canceled", str)
			}
		}
	}
	return e.res.val, e.res.err
}

//copy from memo5/main.go, for performance compare
func (memo *Memo) Get1(key string, dones ...chan struct{}) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.val, e.res.err = memo.f(key)

		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	return e.res.val, e.res.err
}
