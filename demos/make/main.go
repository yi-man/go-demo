package main

import (
	"fmt"
	"time"
)

func main() {
	slice := make([]int, 10, 100)
	for i := 0; i < 10; i++ {
		slice[i] = i
	}
	for i, s := range slice {
		fmt.Println(i, s)
	}
	fmt.Println(slice, len(slice))

	hash := make(map[int]float32, 2)
	hash[1] = 1.0
	hash[2] = 2.0
	hash[3] = 3.0
	hash[4] = 4.0
	for key, value := range hash {
		fmt.Printf("key is: %d - value is: %f\n", key, value)
	}
	fmt.Println(hash, len(hash))

	// 容量(capacity)代表Channel容纳的最多的元素的数量，代表Channel的缓存的大小。
	// 如果没有设置容量，或者容量设置为0, 说明Channel没有缓存，只有sender和receiver都准备好了后它们的通讯(communication)才会发生(Blocking)。如果设置了缓存，就有可能不发生阻塞， 只有buffer满了后 send才会阻塞， 而只有缓存空了后receive才会阻塞。一个nil channel不会通信
	ch := make(chan int, 5)
	ch <- 100
	v, ok := <-ch
	fmt.Println(v, ok)

	c := make(chan int)
	defer close(c)
	go func() { c <- 3 + 4 }()
	i, ok2 := <-c
	fmt.Println(i, ok2)

	fmt.Println(`-------------sum-----------`)
	s := []int{7, 2, 8, -9, 4, 0}
	c2 := make(chan int)
	go sum(s[:len(s)/2], "first half", c2)
	go sum(s[len(s)/2:], "last half", c2)
	x, y := <-c2, <-c2 // receive from c2
	fmt.Println(x, y, x+y)

	fmt.Println(`-------------range-----------`)
	go func() {
		time.Sleep(1 * time.Hour)
	}()
	c3 := make(chan int)
	go func() {
		for i := 0; i < 10; i = i + 1 {
			c3 <- i
		}
		close(c3)
	}()
	// range c3产生的迭代值为Channel中发送的值，它会一直迭代直到channel被关闭。如果把close(c)注释掉，程序会一直阻塞在for …… range那一行。
	for i := range c3 {
		fmt.Println(i)
	}
	fmt.Println("Finished")

	fmt.Println(`-------------select-----------`)
	c4 := make(chan int, 20)
	quit := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i, <-c4)
		}
		quit <- 0
	}()
	fibonacci(c4, quit)

	fmt.Println(`-------------timeout-----------`)
	c5 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c5 <- "result 1"
	}()
	select {
	case res := <-c5:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}

	fmt.Println(`-------------Timer和Ticker-----------`)
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C
	fmt.Println("Timer 1 expired")
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
}

func sum(s []int, id string, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	fmt.Println(s, id)
	c <- sum // send sum to c
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
