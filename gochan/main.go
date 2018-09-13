package main

import (
	"bytes"
	"github.com/langgo/performance/ttool"
	"sync"
	"time"
	"fmt"
)

func main() {
	data := ttool.GenTestData(1000*10000, 32, 8)

	c := make(chan S, 64)

	st := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		comsumer(c)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		producer(c, data)
	}()
	wg.Wait()

	s := time.Now().Sub(st).Seconds()
	fmt.Printf("%f s, %f w op/s", s, float64(len(data))/s/10000)
}

type S struct {
	key   []byte
	value []byte
}

func comsumer(c <-chan S) {
	for s := range c {
		bytes.Equal(s.key, s.value)
	}
}

func producer(c chan<- S, data []ttool.KV) {
	for i := range data {
		c <- S{
			key:   data[i].Key,
			value: data[i].Value,
		}
	}
	close(c)
}
