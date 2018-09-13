package main

import (
	"sync"
	"reflect"
	"unsafe"
	"github.com/pkg/errors"
	"context"
)

var ErrKeyNotExists = errors.New("key not exists")

type Storage interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte)
}

// 有锁的map
type Storage1 struct {
	rwmu sync.RWMutex
	data map[string][]byte
}

func NewStorage1() *Storage1 {
	return &Storage1{
		data: make(map[string][]byte, 64),
	}
}

func (s *Storage1) Put(key []byte, value []byte) {
	k := bytes2string(key)

	s.rwmu.Lock()
	s.data[k] = value
	s.rwmu.Unlock()
}

func (s *Storage1) Get(key []byte) ([]byte, error) {
	k := bytes2string(key)

	s.rwmu.RLock()
	v, ok := s.data[k]
	s.rwmu.RUnlock()

	if !ok {
		return nil, ErrKeyNotExists
	}
	return v, nil
}

func bytes2string(bs []byte) string {
	bssh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))

	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: bssh.Data,
		Len:  bssh.Len,
	}))
}

func string2bytes(s string) []byte {
	ssh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: ssh.Data,
		Len:  ssh.Len,
		Cap:  ssh.Len,
	}))
}

///////////

type KV struct {
	key   []byte
	value []byte
}

type VC struct {
	value []byte
	err   error
}

type KVC struct {
	key []byte
	vc  chan VC
}

type Storage2Runtime struct {
	putC chan KV
	getC chan KVC
	data map[string][]byte
}

func NewStorage2Runtime() *Storage2Runtime {
	return &Storage2Runtime{
		putC: make(chan KV, 16),
		getC: make(chan KVC, 64),
		data: make(map[string][]byte, 64),
	}
}

// TODO 从设计上 ctx 应该是参数，还是结构体上面的字段
func (s2r *Storage2Runtime) Run(ctx context.Context) {
	for {
		select {
		case kvc := <-s2r.getC:
			k := bytes2string(kvc.key)
			value, ok := s2r.data[k]

			var vc VC
			if ok {
				vc = VC{
					value: value,
					err:   nil,
				}
			} else {
				vc = VC{
					value: nil,
					err:   ErrKeyNotExists,
				}
			}
			kvc.vc <- vc
		case kv := <-s2r.putC:
			k := bytes2string(kv.key)
			s2r.data[k] = kv.value
		case <-ctx.Done():
			return
		}
	}
}

func (s2r *Storage2Runtime) Get(key []byte) ([]byte, error) {
	kvc := KVC{
		key: key,
		vc:  make(chan VC, 1),
	}

	s2r.getC <- kvc
	vc := <-kvc.vc
	return vc.value, vc.err
}

func (s2r *Storage2Runtime) Put(key []byte, value []byte) {
	s2r.putC <- KV{
		key:   key,
		value: value,
	}
}

type Storage2 struct {
	runtimes []*Storage2Runtime
}

func NewStorage2(count int) *Storage2 {
	s := &Storage2{
		runtimes: make([]*Storage2Runtime, count),
	}
	for i := 0; i < len(s.runtimes); i++ {
		s.runtimes[i] = NewStorage2Runtime()
	}
	return s
}

func (s *Storage2) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < len(s.runtimes); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			s.runtimes[i].Run(ctx)
		}(i)
	}
	wg.Wait()
}

func (s *Storage2) Get(key []byte) ([]byte, error) {
	return s.r(key).Get(key)
}

func (s *Storage2) Put(key []byte, value []byte) {
	s.r(key).Put(key, value)
}

func (s *Storage2) r(key []byte) *Storage2Runtime {
	h := hash(key)
	return s.runtimes[h%len(s.runtimes)]
}

func hash(key []byte) int {
	h := 5381

	l := 2
	if len(key) < l {
		l = len(key)
	}

	for i := 0; i < l; i++ {
		h = ((h << 5) + h) + int(key[i]) // h * 33 + c
	}
	return h
}

type Storage3 struct {
	data map[string][]byte
}

func NewStorage3() *Storage1 {
	return &Storage1{
		data: make(map[string][]byte, 64),
	}
}

func (s *Storage3) Put(key []byte, value []byte) {
	k := bytes2string(key)

	s.data[k] = value
}

func (s *Storage3) Get(key []byte) ([]byte, error) {
	k := bytes2string(key)

	v, ok := s.data[k]

	if !ok {
		return nil, ErrKeyNotExists
	}
	return v, nil
}
