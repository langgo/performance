package main

import (
	"time"
	"fmt"
	"sync"
)

// 构造测试数据
// 怎么描述测试数据
// 怎么接口化，怎么配置化

// 测试数据的格式
//    key-value key的长度，value的长度

// 测试数据的分布
//    28原则 20%的数据占据80%的访问，80%的数据占据20%的访问

// 测试数据总量

func main() {
	data := genTestData(80*10000, 32, 8)
	//for i := range kvs {
	//	fmt.Printf("%s: %v\n", string(kvs[i].key), string(kvs[i].value))
	//}

	//{
	//	{
	//		fmt.Println("storage1")
	//
	//		storage1 := NewStorage1()
	//
	//		test(storage1, data)
	//	}
	//
	//	time.Sleep(2 * time.Second)
	//
	//	{
	//		fmt.Println("storage2")
	//
	//		storage2 := NewStorage2(1)
	//		go storage2.Run(context.Background())
	//
	//		test(storage2, data)
	//	}
	//}

	//time.Sleep(2 * time.Second)

	//
	//{
	//	{
	//		fmt.Println("storage1")
	//
	//		storage1 := NewStorage1()
	//
	//		testc(4, storage1, data)
	//
	//		_ = storage1
	//	}
	//
	//	time.Sleep(2 * time.Second)
	//
	//	{
	//		fmt.Println("storage2")
	//
	//		storage2 := NewStorage2(1)
	//		go storage2.Run(context.Background())
	//
	//		testc(4, storage2, data)
	//
	//		_ = storage2
	//	}
	//}

	//test3(func() Storage {
	//	return NewStorage1()
	//}, 4, data)
	//
	//time.Sleep(2 * time.Second)

	test3(func() Storage {
		return NewStorage3()
	}, 4, data)
}

func test3(newf func() Storage, c int, data []KV) {
	datac := len(data) / c

	fmt.Printf("Count: %d\n", len(data))

	storages := make([]Storage, c)
	for i := 0; i < c; i++ {
		storages[i] = newf()
	}

	var wg sync.WaitGroup
	{
		st := time.Now()

		for i := 0; i < c; i++ {
			wg.Add(1)

			d := data[datac*i : datac*(i+1)]
			storage1 := storages[i]

			go func() {
				defer wg.Done()

				for _, kv := range d {
					storage1.Put(kv.key, kv.value)
				}
			}()
		}
		wg.Wait()

		s := time.Now().Sub(st).Seconds()
		fmt.Printf("PUT: %f ms, %f w op/s\n", s*1000, float64(len(data))/s/10000)
	}

	{
		st := time.Now()

		for i := 0; i < c; i++ {
			wg.Add(1)

			d := data[datac*i : datac*(i+1)]
			storage1 := storages[i]

			go func() {
				defer wg.Done()

				for _, kv := range d {
					_, err := storage1.Get(kv.key)
					if err != nil {
						panic(err)
					}
				}
			}()
		}
		wg.Wait()

		s := time.Now().Sub(st).Seconds()
		fmt.Printf("GET: %f ms, %f w op/s\n", s*1000, float64(len(data))/s/10000)
	}
}

func testc(c int, storage Storage, data []KV) {
	fmt.Printf("Count: %d\n", len(data))
	var wg sync.WaitGroup

	data4 := len(data) / c

	{
		st := time.Now()

		for i := 0; i < c; i++ {
			wg.Add(1)

			d := data[data4*i : data4*(i+1)]
			go func(data []KV) {
				defer wg.Done()

				for _, kv := range data {
					storage.Put(kv.key, kv.value)
				}
			}(d)
		}

		wg.Wait()
		s := time.Now().Sub(st).Seconds()
		fmt.Printf("PUT: %f ms, %f w op/s\n", s*1000, float64(len(data))/s/10000)
	}

	{
		st := time.Now()

		for i := 0; i < c; i++ {
			wg.Add(1)

			d := data[data4*i : data4*(i+1)]
			go func(data []KV) {
				defer wg.Done()

				for _, kv := range data {
					_, err := storage.Get(kv.key)
					if err != nil {
						panic(err)
					}
				}
			}(d)
		}

		wg.Wait()
		s := time.Now().Sub(st).Seconds()
		fmt.Printf("GET: %f ms, %f w op/s\n", s*1000, float64(len(data))/s/10000)
	}
}

func test(storage Storage, data []KV) {
	fmt.Printf("Count: %d\n", len(data))
	{
		st := time.Now()

		for _, kv := range data {
			storage.Put(kv.key, kv.value)
		}

		s := time.Now().Sub(st).Seconds()
		fmt.Printf("PUT: %f ms, %f w op/s\n", s*1000, float64(len(data))/s/10000)
	}

	{
		st := time.Now()

		for _, kv := range data {
			_, err := storage.Get(kv.key)
			if err != nil {
				panic(err)
			}
		}

		s := time.Now().Sub(st).Seconds()
		fmt.Printf("GET: %f ms, %f w op/s\n", s*1000, float64(len(data))/s/10000)
	}
}
