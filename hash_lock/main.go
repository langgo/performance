package main

import (
	"time"
	"fmt"
	"context"
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

	// storage1 := NewStorage1()

	storage2 := NewStorage2(1)
	go storage2.Run(context.Background())

	test(storage2, data)

	//time.Sleep(5 * time.Second)
	//
	//test(storage2, data)
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
