package main

import (
	"fmt"
	"time"
)

func main() {
	// 创建一个每秒钟触发一次的定时器
	ticker := time.NewTicker(1 * time.Second)

	// 阻塞直到定时器通道可读
	for t := range ticker.C {
		fmt.Println("Tick at", t)
		// 如果需要停止定时器，可以在这里进行判断并调用 ticker.Stop()
	}
}
