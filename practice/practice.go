package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Shuffle1(s []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(s))
	n := len(s)
	for i := 0; i < n; i++ {
		randIndex := r.Intn(len(s))
		ret[i] = s[randIndex]
		s = append(s[:randIndex], s[randIndex+1:]...)
	}
	return ret
}

func Shuffle(s []int) []int {
	return s
}

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	shuffle := Shuffle(s)
	fmt.Println(shuffle)
}
