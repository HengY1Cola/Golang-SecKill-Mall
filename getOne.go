package main

import (
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

//预存商品数量
var productNum int64 = 1000000

//互斥锁
var mutex sync.Mutex

//计数
var count int64 = 0

// GetOneProduct 获取秒杀商品
func GetOneProduct() bool {
	mutex.Lock() //加锁
	defer mutex.Unlock()
	count += 1
	if count%100 == 0 { //判断数据是否超限(运气)
		if sum < productNum {
			sum += 1
			return true
		}
	}
	return false
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("Err:", err)
	}
}
