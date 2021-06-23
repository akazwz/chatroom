package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan []byte)
	go func() {
		str := "this"
		for i := 0; i < 5; i++ {
			str = str[1:]
			ch <- []byte(str)
			time.Sleep(time.Second)
		}
	}()

	for bytes := range ch {
		fmt.Println(string(bytes))
		if string(bytes) == "s" {
			break
		}
	}
}
