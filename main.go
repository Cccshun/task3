package main

import (
	"fmt"
	"time"

	"sysu.com/task3/algo"
	"sysu.com/task3/im"
)

func main() {
	im.Init("network/B_BA_200.txt")
	startTime := time.Now()

	ma := &algo.Ma{}
	ma.FindSeed()

	elapsedTime := time.Since(startTime)
	fmt.Printf("运行时间: %s\n", elapsedTime)
}
