package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type KV struct {
	k int
	v int
}

func useMu(cnt, cnt2 int) {
	mp := map[int]int{}
	mu := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(cnt)
	st := time.Now()
	for i := 0; i < cnt; i++ {
		go func(idx int) {
			wg.Done()
			for j := 0; j < cnt2; j++ {
				mu.Lock()
				mp[idx] = j
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("use mu      : ", time.Since(st))
}

func useSyncMap(cnt, cnt2 int) {
	mp := sync.Map{}
	wg := sync.WaitGroup{}
	wg.Add(cnt)
	st := time.Now()
	for i := 0; i < cnt; i++ {
		go func(idx int) {
			wg.Done()
			for j := 0; j < cnt2; j++ {
				mp.Store(idx, j)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("use sync map: ", time.Since(st))
}


func useChan(cnt, cnt2 int) {
	mp := map[int]int{}

	wg := sync.WaitGroup{}
	ch := make(chan KV, cnt*100)
	go func() {
		for i := cnt*cnt2; i > 0; i-- {
			kv := <-ch
			mp[kv.k] = kv.v
		}
	}()
	wg.Add(cnt)
	st := time.Now()
	for i := 0; i < cnt; i++ {
		go func(idx int) {
			wg.Done()
			for j := 0; j < cnt2; j++ {
				ch <- KV{idx, j}
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("use chan    : ", time.Since(st))
}

func useAtomic(cnt, cnt2 int) {
	vec := make([]int64, cnt)

	wg := sync.WaitGroup{}
	wg.Add(cnt)
	tcnt2 := int64(cnt2)
	st := time.Now()
	for i := 0; i < cnt; i++ {
		go func(idx int) {
			wg.Done()
			for j := int64(0); j < tcnt2; j++ {
				atomic.StoreInt64(&(vec[idx]), j)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("use Atomic  : ", time.Since(st))
}

func main() {
	cnts := []int{100, 1000, 10000}
	cnt2s := []int{100000000, 1000000000, 10000000000}
	for _, cnt := range cnts {
		for _, cnt2 := range cnt2s {
			fmt.Println("cnt: ", cnt, "cnt2: ", cnt2)
			useMu(cnt, cnt2)
			useChan(cnt, cnt2)
			useSyncMap(cnt, cnt2)
			useAtomic(cnt, cnt2)
		}
	}
}