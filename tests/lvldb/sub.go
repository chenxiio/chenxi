package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/syndtr/goleveldb/leveldb"
// 	"github.com/syndtr/goleveldb/leveldb/util"
// )

// func main() {
// 	db, err := leveldb.OpenFile("./testdb", nil)
// 	if err != nil {
// 		fmt.Println("Error opening database:", err)
// 		return
// 	}
// 	defer db.Close()
// 	// // 创建一个前缀为“example”的迭代器
// 	// iter := db.NewIterator(util.BytesPrefix([]byte("example")), nil)
// 	// for iter.Next() {
// 	// 	key := iter.Key()
// 	// 	value := iter.Value()
// 	// 	fmt.Printf("Key: %s, Value: %s\n", key, value)
// 	// }
// 	// iter.Release()
// 	// err = iter.Error()
// 	// if err != nil {
// 	// 	fmt.Println("Error iterating over database:", err)
// 	// 	return
// 	// }
// 	// 创建一个监听器
// 	changes, err := db.CompactRange(util.Range{})
// 	if err != nil {
// 		fmt.Println("Error creating changes listener:", err)
// 		return
// 	}
// 	defer changes.Close()
// 	// 创建一个等待组，以确保所有线程都完成后退出
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	// 启动一个线程不断地修改值
// 	go func() {
// 		for {
// 			err := db.Put([]byte("example_key"), []byte(time.Now().String()), nil)
// 			if err != nil {
// 				fmt.Println("Error putting value:", err)
// 			}
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	// 监听值的变化
// 	for {
// 		select {
// 		case change := <-changes.C:
// 			fmt.Printf("Key: %s, Value: %s\n", change.Key, change.Value)
// 		}
// 	}
// 	// 等待所有线程完成
// 	wg.Wait()
// }
